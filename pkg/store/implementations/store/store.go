package store

import storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"

// store is a store.Store implementation.
type store struct {
	machines     storepkg.Machines
	orchestrator storepkg.Orchestrator
	ignition     storepkg.Ignition
}

// Orchestrator returns a store.Orchestrator implementation that reads configuration from env vars.
func (e *store) Orchestrator() storepkg.Orchestrator {
	return e.orchestrator
}

// Ignition returns a store.Ignition implementation that reads configuration from env vars.
func (e *store) Ignition() storepkg.Ignition {
	return e.ignition
}

// Machines returns a store.Machines implementation that reads configuration from env vars.
func (e *store) Machines() storepkg.Machines {
	return e.machines
}

// NewStoreFromEnvVars initializes a new store.Store implementation using store.
func NewStoreFromEnvVars() storepkg.Store {
	return &store{
		machines:     newMachinesStoreFromEnvVars(),
		ignition:     newIgnitionStoreFromEnvVars(),
		orchestrator: newOrchestratorStoreFromEnvVars(),
	}
}
