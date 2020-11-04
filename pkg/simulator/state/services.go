package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services exposes a method to access the application services.
type Services interface {
	Services() application.Services
}
