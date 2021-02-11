package application

// This file defines the services used by this application. A service is <TODO: WHAT IS A SERVICE>??

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services groups the services needed by an application to launch simulations.
// The default set of services can be found in `pkg/application/services.go`.
type Services interface {
	application.Services
}

// services is a Services implementation.
type services struct {
	application.Services
}

// NewServices initializes a new services implementation using a base generic service.
func NewServices(base application.Services) Services {
	return &services{
		Services: base,
	}
}
