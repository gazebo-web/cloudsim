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
	// HasIPStatusCondition is used to indicate that pods have ips available.
	HasIPStatusCondition = Condition{
		Type:   "HasIPStatus",
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
	NetworkPolicies() NetworkPolicies
}

// Selector is used to identify a certain resource.
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
	// Name returns the name of the resource
	Name() string
	// Selector returns the resource's Selector.
	Selector() Selector
	// Namespace returns the namespace where the resource lives in.
	Namespace() string
}
