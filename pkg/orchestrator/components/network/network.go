package network

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
)

// IngressRule groups a set of rules to be applied on a certain resource.
type IngressRule struct {
	// Ports are all the ports that will be opened.
	Ports []int32
	// IPBlocks are all the IP blocks that can access from outside.
	// Each IPBlock should be expressed using CIDR notation.
	IPBlocks []string
}

// EgressRule groups a set of rules to be applied on a certain resource.
type EgressRule struct {
	// Ports are all the ports that will be opened.
	Ports []int32
	// IPBlocks are all the IP blocks that the resource can communicate to.
	IPBlocks []string
	// AllowOutbound allows outbound traffic to the internet on this pod.
	AllowOutbound bool
}

// CreateNetworkPolicyInput is the input for creating a new network policy.
type CreateNetworkPolicyInput struct {
	// Name is the name of the network policy.
	Name string

	// Namespace is the namespace where this network policy will be created.
	Namespace string

	// Labels is the group of key-value pairs that will identify this policy.
	Labels map[string]string

	// PodSelector are the labels of the pods that this policy applies to.
	PodSelector resource.Selector

	// PeersFrom is the group of pod selectors that are allowed to access the pods covered by this network policy.
	PeersFrom []resource.Selector

	// PeersTo is the group of pod selectors that the pods covered by this network policy are allowed to access.
	PeersTo []resource.Selector

	// Ingresses groups the set of rules to apply to the ingress policy.
	Ingresses IngressRule

	// Egresses groups the set of rules to apply to the egress policy.
	Egresses EgressRule
}

// Policies groups a set of methods to manage network policies.
type Policies interface {
	// Create creates a new network policy.
	Create(input CreateNetworkPolicyInput) (resource.Resource, error)
}
