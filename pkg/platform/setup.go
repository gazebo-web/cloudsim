package platform

import (
	"context"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
)

type ISetup interface {
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
}

func (p *Platform) setupLogger() *Platform {
	l, err := logger.New()
	if err != nil {
		log.Fatalf("Error parsing environment variables for Logger. %+v\n", err)
	}
	p.Logger = l
	return p
}

func (p *Platform) setupContext() *Platform {
	ctx := ign.NewContextWithLogger(context.Background(), p.Logger)
	p.Context = ctx
	return p
}

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

func (p *Platform) setupRouter() *Platform {
	cfg := router.Config{
		Version: "1.0",
	}
	r := router.New(cfg)
	p.Server.SetRouter(r)
	return p
}

func (p *Platform) setupValidator() *Platform {
	validate := validator.New()
	p.Validator = validate
	return p
}

func (p *Platform) setupFormDecoder() *Platform {
	p.FormDecoder = form.NewDecoder()
	return p
}

func (p *Platform) setupPermissions() *Platform {
	per := &permissions.Permissions{}
	err := per.Init(p.Server.Db, p.Config.SysAdmin)
	if err != nil {
		// TODO: Throw error
	}
	p.Permissions = per
	return p
}

func (p *Platform) setupUserService() *Platform {
	s, err := users.NewService(p.Permissions, p.Config.SysAdmin)
	if err != nil {
		// TODO: Throw error
	}
	p.UserService = s
	return p
}

func (p *Platform) setupDatabase() *Platform {
	db.Migrate(p.Context, p.Server.Db)
	db.AddDefaultData(p.Context, p.Server.Db)
	db.AddCustomIndexes(p.Context, p.Server.Db)
	return p
}

func (p *Platform) setupCloudProvider() *Platform {
	p.CloudProvider = cloud.New()
	return p
}

func (p *Platform) setupOrchestrator() *Platform {
	p.Orchestrator = orchestrator.New()
	return p
}