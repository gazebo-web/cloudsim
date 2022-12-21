package store

import (
	"github.com/gazebo-web/cloudsim/pkg/defaults"
	"github.com/gazebo-web/cloudsim/pkg/factory"
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
		ignition:     &typeConfig.IgnitionStore,
	}
	if err := factory.SetValue(out, store); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
