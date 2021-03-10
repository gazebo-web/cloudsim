package store

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
)

// NewFunc is the factory creation function for the Kubernetes orchestrator.Cluster implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return err
	}
	if err := defaults.SetDefaults(&typeConfig); err != nil {
		return err
	}

	// Create instance
	store := &store{
		machines:     &typeConfig.MachinesStore,
		orchestrator: &typeConfig.OrchestratorStore,
		ignition:     &typeConfig.IgnitionStore,
	}
	if err := factory.SetValue(out, store); err != nil {
		return err
	}

	return nil
}
