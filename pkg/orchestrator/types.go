package orchestrator

import (
	"fmt"
	"strings"
)

// resource is an orchestrator.Resource implementation of Kubernetes resources.
type resource struct {
	// name represents the name of the service.
	name string
	// selector defines a set of key-value pairs that identifies this resource.
	selector Selector
	// namespace is the environment where the resource is currently running.
	namespace string
}

// Name returns the resource's name.
func (s *resource) Name() string {
	return s.name
}

// Selector returns the resource's selector.
func (s *resource) Selector() Selector {
	return s.selector
}

// Namespace returns the resource's namespace.
func (s *resource) Namespace() string {
	return s.namespace
}

// NewResource initializes a new orchestrator.Resource using a kubernetes service implementation.
func NewResource(name, namespace string, selector Selector) Resource {
	return &resource{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}

// selector is a group of key-pair values that identify a resource.
type selector map[string]string

// Set sets the given value to the given key. If the key already exists, it will be overwritten.
func (s selector) Set(key string, value string) {
	s[key] = value
}

// Extend extends the underlying base map with the extension selector.
// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
func (s selector) Extend(extension Selector) Selector {
	for k, v := range extension.Map() {
		s[k] = v
	}
	return s
}

// Map returns the selector in map format.
func (s selector) Map() map[string]string {
	return s
}

// String returns the selector in string format.
func (s selector) String() string {
	var out string
	var labels []string
	for key, value := range s {
		out = fmt.Sprintf("%s=%s", key, value)
		labels = append(labels, out)
	}
	return strings.Join(labels, ",")
}

// NewSelector initializes a new orchestrator.Selector from the given map.
// If `nil` is passed as input, an empty selector will be returned.
func NewSelector(input map[string]string) Selector {
	if input == nil {
		input = map[string]string{}
	}
	return selector(input)
}
