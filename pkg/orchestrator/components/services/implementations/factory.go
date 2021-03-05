package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	kubernetesServices "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Services implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesServices.NewFunc,
}
