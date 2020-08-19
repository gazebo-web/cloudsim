package ingress

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

// Resource is a Kubernetes ingress. It extends the generic orchestrator.Resource interface.
type Resource interface {
	orchestrator.Resource
}

// ingress is an Resource implementation.
type ingress struct {
	// name is the ingress' name.
	name string
	// namespace is the ingress' namespace.
	namespace string
}

// Selector returns the ingress selector. We aren't using selectors for ingresses right now.
func (i *ingress) Selector() string {
	return ""
}

// Namespace returns the ingress namespace.
func (i *ingress) Namespace() string {
	return i.namespace
}

// Name returns the ingress name.
func (i *ingress) Name() string {
	return i.name
}

// NewIngress initializes a new Resource using Kubernetes.
func NewIngress(name string, namespace string) Resource {
	return &ingress{
		name:      name,
		namespace: namespace,
	}
}
