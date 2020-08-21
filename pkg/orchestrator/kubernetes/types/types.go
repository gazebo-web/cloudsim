package types

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// resource is an orchestrator.Resource implementation of Kubernetes resources.
type resource struct {
	// name represents the name of the service.
	name string
	// selector defines a set of key-value pairs that identifies this resource.
	selector orchestrator.Selector
	// namespace is the environment where the resource is currently running.
	namespace string
}

// Name returns the resource's name.
func (s *resource) Name() string {
	return s.name
}

// Selector returns the resource's selector.
func (s *resource) Selector() orchestrator.Selector {
	return s.selector
}

// Namespace returns the resource's namespace.
func (s *resource) Namespace() string {
	return s.namespace
}

// NewResource initializes a new orchestrator.Resource using a kubernetes service implementation.
func NewResource(name, namespace string, selector orchestrator.Selector) orchestrator.Resource {
	return &resource{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}

// selector is a kubernetes resource selector.
type selector map[string]string

// String returns the kubernetes selector in string format.
func (s selector) String() string {
	var out string
	count := len(s)
	for key, value := range s {
		out += fmt.Sprintf("%s=%s", key, value)
		count--
		if count > 1 {
			out += ","
		}
	}
	return out
}

// NewSelector initializes a new orchestrator.Selector from the given map.
// If `nil` is passed as input, an empty selector will be returned.
func NewSelector(input map[string]string) orchestrator.Selector {
	var output selector
	if input == nil {
		input = map[string]string{}
	}
	output = input
	return &output
}
