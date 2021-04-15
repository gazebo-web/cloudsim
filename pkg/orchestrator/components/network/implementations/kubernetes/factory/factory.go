package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	kubernetesNetpols "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations/kubernetes"
)

// NewFunc is the factory creation function for the Kubernetes network.Policies implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create instance
	netPols := kubernetesNetpols.NewNetworkPolicies(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, netPols); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}
