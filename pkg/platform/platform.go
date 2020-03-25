package platform

import (
	"context"
	"github.com/go-playground/form"
	"github.com/go-playground/validator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/handlers"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/manager"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transporter"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IPlatform interface {
	Name() string
}

type Platform struct {
	Server *ign.Server
	Logger ign.Logger
	Context context.Context
	Validator *validator.Validate
	FormDecoder *form.Decoder
	Transporter *transporter.Transporter
	Orchestrator *orchestrator.Kubernetes
	CloudProvider *cloud.AmazonWS
	Permissions *permissions.Permissions
	UserService *users.Service
	Config Config
	HTTPHandlers *handlers.HTTPHandler
	Manager *manager.Manager
	Applications []*application.IApplication
}

// Name returns the platform name
func (p Platform) Name() string {
	return "cloudsim"
}