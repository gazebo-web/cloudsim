package pods

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

// resource is a Pod implementation that contains the basic information to identify a resource in a Kubernetes cluster.
type resource struct {
	name      string
	selector  string
	namespace string
}

// Name returns the name of the resource.
func (p resource) Name() string {
	return p.name
}

// Selector returns the selector of the resource.
func (p resource) Selector() string {
	return p.selector
}

// Namespace returns the namespace of the resource.
func (p resource) Namespace() string {
	return p.namespace
}

// NewPod returns a new Resource implementation.
func NewPod(name string, namespace string, selector string) orchestrator.Resource {
	return &resource{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
