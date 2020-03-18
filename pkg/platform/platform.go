package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transporter"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
)

type Platform struct {
	Server *ign.Server
	Logger ign.Logger
	Context context.Context
	Validator *validator.Validate
	Decoder *form.Decoder
	Transporter *transporter.Transporter
	UserAccessor *users.UserAccessor
}

func New(config Config) Platform {
	p := Platform{
		Server:       nil,
		Logger:       config.logger,
		Context:      nil,
		Validator:    nil,
		Decoder:      nil,
		Transporter:  nil,
		UserAccessor: nil,
	}

	p = initializeLogger(p)
	p = initializeContext(p)
	p = initializeServer(p, config)
	p.Logger.Info(fmt.Sprintf("Using HTTP port [%s] and SSL port [%s]", p.Server.HTTPPort, p.Server.SSLport))

	return p
}
