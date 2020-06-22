package platform

import (
	"context"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db/migrations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/pool"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/queue"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/workers"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

// Setup represent a set of methods to initialize the platform.
type Setup interface {
	setupLogger() Platform
	setupContext() Platform
	setupServer() Platform
	setupRouter() Platform
	setupEmail() Platform
	setupValidator() Platform
	setupFormDecoder() Platform
	setupPermissions() Platform
	setupUserService() Platform
	setupDatabase() Platform
	setupCloudProvider() Platform
	setupOrchestrator() Platform
	setupSimulator() Platform
	setupPoolFactory() Platform
	setupScheduler() Platform
	setupQueues() Platform
	setupWorkers() (Platform, error)
	setupTransport() (Platform, error)
}

// setupLogger initializes the logger.
func (p *platform) setupLogger() Platform {
	l, err := logger.New()
	if err != nil {
		log.Fatalf("Error parsing environment variables for Logger. %+v\n", err)
	}
	p.logger = l
	return p
}

// setupContext initializes the context.
func (p *platform) setupContext() Platform {
	ctx := ign.NewContextWithLogger(context.Background(), p.logger)
	p.context = ctx
	return p
}

// setupServer initializes the HTTP server.
// If there is an error, it will panic.
func (p *platform) setupServer() Platform {
	cfg := server.Config{
		Auth0:    p.Config.Auth0,
		HTTPport: p.Config.HTTPport,
		SSLport:  p.Config.SSLport,
	}
	s, err := server.New(cfg)
	if err != nil {
		p.Logger().Critical(err)
		log.Fatalf("Error while initializing server. %v\n", err)
	}
	p.server = s
	p.Server().DbConfig = db.NewConfig()
	return p
}

// setupRouter initializes the server's router.
func (p *platform) setupRouter() Platform {
	r := router.New()
	p.Server().SetRouter(r)
	return p
}

// setupEmail initializes the email service.
func (p *platform) setupEmail() Platform {
	e := email.New()
	p.email = e
	return p
}

// setupValidator initializes the validator.
func (p *platform) setupValidator() Platform {
	validate := validator.New()
	p.validator = validate
	return p
}

// setupFormDecoder initializes the form decoder.
func (p *platform) setupFormDecoder() Platform {
	p.formDecoder = form.NewDecoder()
	return p
}

// setupPermissions initializes the platform permissions.
func (p *platform) setupPermissions() Platform {
	per := &permissions.Permissions{}
	err := per.Init(p.Database(), p.Config.SysAdmin)
	if err != nil {
		p.Logger().Critical(err)
		log.Fatalf("Error while initializing server. %v\n", err)
	}
	p.Permissions = per
	return p
}

// setupUserService initializes the User service.
func (p *platform) setupUserService() Platform {
	s, err := users.NewService(p.Permissions, p.Config.SysAdmin)
	if err != nil {
		p.Logger().Critical(err)
		log.Fatalf("Error while configuring user service. %v\n", err)
	}
	p.services.user = s
	return p
}

// setupDatabase performs migrations, adds default data and adds custom indexes.
func (p *platform) setupDatabase() Platform {
	migrations.Migrate(p.Context(), p.Database())
	migrations.AddDefaultData(p.Context(), p.Database())
	migrations.AddCustomIndexes(p.Context(), p.Database())
	return p
}

// setupCloudProvider initializes the Cloud Provider.
func (p *platform) setupCloudProvider() Platform {
	p.aws = cloud.New()
	return p
}

// setupOrchestrator initializes the container k8s.
func (p *platform) setupOrchestrator() Platform {
	p.k8s = orchestrator.New()
	return p
}

// setupSimulator initializes the simulator.
func (p *platform) setupSimulator() Platform {
	input := simulator.NewSimulatorInput{}
	p.simulator = simulator.NewSimulator(input)
	return p
}

// setupPoolFactory initializes the Default Pool Factory.
func (p *platform) setupPoolFactory() Platform {
	p.PoolFactory = pool.NewPool
	return p
}

// setupScheduler gets the instance from the scheduler package.
func (p *platform) setupScheduler() Platform {
	p.scheduler = scheduler.GetInstance()
	return p
}

// setupQueues initializes the RequestLaunch and Termination queues.
func (p *platform) setupQueues() Platform {
	p.LaunchQueue = queue.NewQueue()
	p.TerminationQueue = make(chan workers.TerminateInput, 1000)
	return p
}

// setupWorkers configures the RequestLaunch and the Termination Pool
// If there is an error during the PoolFactory execution, it returns an error.
func (p *platform) setupWorkers() (Platform, error) {
	var err error

	p.LaunchPool, err = p.PoolFactory(p.Config.PoolSizeLaunchSim, workers.Launch)
	if err != nil {
		return nil, err
	}

	p.TerminationPool, err = p.PoolFactory(p.Config.PoolSizeErrorHandler, workers.Terminate)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// setupTransport initializes Ignition transport.
func (p *platform) setupTransport() (Platform, error) {
	t, err := transport.New()
	if err != nil {
		return nil, err
	}
	p.transport = t
	return p, nil
}

func (p *platform) setupControllers() Platform {
	queueService := queue.NewService(p.LaunchQueue, p.Services().User())
	p.controllers = &controllers{
		queue: queue.NewController(queueService),
	}
	return p
}
