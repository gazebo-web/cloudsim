package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	kubernetesConfigMap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Configurations implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesConfigMap.NewFunc,
}
