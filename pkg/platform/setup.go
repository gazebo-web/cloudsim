package platform

import (
	"context"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/manager"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/pool"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
)

// IPlatformSetup represent a set of methods to initialize the Platform.
type IPlatformSetup interface {
	setupLogger() *Platform
	setupContext() *Platform
	setupServer() *Platform
	setupRouter() *Platform
	setupValidator() *Platform
	setupFormDecoder() *Platform
	setupPermissions() *Platform
	setupUserService() *Platform
	setupDatabase() *Platform
	setupCloudProvider() *Platform
	setupOrchestrator() *Platform
	setupManager() *Platform
	setupPoolFactory() *Platform
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
	cfg := router.Config{
		Version: "1.0",
	}
	r := router.New(cfg)
	p.Server.SetRouter(r)
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

func (p *Platform) setupManager() *Platform {
	p.Manager = manager.New(p.Orchestrator, p.CloudProvider)
	return p
}

func (p *Platform) setupPoolFactory() *Platform {
	p.PoolFactory = pool.DefaultFactory
	return p
}