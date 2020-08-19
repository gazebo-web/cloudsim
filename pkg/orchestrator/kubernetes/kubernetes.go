package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// k8s is a orchestrator.Cluster implementation.
type k8s struct {
	// nodes has a reference to an orchestrator.Nodes implementation.
	nodes orchestrator.Nodes

	// pods has a reference to an orchestrator.Pods implementation.
	pods orchestrator.Pods

	// ingressRules has a reference to an orchestrator.IngressRules implementation.
	ingressRules orchestrator.IngressRules

	// serviceManager has a reference to an orchestrator.Services implementation.
	serviceManager orchestrator.Services

	// ingresses has a reference to an orchestrator.Ingresses implementation.
	ingresses orchestrator.Ingresses
}

// IngressRules returns the Kubernetes orchestrator.IngressRules implementation.
func (k k8s) IngressRules() orchestrator.IngressRules {
	return k.ingressRules
}

// Nodes returns the Kubernetes orchestrator.Nodes implementation.
func (k k8s) Nodes() orchestrator.Nodes {
	return k.nodes
}

// Pods returns the Kubernetes orchestrator.Pods implementation.
func (k k8s) Pods() orchestrator.Pods {
	return k.pods
}

// Services returns the Kubernetes orchestrator.Services implementation.
func (k k8s) Services() orchestrator.Services {
	return k.serviceManager
}

// Ingresses returns the Kubernetes orchestrator.Ingresses implementation.
func (k k8s) Ingresses() orchestrator.Ingresses {
	return k.ingresses
}

// NewKubernetes returns a orchestrator.Cluster implementation using Kubernetes.
func NewKubernetes(nodes orchestrator.Nodes, pods orchestrator.Pods, ingresses orchestrator.Ingresses, ingressRules orchestrator.IngressRules) orchestrator.Cluster {
	return &k8s{
		nodes:        nodes,
		pods:         pods,
		ingresses:    ingresses,
		ingressRules: ingressRules,
	}
}
