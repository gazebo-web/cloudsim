package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// k8s is a orchestrator.ClusterManager implementation.
type k8s struct {
	// nodeManager has a reference to a orchestrator.NodeManager implementation.
	nodeManager orchestrator.NodeManager

	// podManager has a reference to a orchestrator.PodManager implementation.
	podManager orchestrator.PodManager
}

// Nodes returns the Kubernetes orchestrator.NodeManager implementation.
func (k k8s) Nodes() orchestrator.NodeManager {
	return k.nodeManager
}

// Pods returns the Kubernetes orchestrator.PodManager implementation.
func (k k8s) Pods() orchestrator.PodManager {
	return k.podManager
}

// Services returns the Kubernetes orchestrator.ServiceManager implementation.
func (k k8s) Services() orchestrator.ServiceManager {
	panic("implement me")
}

// Ingresses returns the Kubernetes orchestrator.IngressManager implementation.
func (k k8s) Ingresses() orchestrator.IngressManager {
	panic("implement me")
}

// NewKubernetes returns a orchestrator.ClusterManager implementation using Kubernetes.
func NewKubernetes(nodeManager orchestrator.NodeManager, podManager orchestrator.PodManager) orchestrator.ClusterManager {
	return &k8s{
		nodeManager: nodeManager,
		podManager:  podManager,
	}
}
