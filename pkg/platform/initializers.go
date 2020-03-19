package platform

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
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

func (p *Platform) initializeServer(config Config) *Platform {
	serverConfig := server.Config{
		Auth0:    config.Auth0,
		HTTPport: config.HTTPport,
		SSLport:  config.SSLport,
	}

	s, err := server.New(serverConfig)
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

func (p *Platform) initializeEmail() *Platform {
	p.Email = email.New()
	return p
}