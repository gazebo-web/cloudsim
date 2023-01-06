package application

import (
	"github.com/gazebo-web/cloudsim/pkg/billing"
	"github.com/gazebo-web/cloudsim/pkg/simulations"
	"github.com/gazebo-web/cloudsim/pkg/users"
)

// Services groups the services needed by an application to launch simulations.
type Services interface {
	// Simulations provides access to a different set of methods for managing simulations.
	Simulations() simulations.Service

	// Users provides access to a different set of methods for managing users.
	Users() users.Service

	// Billing provides access to a different set of methods for managing credits.
	Billing() billing.Service
}

// services is a Services implementation.
type services struct {
	simulation simulations.Service
	user       users.Service
	billing    billing.Service
}

// Billing returns the underlying Billing service.
func (s *services) Billing() billing.Service {
	return s.billing
}

// Users returns the underlying User service.
func (s *services) Users() users.Service {
	return s.user
}

// Simulations returns the underlying Simulation service.
func (s *services) Simulations() simulations.Service {
	return s.simulation
}

// NewServices initializes a new Application Services implementation.
func NewServices(simulation simulations.Service, user users.Service, billing billing.Service) Services {
	return &services{
		simulation: simulation,
		user:       user,
		billing:    billing,
	}
}
