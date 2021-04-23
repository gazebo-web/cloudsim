package orchestrator

import "time"

// Resource groups a set of method to identify a resource in a cluster.
type Resource interface {
	// Name returns the name of the Resource
	Name() string
	// Namespace returns the namespace where the Resource lives in.
	Namespace() string
	// Selector returns the Resource's Selector.
	Selector() Selector
}

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

//--------------------------------------------------------------------------------------

// ResourcePhase has a method to return the phase of a certain Resource.
type ResourcePhase interface {
	// Phase is a simple, high-level summary of where the Resource is in its lifecycle.
	Phase() Phase
}

type resourcePhase struct {
	// phase is a simple, high-level summary of where the Resource is in its lifecycle.
	phase Phase
}

// Phase is a simple, high-level summary of where the ResourcePhase is in its lifecycle.
func (r *resourcePhase) Phase() Phase {
	return r.phase
}

// NewResourcePhase initializes a new ResourcePhase implementation.
func NewResourcePhase(phase Phase) ResourcePhase {
	return &resourcePhase{phase: phase}
}

//--------------------------------------------------------------------------------------

// ResourceTimestamp has some methods to expose the creation and deletion timestamps of a Resource.
type ResourceTimestamp interface {
	// CreationTimestamp is a timestamp representing the server time when this object was created.
	CreationTimestamp() time.Time
	// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
	// graceful deletion is requested.
	DeletionTimestamp() *time.Time
}

type resourceTimestamp struct {
	// creationTimestamp is the ResourceTimestamp.CreationTimestamp.
	creationTimestamp time.Time
	// deletionTimestamp is the ResourceTimestamp.DeletionTimestamp.
	deletionTimestamp *time.Time
}

// CreationTimestamp is a timestamp representing the server time when this object was created.
func (s *resourceTimestamp) CreationTimestamp() time.Time {
	return s.creationTimestamp
}

// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
// graceful deletion is requested.
func (s *resourceTimestamp) DeletionTimestamp() *time.Time {
	return s.deletionTimestamp
}

// NewResourceTimestamp is used to initialize a new ResourceTimestamp implementation.
func NewResourceTimestamp(creation time.Time, deletion *time.Time) ResourceTimestamp {
	return &resourceTimestamp{
		creationTimestamp: creation,
		deletionTimestamp: deletion,
	}
}
