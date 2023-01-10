package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesNodes "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/nodes/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Nodes implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesNodes.NewFunc,
}
