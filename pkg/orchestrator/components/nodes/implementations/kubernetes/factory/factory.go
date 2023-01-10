package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	kubernetesNodes "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/nodes/implementations/kubernetes"
)

// NewFunc is the factory creation function for the Kubernetes nodes.Nodes implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create instance
	nodes := kubernetesNodes.NewNodes(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, nodes); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
