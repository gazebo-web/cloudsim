package network

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// networkPolicies is an orchestrator.NetworkPolicies implementation.
type networkPolicies struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Create creates a network policy.
func (np *networkPolicies) Create(input orchestrator.CreateNetworkPolicyInput) (orchestrator.Resource, error) {
	// Prepare ingress spec
	specIngress := np.createIngressSpec(input.Ingresses, input.PeersFrom)

	// Prepare egress spec
	specEgress := np.createEgressSpec(input.Egresses, input.PeersTo)

	// Prepare input for Kubernetes
	createNetworkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: input.PodSelector.Map(),
			},
			Ingress:     specIngress,
			Egress:      specEgress,
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
		},
	}

	// Create network policy
	_, err := np.API.NetworkingV1().NetworkPolicies(input.Namespace).Create(createNetworkPolicy)
	if err != nil {
		return nil, err
	}

	return orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels)), nil
}

// createEgressSpec creates all the egress rules needed by the networkPolicies.Create method.
func (np *networkPolicies) createEgressSpec(egressRule orchestrator.NetworkEgressRule,
	to []orchestrator.Selector) []networkingv1.NetworkPolicyEgressRule {

	// Calculate NetworkPolicyEgressRule slice capacity
	capacity := len(egressRule.Ports) + len(egressRule.IPBlocks) + len(to) + 1

	// Define specEgress slice with given capacity
	specEgress := make([]networkingv1.NetworkPolicyEgressRule, 0, capacity)

	// Add ports
	for _, port := range egressRule.Ports {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: port,
					},
				},
			},
		})
	}

	// Add IP Blocks
	for _, cidr := range egressRule.IPBlocks {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

	// Add peers that will point to this pod
	for _, t := range to {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: t.Map(),
					},
				},
			},
		})
	}

	// Allow outbound traffic enabling connection to the internet on this pod.
	if egressRule.AllowOutbound {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{})
	}

	return specEgress
}

// createIngressSpec creates all the ingress rules needed by the networkPolicies.Create method.
func (np *networkPolicies) createIngressSpec(ingressRule orchestrator.NetworkIngressRule,
	from []orchestrator.Selector) []networkingv1.NetworkPolicyIngressRule {
	// Calculate NetworkPolicyIngressRule slice capacity
	capacity := len(ingressRule.Ports) + len(ingressRule.IPBlocks) + len(from)

	// Define slice with the given capacity
	specIngress := make([]networkingv1.NetworkPolicyIngressRule, 0, capacity)

	// Add ports
	for _, port := range ingressRule.Ports {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: port,
					},
				},
			},
		})
	}

	// Add ip blocks
	for _, cidr := range ingressRule.IPBlocks {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

	// Add peers that will point from this pod.
	for _, f := range from {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: f.Map(),
					},
				},
			},
		})
	}

	return specIngress
}

// NewNetworkPolicies initializes a new orchestrator.NetworkPolicies using Kubernetes.
func NewNetworkPolicies(api kubernetes.Interface, logger ign.Logger) orchestrator.NetworkPolicies {
	return &networkPolicies{
		API:    api,
		Logger: logger,
	}
}
