package nodes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// Resource is a Kubernetes resource. It extends the generic orchestrator.Resource interface.
type Resource interface {
	orchestrator.Resource
}

// resource is a Resource implementation that contains the basic information to identify a resource in a Kubernetes cluster.
type resource struct {
	name      string
	selector  string
	namespace string
}

// Name returns the name of the resource.
func (n *resource) Name() string {
	return n.name
}

// Selector returns the selector of the resource.
func (n *resource) Selector() string {
	return n.selector
}

// Namespace returns the namespace of the resource.
func (n *resource) Namespace() string {
	return n.namespace
}

// NewNodeResource returns a new Resource implementation using resource.
func NewNodeResource(name string, namespace string, selector string) Resource {
	return &resource{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
