package orchestrator

type NetworkIngressRule struct {
	Ports    []int32
	IPBlocks []string
}

type NetworkEgress struct {
	Ports    []int32
	IPBlocks []string
}

// CreateNetworkPolicyInput is the input for creating a new network policy.
// TODO: Make this struct more generic.
type CreateNetworkPolicyInput struct {
	// Name is the name of the network policy.
	Name string
	// Namespace is the namespace where this network policy will be created.
	Namespace string
	// Labels is the group of key-value pairs that will identify this policy.
	Labels map[string]string
	// PodSelector are the labels of the pods that this policy should match to.
	PodSelector Selector
	// CIDR is the range of ips that this network policy will be applied to.
	CIDR string
	// WebsocketPort is the port where the websocket server will listen from.
	WebsocketPort int32
	// PeersFrom is the group of pod selectors that are allowed to access the pods covered by this network policy.
	PeersFrom []Selector
	// PeersTo is the group of pod selectors that the pods covered by this network policy are allowed to access to.
	PeersTo []Selector

	Ingresses NetworkIngressRule

	Egresses NetworkIngressRule
}

// NetworkPolicies groups a set of methods to manage network policies.
type NetworkPolicies interface {
	// Create creates a new network policy.
	Create(input CreateNetworkPolicyInput) (Resource, error)
}
