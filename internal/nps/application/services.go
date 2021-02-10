package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

type Services interface {
	application.Services
}

// services is a Services implementation.
type services struct {
	application.Services
}

// NewServices initializes a new Services implementation using a base generic service.
func NewServices(base application.Services) Services {
	return &services{
		Services: base,
	}
}
