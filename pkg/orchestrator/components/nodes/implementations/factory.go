package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	kubernetesNodes "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Nodes implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesNodes.NewFunc,
}
