package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	kubernetesIngresses "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/kubernetes"
	kubernetesIngressRules "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/kubernetes/rules"
)

// IngressesNewFunc is the factory creation function for the Kubernetes ingresses.Ingresses implementation.
func IngressesNewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Create instance
	ingresses := kubernetesIngresses.NewIngresses(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, ingresses); err != nil {
		return err
	}

	return nil
}

// IngressRulesNewFunc is the factory creation function for the Kubernetes ingresses.IngressRules implementation.
func IngressRulesNewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Create instance
	ingressRules := kubernetesIngressRules.NewIngressRules(typeDependencies.API, typeDependencies.Logger)
	if err := factory.SetValue(out, ingressRules); err != nil {
		return err
	}

	return nil
}
