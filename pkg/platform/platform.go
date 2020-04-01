package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/pool"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/queue"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transporter"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
)

type IPlatform interface {
	Name() string
}

type Platform struct {
	Server        *ign.Server
	Logger        ign.Logger
	Context       context.Context
	Validator     *validator.Validate
	FormDecoder   *form.Decoder
	Transporter   *transporter.Transporter
	Orchestrator  *orchestrator.Kubernetes
	CloudProvider *cloud.AmazonWS
	Permissions   *permissions.Permissions
	UserService   *users.Service
	Config        Config
	NodeManager   *simulator.Manager
	Applications  map[string]*application.IApplication
	PoolFactory   pool.Factory
	Scheduler *scheduler.Scheduler
	LaunchQueue *queue.IQueue
	TerminationQueue chan string
	LaunchPool pool.IPool
	TerminationPool pool.IPool
}

// Name returns the platform name
func (p Platform) Name() string {
	return "cloudsim"
}

// New returns a new application from the given configuration.
func New(config Config) Platform {
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

	// TODO: Move http handler instance logic to applications

	// TODO: Initialize applications? We might need to register them instead.

	p.setupNodeManager()
	p.Logger.Debug("[INIT] Manager initialized. Using: AWS and Kubernetes.")

	p.setupScheduler()
	p.Logger.Debug("[INIT] Scheduler initialized.")

	p.setupQueues()
	p.Logger.Debug("[INIT] Launch and termination queues have been initialized.")

	p.setupPoolFactory()
	p.setupWorkers()
	p.Logger.Debug("[INIT] Launch and termination workers have been initialized.")
	return p
}