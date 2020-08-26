package orchestrator

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
)

// Cluster groups a set of methods for managing a cluster.
type Cluster interface {
	Nodes() Nodes
	Pods() Pods
	Services() Services
	Ingresses() Ingresses
	IngressRules() IngressRules
}

// Selector is used to identify a certain resource.
type Selector interface {
	// String returns the selector represented in string format.
	String() string
}

// Resource groups a set of method to identify a resource in a cluster.
type Resource interface {
	// Name returns the name of the resource
	Name() string
	// Selector returns the resource's Selector.
	Selector() Selector
	// Namespace returns the namespace where the resource lives in.
	Namespace() string
}