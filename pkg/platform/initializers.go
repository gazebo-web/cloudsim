package platform

import (
	"context"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
)

func (p *Platform) initializeLogger() *Platform {
	l, err := logger.New()
	if err != nil {
		log.Fatalf("Error parsing environment variables for Logger. %+v\n", err)
	}
	p.Logger = l
	return p
}

func (p *Platform) initializeContext() *Platform {
	ctx := ign.NewContextWithLogger(context.Background(), p.Logger)
	p.Context = ctx
	return p
}

func (p *Platform) initializeServer() *Platform {
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

func (p *Platform) initializeRouter() *Platform {
	cfg := router.Config{
		Version: "1.0",
	}
	r := router.New(cfg)
	p.Server.SetRouter(r)
	return p
}

func (p *Platform) initializeValidator() *Platform {
	validate := validator.New()
	p.Validator = validate
	return p
}

func (p *Platform) initializeFormDecoder() *Platform {
	p.FormDecoder = form.NewDecoder()
}

func (p *Platform) initializePermissions() *Platform {
	per := &permissions.Permissions{}
	per.Init(p.Server.Db, p.Config.SysAdmin)
	return p
}