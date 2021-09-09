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
		return factory.ErrorWithContext(err)
	}
	if err := defaults.SetValues(&typeConfig); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create instance
	store := &store{
		machines:     &typeConfig.MachinesStore,
		orchestrator: &typeConfig.OrchestratorStore,
		mole:         &typeConfig.MoleStore,
		ignition:     &typeConfig.IgnitionStore,
	}
	if err := factory.SetValue(out, store); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
