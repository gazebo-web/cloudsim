package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// Services groups the services needed by an application to launch simulations.
type Services interface {
	// Simulations returns a service to operate over different simulations.
	Simulations() simulations.Service

	// ConfigStore returns a service to handle configuration.
	ConfigStore() store.ConfigStore
}

// services is a Services implementation.
type services struct {
	simulation  simulations.Service
	configStore store.ConfigStore
}

// ConfigStore returns the underlying Config Store's service.
func (s *services) ConfigStore() store.ConfigStore {
	return s.configStore
}

// Simulations returns the underlying Simulation's service.
func (s *services) Simulations() simulations.Service {
	return s.simulation
}

// NewServices initializes a new Application Services implementation.
func NewServices(simulation simulations.Service, configStore store.ConfigStore) Services {
	return &services{
		simulation:  simulation,
		configStore: configStore,
	}
}
