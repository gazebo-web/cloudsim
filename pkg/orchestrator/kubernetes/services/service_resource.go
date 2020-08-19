package services

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

// resource is an orchestrator.Resource implementation of Kubernetes Services.
type resource struct {
	// name represents the name of the service.
	name string
	// selector defines a set of key-value pairs that identifies this resource.
	selector string
	// namespace is the environment where the resource is currently running.
	namespace string
}

// Name returns the resource's name.
func (s *resource) Name() string {
	return s.name
}

// Selector returns the resource's selector.
func (s *resource) Selector() string {
	return s.selector
}

// Namespace returns the resource's namespace.
func (s *resource) Namespace() string {
	return s.namespace
}

// NewResource initializes a new orchestrator.Resource using a kubernetes service implementation.
func NewResource(name, selector, namespace string) orchestrator.Resource {
	return &resource{
		name:      name,
		selector:  selector,
		namespace: namespace,
	}
}
