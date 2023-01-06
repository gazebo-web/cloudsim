package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	kubernetesServices "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/services/implementations/kubernetes"
)

// NewFunc is the factory creation function for the Kubernetes services.Services implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create instance
	services := kubernetesServices.NewServices(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, services); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
