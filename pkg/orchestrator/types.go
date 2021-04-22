package orchestrator

import (
	"fmt"
	"strings"
	"time"
)

// Condition represents a state that should be reached.
type Condition struct {
	Type   string
	Status string
}

var (
	// ReadyCondition is used to indicate that Nodes and Pods are ready.
	ReadyCondition = Condition{
		Type:   "Ready",
		Status: "True",
	}
	// HasIPStatusCondition is used to indicate that pods have ips available.
	HasIPStatusCondition = Condition{
		Type:   "HasIPStatus",
		Status: "True",
	}
)

// Phase represents a certain point in the lifecycle of a Resource.
type Phase string

const (
	// PhasePending is used to represent when a Resource is on a Pending Phase.
	// Used by: Pods, Nodes.
	PhasePending Phase = "Pending"
	// PhaseRunning is used to represent when a Resource is on a Running Phase.
	// Used by Pods, Nodes.
	PhaseRunning Phase = "Running"
	// PhaseSucceeded is used to represent when a Resource is on a Succeeded Phase.
	// Used by: Pods.
	PhaseSucceeded Phase = "Succeeded"
	// PhaseFailed is used to represent when a Resource is on a Failed Phase.
	// Used by: Pods.
	PhaseFailed Phase = "Failed"
	// PhaseUnknown is used to represent when a Resource is on a Unknown Phase.
	// Used by: Pods.
	PhaseUnknown Phase = "Unknown"
	// PhaseTerminated is used to represent when a Resource is on a Terminated Phase.
	// Used by: Nodes.
	PhaseTerminated Phase = "Terminated"
)

// Selector is used to represent the state a certain resource.
type Selector interface {
	// String returns the selector represented in string format.
	String() string
	// Map returns the underlying selector's map.
	Map() map[string]string
	// Extend extends the underlying base map with the extension selector.
	// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
	Extend(extension Selector) Selector
	// Set sets the given value to the given key. If the key already exists, it will be overwritten.
	Set(key string, value string)
}

// Resource groups a set of method to identify a resource in a cluster.
type Resource interface {
	// Name returns the name of the Resource
	Name() string
	// Selector returns the Resource's Selector.
	Selector() Selector
	// Namespace returns the namespace where the Resource lives in.
	Namespace() string
	// CreationTimestamp is a timestamp representing the server time when this object was created.
	CreationTimestamp() time.Time
	// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
	// graceful deletion is requested.
	DeletionTimestamp() *time.Time
	// Phase is a simple, high-level summary of where the Resource is in its lifecycle.
	Phase() Phase
}

// resource is an orchestrator.Resource implementation of Kubernetes resources.
type resource struct {
	// name represents the name of the service.
	name string
	// selector defines a set of key-value pairs that identifies this resource.
	selector Selector
	// namespace is the environment where the resource is currently running.
	namespace string
	// phase is a simple, high-level summary of where the Resource is in its lifecycle.
	phase Phase
	// creationTimestamp is a timestamp representing the server time when this object was created.
	creationTimestamp time.Time
	// deletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
	// graceful deletion is requested.
	deletionTimestamp *time.Time
}

// CreationTimestamp is a timestamp representing the server time when this object was created.
func (s *resource) CreationTimestamp() time.Time {
	return s.creationTimestamp
}

// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
// graceful deletion is requested.
func (s *resource) DeletionTimestamp() *time.Time {
	return s.deletionTimestamp
}

// Phase is a simple, high-level summary of where the Resource is in its lifecycle.
func (s *resource) Phase() Phase {
	return s.phase
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

// ResourceOptions is used to set the different values for a Resource.
type ResourceOptions struct {
	// Name is the Resource.Name.
	Name string
	// Namespace is the Resource.Namespace.
	Namespace string
	// Selector is the Resource.Selector.
	Selector Selector
	// Phase is the Resource.Phase.
	Phase Phase
	// CreationTimestamp is the Resource.CreationTimestamp.
	CreationTimestamp time.Time
	// DeletionTimestamp is the Resource.DeletionTimestamp.
	DeletionTimestamp *time.Time
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
