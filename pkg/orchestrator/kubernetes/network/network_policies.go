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
func (n networkPolicies) Create(input orchestrator.CreateNetworkPolicyInput) (orchestrator.Resource, error) {
	specIngress := createIngressSpec(input)
	specEgress := createEgressSpec(input)

	np := &networkingv1.NetworkPolicy{
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

	_, err := n.API.NetworkingV1().NetworkPolicies(input.Namespace).Create(np)
	if err != nil {
		return nil, err
	}
	return orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels)), nil
}

// createEgressSpec creates all the egress rules needed by the networkPolicies.Create method.
func createEgressSpec(input orchestrator.CreateNetworkPolicyInput) []networkingv1.NetworkPolicyEgressRule {
	size := len(input.Egresses.Ports) + len(input.Egresses.IPBlocks) + len(input.PeersTo)
	specEgress := make([]networkingv1.NetworkPolicyEgressRule, size)
	for _, port := range input.Egresses.Ports {
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
	for _, cidr := range input.Egresses.IPBlocks {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

	for _, to := range input.PeersTo {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: to.Map(),
					},
				},
			},
		})
	}
	return specEgress
}

// createIngressSpec creates all the ingress rules needed by the networkPolicies.Create method.
func createIngressSpec(input orchestrator.CreateNetworkPolicyInput) []networkingv1.NetworkPolicyIngressRule {
	size := len(input.Ingresses.Ports) + len(input.Ingresses.IPBlocks) + len(input.PeersFrom)
	specIngress := make([]networkingv1.NetworkPolicyIngressRule, size)

	for _, port := range input.Ingresses.Ports {
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
	for _, cidr := range input.Ingresses.IPBlocks {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

	for _, from := range input.PeersFrom {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: from.Map(),
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
