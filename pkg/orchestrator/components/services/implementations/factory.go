package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesServices "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/services/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Services implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesServices.NewFunc,
}
