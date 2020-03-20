package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
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
	Config Config
}

func New(config Config) Platform {
	p := Platform{}
	p.Config = config
	p.initializeLogger()
	// TODO: Decide where the score generation should go
	p.initializeContext()
	p.initializeServer()
	p.initializeRouter()
	p.initializeValidator() // TODO: Decide where should the custom validators should go
	p.initializeFormDecoder()
	p.initializePermissions()

	p.Logger.Info(fmt.Sprintf("Using HTTP port [%s] and SSL port [%s]", p.Server.HTTPPort, p.Server.SSLport))
	return p
}
