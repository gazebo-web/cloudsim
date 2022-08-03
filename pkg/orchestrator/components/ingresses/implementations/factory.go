package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	kubernetesIngresses "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// IngressesFactory provides a factory to create Ingresses implementations.
var IngressesFactory = factorymap.Map{
	Kubernetes: kubernetesIngresses.IngressesNewFunc,
}

// IngressRulesFactory provides a factory to create IngressRules implementations.
var IngressRulesFactory = factorymap.Map{
	Kubernetes: kubernetesIngresses.IngressRulesNewFunc,
}
