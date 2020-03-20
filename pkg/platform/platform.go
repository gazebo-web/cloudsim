package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transporter"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Platform struct {
	Server *ign.Server
	Logger ign.Logger
	Context context.Context
	Validator *validator.Validate
	FormDecoder *form.Decoder
	Transporter *transporter.Transporter
	UserAccessor *users.UserAccessor
	Orchestrator *orchestrator.Kubernetes
	CloudProvider *cloud.AmazonWS
	Config Config
}

func New(config Config) Platform {
	p := Platform{}
	p.Config = config

	p.initializeLogger()
	p.Logger.Debug("[INIT] Logger initialized.")

	// TODO: Decide where the score generation should go

	p.initializeContext()
	p.Logger.Debug("[INIT] Context initialized.")

	p.initializeServer()
	p.Logger.Debug(fmt.Sprintf("[INIT] Server initialized using HTTP port [%s] and SSL port [%s].", p.Server.HTTPPort, p.Server.SSLport))
	p.Logger.Debug(fmt.Sprintf("[INIT] Database [%s] initialized", p.Server.DbConfig.Name))

	p.initializeRouter()
	p.Logger.Debug("[INIT] Router initialized.")

	p.initializeValidator() // TODO: Decide where should the custom validators should go
	p.Logger.Debug("[INIT] Validators initialized.")

	p.initializeFormDecoder()
	p.Logger.Debug("[INIT] Form decoder initialized.")

	p.initializePermissions()
	p.Logger.Debug("[INIT] Permissions initialized.")

	p.initializeDatabase()
	p.Logger.Debug("[INIT] Database initialized: Migration, default data and custom indexes.")

	p.initializeCloudProvider()
	p.initializeOrchestrator()

	return p
}
