package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	kubernetesNodes "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations/kubernetes"
)

// NewFunc is the factory creation function for the Kubernetes nodes.Nodes implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Create instance
	nodes := kubernetesNodes.NewNodes(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, nodes); err != nil {
		return err
	}

	return nil
}
