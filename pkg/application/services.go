package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Services groups the services needed by an application to launch simulations.
type Services interface {
	// Simulations returns a service to operate over different simulations.
	Simulations() simulations.Service
}

// services is a Services implementation.
type services struct {
	simulation simulations.Service
}

// Simulations returns the underlying Simulation's service.
func (s *services) Simulations() simulations.Service {
	return s.simulation
}

// NewServices initializes a new Application Services implementation.
func NewServices(simulation simulations.Service) Services {
	return &services{
		simulation: simulation,
	}
}
