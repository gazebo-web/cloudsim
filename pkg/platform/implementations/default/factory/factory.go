package factory

import (
	email "github.com/gazebo-web/cloudsim/v4/pkg/email/implementations"
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	machines "github.com/gazebo-web/cloudsim/v4/pkg/machines/implementations"
	orchestrator "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/implementations"
	"github.com/gazebo-web/cloudsim/v4/pkg/platform"
	secrets "github.com/gazebo-web/cloudsim/v4/pkg/secrets/implementations"
	storage "github.com/gazebo-web/cloudsim/v4/pkg/storage/implementations"
	store "github.com/gazebo-web/cloudsim/v4/pkg/store/implementations"
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
