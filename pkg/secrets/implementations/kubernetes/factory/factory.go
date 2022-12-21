package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/implementations/kubernetes/client"
	kubernetesSecrets "github.com/gazebo-web/cloudsim/pkg/secrets/implementations/kubernetes"
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
		initializeAPI,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	secrets := kubernetesSecrets.NewKubernetesSecrets(typeDependencies.API.CoreV1())
	if err := factory.SetValue(out, secrets); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}

// initializeAPI initializes the API dependency.
func initializeAPI(config *Config, dependencies *Dependencies) error {
	if dependencies.API != nil {
		return nil
	}
	if config == nil {
		return factory.ErrorWithContext(factory.ErrNilConfig)
	}

	// Get the Kubernetes config
	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create the API
	api, err := client.NewAPI(kubeconfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	if err = factory.SetValue(&dependencies.API, api); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
