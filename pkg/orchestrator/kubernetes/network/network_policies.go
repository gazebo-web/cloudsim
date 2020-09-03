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
	specIngress := np.createIngressSpec(input.Ingresses, input.PeersFrom)

	specEgress := np.createEgressSpec(input.Egresses, input.PeersTo)

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

	_, err := np.API.NetworkingV1().NetworkPolicies(input.Namespace).Create(createNetworkPolicy)
	if err != nil {
		return nil, err
	}

	return orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels)), nil
}

// createEgressSpec creates all the egress rules needed by the networkPolicies.Create method.
func (np *networkPolicies) createEgressSpec(egressRule orchestrator.NetworkEgressRule,
	to []orchestrator.Selector) []networkingv1.NetworkPolicyEgressRule {

	size := len(egressRule.Ports) + len(egressRule.IPBlocks) + len(to)

	specEgress := make([]networkingv1.NetworkPolicyEgressRule, size)

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

	for _, cidr := range egressRule.IPBlocks {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

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

	if egressRule.AllowOutbound {
		specEgress = append(specEgress, networkingv1.NetworkPolicyEgressRule{})
	}

	return specEgress
}

// createIngressSpec creates all the ingress rules needed by the networkPolicies.Create method.
func (np *networkPolicies) createIngressSpec(ingressRule orchestrator.NetworkIngressRule,
	from []orchestrator.Selector) []networkingv1.NetworkPolicyIngressRule {

	size := len(ingressRule.Ports) + len(ingressRule.IPBlocks) + len(from)

	specIngress := make([]networkingv1.NetworkPolicyIngressRule, size)

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

	for _, cidr := range ingressRule.IPBlocks {
		specIngress = append(specIngress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{CIDR: cidr},
				},
			},
		})
	}

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
