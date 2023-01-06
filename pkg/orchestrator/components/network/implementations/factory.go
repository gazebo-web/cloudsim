package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesNetPols "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/network/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Policies implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesNetPols.NewFunc,
}
