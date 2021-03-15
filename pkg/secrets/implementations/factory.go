package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	kubernetesfactory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations/kubernetes/factory"
)

const (
	// Kubernetes is the Kubernetes implementation factory identifier.
	Kubernetes = "kubernetes"
)

// Factory provides a factory to create Cluster implementations.
var Factory = factorymap.Map{
	Kubernetes: kubernetesfactory.NewFunc,
}
