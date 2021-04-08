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
		return err
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Load components
	components := platform.Components{}
	factoryCalls := factory.Calls{
		// Machines
		{
			Factory:      machines.Factory,
			Config:       typeConfig.Machines,
			Dependencies: dependencies,
			Out:          &components.Machines,
		},
		// Storage
		{
			Factory:      storage.Factory,
			Config:       typeConfig.Storage,
			Dependencies: dependencies,
			Out:          &components.Storage,
		},
		// Orchestrator
		{
			Factory:      orchestrator.Factory,
			Config:       typeConfig.Orchestrator,
			Dependencies: dependencies,
			Out:          &components.Cluster,
		},
		// Store
		{
			Factory:      store.Factory,
			Config:       typeConfig.Store,
			Dependencies: dependencies,
			Out:          &components.Store,
		},
		// Secrets
		{
			Factory:      secrets.Factory,
			Config:       typeConfig.Secrets,
			Dependencies: dependencies,
			Out:          &components.Secrets,
		},
		// Email Sender
		{
			Factory:      email.Factory,
			Config:       typeConfig.EmailSender,
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
	platform := platform.NewPlatform(components)
	if err := factory.SetValue(out, platform); err != nil {
		return err
	}

	return nil
}
