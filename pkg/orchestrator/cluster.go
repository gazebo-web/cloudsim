package orchestrator

// Cluster groups a set of methods for managing a cluster.
type Cluster interface {
	Nodes() Nodes
	Pods() Pods
	Services() Services
	Ingresses() Ingresses
	IngressRules() IngressRules
	NetworkPolicies() NetworkPolicies
}
