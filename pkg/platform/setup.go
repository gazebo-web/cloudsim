package platform

import (
	"context"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db"
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

// IPlatformSetup represent a set of methods to initialize the Platform.
type IPlatformSetup interface {
	setupLogger() *Platform
	setupContext() *Platform
	setupServer() *Platform
	setupRouter() *Platform
	setupEmail() *Platform
	setupValidator() *Platform
	setupFormDecoder() *Platform
	setupPermissions() *Platform
	setupUserService() *Platform
	setupDatabase() *Platform
	setupCloudProvider() *Platform
	setupOrchestrator() *Platform
	setupSimulator() *Platform
	setupPoolFactory() *Platform
	setupScheduler() *Platform
	setupQueues() *Platform
	setupWorkers() (*Platform, error)
	setupTransport() (*Platform, error)
}

// setupLogger initializes the logger.
func (p *Platform) setupLogger() *Platform {
	l, err := logger.New()
	if err != nil {
		log.Fatalf("Error parsing environment variables for Logger. %+v\n", err)
	}
	p.Logger = l
	return p
}

// setupContext initializes the context.
func (p *Platform) setupContext() *Platform {
	ctx := ign.NewContextWithLogger(context.Background(), p.Logger)
	p.Context = ctx
	return p
}

// setupServer initializes the HTTP server.
// If there is an error, it will panic.
func (p *Platform) setupServer() *Platform {
	cfg := server.Config{
		Auth0:    p.Config.Auth0,
		HTTPport: p.Config.HTTPport,
		SSLport:  p.Config.SSLport,
	}
	s, err := server.New(cfg)
	if err != nil {
		p.Logger.Critical(err)
		log.Fatalf("Error while initializing server. %v\n", err)
	}
	p.Server = s
	return p
}

// setupRouter initializes the server's router.
func (p *Platform) setupRouter() *Platform {
	r := router.New()
	p.Server.SetRouter(r)
	return p
}

// setupEmail initializes the email service.
func (p *Platform) setupEmail() *Platform {
	e := email.New()
	p.Email = e
	return p
}

// setupValidator initializes the validator.
func (p *Platform) setupValidator() *Platform {
	validate := validator.New()
	p.Validator = validate
	return p
}

// setupFormDecoder initializes the form decoder.
func (p *Platform) setupFormDecoder() *Platform {
	p.FormDecoder = form.NewDecoder()
	return p
}

// setupPermissions initializes the platform permissions.
func (p *Platform) setupPermissions() *Platform {
	per := &permissions.Permissions{}
	err := per.Init(p.Server.Db, p.Config.SysAdmin)
	if err != nil {
		p.Logger.Critical(err)
		log.Fatalf("Error while initializing server. %v\n", err)
	}
	p.Permissions = per
	return p
}

// setupUserService initializes the User Service.
func (p *Platform) setupUserService() *Platform {
	s, err := users.NewService(p.Permissions, p.Config.SysAdmin)
	if err != nil {
		p.Logger.Critical(err)
		log.Fatalf("Error while configuring user service. %v\n", err)
	}
	p.UserService = s
	return p
}

// setupDatabase performs migrations, adds default data and adds custom indexes.
func (p *Platform) setupDatabase() *Platform {
	db.Migrate(p.Context, p.Server.Db)
	db.AddDefaultData(p.Context, p.Server.Db)
	db.AddCustomIndexes(p.Context, p.Server.Db)
	return p
}

// setupCloudProvider initializes the Cloud Provider.
func (p *Platform) setupCloudProvider() *Platform {
	p.CloudProvider = cloud.New()
	return p
}

// setupOrchestrator initializes the container Orchestrator.
func (p *Platform) setupOrchestrator() *Platform {
	p.Orchestrator = orchestrator.New()
	return p
}

// setupSimulator initializes the Simulator.
func (p *Platform) setupSimulator() *Platform {
	input := simulator.NewSimulatorInput{}
	p.Simulator = simulator.NewSimulator(input)
	return p
}

// setupPoolFactory initializes the Default Pool Factory.
func (p *Platform) setupPoolFactory() *Platform {
	p.PoolFactory = pool.NewPool
	return p
}

// setupScheduler gets the instance from the scheduler package.
func (p *Platform) setupScheduler() *Platform {
	p.Scheduler = scheduler.GetInstance()
	return p
}

// setupQueues initializes the RequestLaunch and Termination queues.
func (p *Platform) setupQueues() *Platform {
	p.LaunchQueue = queue.NewQueue()
	p.TerminationQueue = make(chan workers.TerminateInput, 1000)
	return p
}

// setupWorkers configures the RequestLaunch and the Termination Pool
// If there is an error during the PoolFactory execution, it returns an error.
func (p *Platform) setupWorkers() (*Platform, error) {
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

// setupTransport initializes Ignition Transport.
func (p *Platform) setupTransport() (*Platform, error) {
	t, err := transport.New()
	if err != nil {
		return nil, err
	}
	p.Transport = t
	return p, nil
}
