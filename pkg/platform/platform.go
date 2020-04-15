package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/pool"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/queue"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/workers"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
	"gopkg.in/go-playground/validator.v9"
)

// IPlatform defines the set of methods of a Platform.
type IPlatform interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Platform represents a set of components to run applications.
type Platform struct {
	Server           *ign.Server
	Logger           ign.Logger
	Context          context.Context
	Email            *email.Email
	Validator        *validator.Validate
	FormDecoder      *form.Decoder
	Transport        *transport.Transport
	Orchestrator     *orchestrator.Kubernetes
	CloudProvider    *cloud.AmazonWS
	Permissions      *permissions.Permissions
	UserService      *users.Service
	Config           Config
	Simulator        simulator.ISimulator
	PoolFactory      pool.Factory
	Scheduler        *scheduler.Scheduler
	LaunchQueue      queue.IQueue
	TerminationQueue chan workers.TerminateDTO
	LaunchPool       pool.IPool
	TerminationPool  pool.IPool
}

// Name returns the platform name
func (p *Platform) Name() string {
	return "cloudsim"
}

// NewSimulator returns a new Platform from the given configuration.
func New(config Config) *Platform {
	p := Platform{}
	p.Config = config

	p.setupLogger()
	p.Logger.Debug("[INIT] Logger initialized.")

	// TODO: Decide where the score generation should go

	p.setupContext()
	p.Logger.Debug("[INIT] Context initialized.")

	p.setupServer()
	p.Logger.Debug(fmt.Sprintf("[INIT] Server initialized using HTTP port [%s] and SSL port [%s].", p.Server.HTTPPort, p.Server.SSLport))
	p.Logger.Debug(fmt.Sprintf("[INIT] Database [%s] initialized", p.Server.DbConfig.Name))

	p.setupRouter()
	p.Logger.Debug("[INIT] Router initialized.")

	p.setupValidator() // TODO: Decide where should the custom validators should go
	p.Logger.Debug("[INIT] Validators initialized.")

	p.setupFormDecoder()
	p.Logger.Debug("[INIT] Form decoder initialized.")

	p.setupPermissions()
	p.Logger.Debug("[INIT] Permissions initialized.")

	p.setupUserService()
	p.Logger.Debug("[INIT] User service initialized")

	p.setupDatabase()
	p.Logger.Debug("[INIT] Database configured: Migration, default data and custom indexes.")

	p.setupCloudProvider()
	p.Logger.Debug("[INIT] Cloud provider initialized: AWS.")

	p.setupOrchestrator()
	p.Logger.Debug("[INIT] Orchestrator initialized: Kubernetes.")

	p.setupSimulator()
	p.Logger.Debug("[INIT] Simulator initialized. Using: AWS and Kubernetes.")

	p.setupScheduler()
	p.Logger.Debug("[INIT] Scheduler initialized.")

	p.setupQueues()
	p.Logger.Debug("[INIT] RequestLaunch and termination queues have been initialized.")

	p.setupPoolFactory()
	if _, err := p.setupWorkers(); err != nil {
		p.Logger.Critical("[INIT|CRITICAL] Could not initialize workers.")
	}
	p.Logger.Debug("[INIT] RequestLaunch and termination workers have been initialized.")

	if _, err := p.setupTransport(); err != nil {
		p.Logger.Critical("[INIT|CRITICAL] Could not initialize transport.")
	}
	p.Logger.Debug("[INIT] Transport initialized. Using: IGN Transport.")
	return &p
}

// Start starts the platform.
func (p *Platform) Start(ctx context.Context) error {
	go func() {
		for {
			var element interface{}
			var err *ign.ErrMsg
			if element, err = p.LaunchQueue.DequeueOrWait(); err != nil {
				continue
			}

			dto, ok := element.(workers.LaunchDTO)
			if !ok {
				continue
			}

			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] About to process launch action. Group ID: [%s]", dto.GroupID))
			if err := p.LaunchPool.Serve(dto); err != nil {
				p.Logger.Error(fmt.Sprintf("[QUEUE|LAUNCH] Error while serving launch action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] The launch action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
		}
	}()

	go func() {
		for dto := range p.TerminationQueue {
			p.Logger.Info(fmt.Sprintf("[QUEUE|TERMINATE] About to process terminate action. Group ID: [%s]", dto.GroupID))
			if err := p.TerminationPool.Serve(dto); err != nil {
				p.Logger.Error(fmt.Sprintf("[QUEUE|TERMINATE] Error while serving terminate action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger.Info(fmt.Sprintf("[QUEUE|TERMINATE] The terminate action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
		}
	}()
	return nil
}

// Stop stops the platform.
func (p *Platform) Stop(ctx context.Context) error {
	close(p.TerminationQueue)
	return nil
}

// RequestLaunch enqueues a launch action to launch a simulation from the given Group ID.
func (p *Platform) RequestLaunch(ctx context.Context, groupID string) {
	job := workers.LaunchDTO{
		GroupID: groupID,
		Action: nil,
	}
	p.LaunchQueue.Enqueue(job)
}

// RequestTermination enqueues a termination action to terminate a simulation from the given Group ID.
func (p *Platform) RequestTermination(ctx context.Context, groupID string) {
	job := workers.TerminateDTO{
		GroupID: groupID,
		Action: nil,
	}
	p.TerminationQueue <- job
}