package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// Services groups the services needed by an application to launch simulations.
type Services interface {
	// Simulations returns a service to operate over different simulations.
	Simulations() simulations.Service

	// Store returns a service to handle configuration.
	Store() store.Store
}

// services is a Services implementation.
type services struct {
	simulation simulations.Service
	store      store.Store
}

// Store returns the underlying Config Store's service.
func (s *services) Store() store.Store {
	return s.store
}

// Simulations returns the underlying Simulation's service.
func (s *services) Simulations() simulations.Service {
	return s.simulation
}

// NewServices initializes a new Application Services implementation.
func NewServices(simulation simulations.Service, store store.Store) Services {
	return &services{
		simulation: simulation,
		store:      store,
	}
}
