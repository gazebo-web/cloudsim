package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// k8s is a orchestrator.Cluster implementation.
type k8s struct {
	// nodeManager has a reference to an orchestrator.Nodes implementation.
	nodeManager orchestrator.Nodes

	// podManager has a reference to an orchestrator.Pods implementation.
	podManager orchestrator.Pods

	// rulesManager has a reference to an orchestrator.IngressRules implementation.
	rulesManager orchestrator.IngressRules

	// serviceManager has a reference to an orchestrator.Services implementation.
	serviceManager orchestrator.Services

	// ingressManager has a reference to an orchestrator.Ingresses implementation.
	ingressManager orchestrator.Ingresses
}

// IngressRules returns the Kubernetes orchestrator.IngressRules implementation.
func (k k8s) IngressRules() orchestrator.IngressRules {
	return k.rulesManager
}

// Nodes returns the Kubernetes orchestrator.Nodes implementation.
func (k k8s) Nodes() orchestrator.Nodes {
	return k.nodeManager
}

// Pods returns the Kubernetes orchestrator.Pods implementation.
func (k k8s) Pods() orchestrator.Pods {
	return k.podManager
}

// Services returns the Kubernetes orchestrator.Services implementation.
func (k k8s) Services() orchestrator.Services {
	return k.serviceManager
}

// Ingresses returns the Kubernetes orchestrator.Ingresses implementation.
func (k k8s) Ingresses() orchestrator.Ingresses {
	return k.ingressManager
}

// NewKubernetes returns a orchestrator.Cluster implementation using Kubernetes.
func NewKubernetes(nodeManager orchestrator.Nodes, podManager orchestrator.Pods, rulesManager orchestrator.IngressRules) orchestrator.Cluster {
	return &k8s{
		nodeManager:  nodeManager,
		podManager:   podManager,
		rulesManager: rulesManager,
	}
}
