package env

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"

// envStore is a store.Store implementation.
type envStore struct {
	machines     store.Machines
	orchestrator store.Orchestrator
	ignition     store.Ignition
}

// Orchestrator returns a store.Orchestrator implementation that reads configuration from env vars.
func (e *envStore) Orchestrator() store.Orchestrator {
	return e.orchestrator
}

// Ignition returns a store.Ignition implementation that reads configuration from env vars.
func (e *envStore) Ignition() store.Ignition {
	return e.ignition
}

// Machines returns a store.Machines implementation that reads configuration from env vars.
func (e *envStore) Machines() store.Machines {
	return e.machines
}

// NewStore initializes a new store.Store implementation using envStore.
func NewStore() store.Store {
	return &envStore{
		machines:     newMachinesStore(),
		ignition:     newIgnitionStore(),
		orchestrator: newOrchestratorStore(),
	}
}
