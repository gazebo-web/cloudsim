package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	kubernetesConfigMaps "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations/implementations/kubernetes"
)

// NewFunc is the factory creation function for the Kubernetes configurations.Configurations implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create instance
	netPols := kubernetesConfigMaps.NewConfigMaps(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, netPols); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
