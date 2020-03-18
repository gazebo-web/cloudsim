package platform

import (
	"context"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transporter"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Platform struct {
	Server *ign.Server
	Logger *ign.Logger
	Context *context.Context
	Validator *validator.Validate
	Decoder *form.Decoder
	Transporter *transporter.Transporter
	UserAccessor *users.UserAccessor

}

func New(config Config) Platform {

}