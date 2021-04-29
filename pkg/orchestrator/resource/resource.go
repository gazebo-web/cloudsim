package resource

// Resource groups a set of method to identify a resource in a cluster.
type Resource interface {
	// Name returns the name of the resource
	Name() string
	// Selector returns the resource's Selector.
	Selector() Selector
	// Namespace returns the namespace where the resource lives in.
	Namespace() string
}

// resource is a resource.Resource implementation of Kubernetes resources.
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

// NewResource initializes a new resource.Resource using a kubernetes service implementation.
func NewResource(name, namespace string, selector Selector) Resource {
	return &resource{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
