package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/handlers"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/pool"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/queue"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/workers"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
	"gopkg.in/go-playground/validator.v9"
)

// IPlatform defines the set of methods of a platform.
type IPlatform interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RequestLaunch(ctx context.Context, groupID string)
	RequestTermination(ctx context.Context, groupID string)
	Logger() ign.Logger
	Context() context.Context
}

// Platform represents a set of components to run applications.
type Platform struct {
	Server           *ign.Server
	logger           ign.Logger
	context          context.Context
	email            email.Email
	Validator        *validator.Validate
	FormDecoder      *form.Decoder
	Transport        *transport.Transport
	Orchestrator     orchestrator.Kubernetes
	CloudProvider    *cloud.AmazonWS
	Permissions      *permissions.Permissions
	UserService      users.Service
	Config           Config
	Simulator        simulator.ISimulator
	PoolFactory      pool.Factory
	Scheduler        *scheduler.Scheduler
	LaunchQueue      queue.Queue
	TerminationQueue chan workers.TerminateInput
	LaunchPool       pool.Pool
	TerminationPool  pool.Pool
	Controllers      controllers
}

func (p *Platform) Logger() ign.Logger {
	return p.logger
}

func (p *Platform) Context() context.Context {
	return p.context
}

func (p *Platform) Email() email.Email {
	return p.email
}

// TODO: Add initializer for queue controller.
type controllers struct {
	Queue queue.IController
}

// Name returns the platform's name
func (p *Platform) Name() string {
	return "cloudsim"
}

// NewSimulator returns a new platform from the given configuration.
func New(config Config) IPlatform {
	p := Platform{}
	p.Config = config

	p.setupLogger()
	p.Logger().Debug("[INIT] Logger initialized.")

	// TODO: Decide where the score generation should go

	p.setupContext()
	p.Logger().Debug("[INIT] Context initialized.")

	p.setupServer()
	p.Logger().Debug(fmt.Sprintf("[INIT] Server initialized using HTTP port [%s] and SSL port [%s].", p.Server.HTTPPort, p.Server.SSLport))
	p.Logger().Debug(fmt.Sprintf("[INIT] Database [%s] initialized", p.Server.DbConfig.Name))

	p.setupRouter()
	p.Logger().Debug("[INIT] Router initialized.")

	// TODO: Decide where should the custom validators should go
	p.setupValidator()
	p.Logger().Debug("[INIT] Validators initialized.")

	p.setupFormDecoder()
	p.Logger().Debug("[INIT] Form decoder initialized.")

	p.setupPermissions()
	p.Logger().Debug("[INIT] Permissions initialized.")

	p.setupUserService()
	p.Logger().Debug("[INIT] User service initialized")

	p.setupDatabase()
	p.Logger().Debug("[INIT] Database configured: Migration, default data and custom indexes.")

	p.setupCloudProvider()
	p.Logger().Debug("[INIT] Cloud provider initialized: AWS.")

	p.setupOrchestrator()
	p.Logger().Debug("[INIT] Orchestrator initialized: k8s.")

	p.setupSimulator()
	p.Logger().Debug("[INIT] Simulator initialized. Using: AWS and k8s.")

	p.setupScheduler()
	p.Logger().Debug("[INIT] Scheduler initialized.")

	p.setupQueues()
	p.Logger().Debug("[INIT] RequestLaunch and termination queues have been initialized.")

	p.setupPoolFactory()
	if _, err := p.setupWorkers(); err != nil {
		p.Logger().Critical("[INIT|CRITICAL] Could not initialize workers.")
	}
	p.Logger().Debug("[INIT] RequestLaunch and termination workers have been initialized.")

	if _, err := p.setupTransport(); err != nil {
		p.Logger().Critical("[INIT|CRITICAL] Could not initialize transport.")
	}
	p.Logger().Debug("[INIT] Transport initialized. Using: IGN Transport.")
	return &p
}

// Launch starts the platform.
func (p *Platform) Start(ctx context.Context) error {
	go func() {
		for {
			var element interface{}
			var err *ign.ErrMsg
			if element, err = p.LaunchQueue.DequeueOrWait(); err != nil {
				continue
			}

			dto, ok := element.(workers.LaunchInput)
			if !ok {
				continue
			}

			p.Logger().Info(fmt.Sprintf("[QUEUE|LAUNCH] About to process launch action. Group ID: [%s]", dto.GroupID))
			if err := p.LaunchPool.Serve(dto); err != nil {
				p.Logger().Error(fmt.Sprintf("[QUEUE|LAUNCH] Error while serving launch action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger().Info(fmt.Sprintf("[QUEUE|LAUNCH] The launch action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
		}
	}()

	go func() {
		for dto := range p.TerminationQueue {
			p.Logger().Info(fmt.Sprintf("[QUEUE|TERMINATE] About to process terminate action. Group ID: [%s]", dto.GroupID))
			if err := p.TerminationPool.Serve(dto); err != nil {
				p.Logger().Error(fmt.Sprintf("[QUEUE|TERMINATE] Error while serving terminate action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger().Info(fmt.Sprintf("[QUEUE|TERMINATE] The terminate action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
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
	job := workers.LaunchInput{
		GroupID: groupID,
		Action:  nil,
	}
	p.LaunchQueue.Enqueue(job)
}

// RequestTermination enqueues a termination action to terminate a simulation from the given Group ID.
func (p *Platform) RequestTermination(ctx context.Context, groupID string) {
	job := workers.TerminateInput{
		GroupID: groupID,
		Action:  nil,
	}
	p.TerminationQueue <- job
}

func (p *Platform) registerRoutes() {
	router.ConfigureRoutes(p.Server, "2.0", "", p.getLaunchQueueRoutes())
}

func (p *Platform) getLaunchQueueRoutes() ign.Routes {
	return ign.Routes{
		ign.Route{
			Name:        "Get all elements from queue",
			Description: "Get all elements from queue. This route should optionally be able to handle pagination parameters.",
			URI:         "/queue",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get all elements from queue. This route should optionally be able to handle pagination parameters",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.GetAll)),
						},
					},
				},
			},
		},
		// Launch queue - Count elements
		ign.Route{
			Name:        "Count elements in the queue",
			Description: "Get the amount of elements in the queue",
			URI:         "/queue/count",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get the amount of elements in the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.Count)),
						},
					},
				},
			},
		},
		// Launch queue - Swap elements
		ign.Route{
			Name:        "Swap queue elements moving A to B and vice versa",
			Description: "Swap queue elements moving A to B and vice versa",
			URI:         "/queue/{groupIDA}/swap/{groupIDB}",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Swap queue elements moving A to B and vice versa",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.Swap)),
						},
					},
				},
			},
		},
		// Launch queue - Move to front
		ign.Route{
			Name:        "Move an element to the front of the queue",
			Description: "Move an element to the front of the queue",
			URI:         "/queue/{groupID}/move/front",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Move an element to the front of the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.MoveToFront)),
						},
					},
				},
			},
		},
		// Launch queue - Move to back
		ign.Route{
			Name:        "Move an element to the back of the queue",
			Description: "Move an element to the back of the queue",
			URI:         "/queue/{groupID}/move/back",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Move an element to the back of the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.MoveToBack)),
						},
					},
				},
			},
		},
		// Launch queue - Remove an element
		ign.Route{
			Name:        "Remove an element from the queue",
			Description: "Remove an element from the queue",
			URI:         "/queue/{groupID}",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "DELETE",
					Description: "Remove an element from the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(handlers.WithUser(p.UserService, p.Controllers.Queue.Remove)),
						},
					},
				},
			},
		},
	}
}
