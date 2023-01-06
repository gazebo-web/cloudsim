package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesPods "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Pods implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesPods.NewFunc,
}
