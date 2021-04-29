package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes/client"
)

// NewFunc is the factory creation function for the Kubernetes pods.Pods implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Initialize dependencies
	dependenciesInitFns := []func(config *Config, dependencies *Dependencies) error{
		initializeSPDY,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	pods := kubernetesPods.NewPods(typeDependencies.API, typeDependencies.SPDY, typeDependencies.Logger)
	if err := factory.SetValue(out, pods); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}

// initializeSPDY initializes the SPDY initializer dependency.
func initializeSPDY(config *Config, dependencies *Dependencies) error {
	// Validate inputs
	if dependencies.SPDY != nil {
		return nil
	}

	// Get the Kubernetes config
	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create the SPDY Initializer
	spdy := spdy.NewSPDYInitializer(kubeconfig)

	if err = factory.SetValue(&dependencies.SPDY, spdy); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
