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
func NewStore() (store.Store, error) {

	machines, err := newMachinesStore()
	if err != nil {
		return nil, err
	}

	ignition, err := newIgnitionStore()
	if err != nil {
		return nil, err
	}

	orchestrator, err := newOrchestratorStore()
	if err != nil {
		return nil, err
	}
	return &envStore{
		machines:     machines,
		ignition:     ignition,
		orchestrator: orchestrator,
	}, nil
}
