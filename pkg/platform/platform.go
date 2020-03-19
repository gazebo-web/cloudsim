package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
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
	Email email.Email
}

func New(config Config) Platform {
	p := Platform{}

	p.initializeLogger()
	p.initializeContext()
	p.initializeServer(config)
	p.initializeRouter()
	p.initializeEmail()
	p.Logger.Info(fmt.Sprintf("Using HTTP port [%s] and SSL port [%s]", p.Server.HTTPPort, p.Server.SSLport))
	return p
}
