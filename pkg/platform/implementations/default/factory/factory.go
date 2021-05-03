package factory

import (
	email "gitlab.com/ignitionrobotics/web/cloudsim/pkg/email/implementations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	machines "gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations"
	orchestrator "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/runsim"
	secrets "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations"
	storage "gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations"
	store "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations"
)

// NewFunc is the factory creation function for the EC2 Machines implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Load components
	components := platform.Components{}
	factoryCalls := factory.Calls{
		// Machines
		{
			Factory:      machines.Factory,
			Config:       typeConfig.Components.Machines,
			Dependencies: dependencies,
			Out:          &components.Machines,
		},
		// Storage
		{
			Factory:      storage.Factory,
			Config:       typeConfig.Components.Storage,
			Dependencies: dependencies,
			Out:          &components.Storage,
		},
		// Orchestrator
		{
			Factory:      orchestrator.Factory,
			Config:       typeConfig.Components.Orchestrator,
			Dependencies: dependencies,
			Out:          &components.Cluster,
		},
		// Store
		{
			Factory:      store.Factory,
			Config:       typeConfig.Components.Store,
			Dependencies: dependencies,
			Out:          &components.Store,
		},
		// Secrets
		{
			Factory:      secrets.Factory,
			Config:       typeConfig.Components.Secrets,
			Dependencies: dependencies,
			Out:          &components.Secrets,
		},
		// Email Sender
		{
			Factory:      email.Factory,
			Config:       typeConfig.Components.EmailSender,
			Dependencies: dependencies,
			Out:          &components.EmailSender,
		},
	}
	if err := factory.CallFactories(factoryCalls); err != nil {
		return err
	}

	// Configure the RunningSimulations component
	components.RunningSimulations = runsim.NewManager()

	// Set output value
	platform, err := platform.NewPlatform(typeConfig.Name, components)
	if err != nil {
		return factory.ErrorWithContext(err)
	}
	if err := factory.SetValue(out, platform); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
