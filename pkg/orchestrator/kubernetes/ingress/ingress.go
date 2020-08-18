package ingress

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

// Ingress is a Kubernetes ingress. It extends the generic orchestrator.Resource interface.
type Ingress interface {
	orchestrator.Resource
}

// ingress is an Ingress implementation.
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

// NewIngress initializes a new Ingress using Kubernetes.
func NewIngress(name string, namespace string) Ingress {
	return &ingress{
		name:      name,
		namespace: namespace,
	}
}
