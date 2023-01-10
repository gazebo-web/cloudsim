package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesfactory "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Cluster implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesfactory.NewFunc,
}
