package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	kubernetesConfigMap "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/configurations/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Configurations implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesConfigMap.NewFunc,
}
