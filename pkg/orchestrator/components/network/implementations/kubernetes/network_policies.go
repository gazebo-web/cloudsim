package kubernetes

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// networkPolicies is a network.Policies implementation.
type networkPolicies struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// RemoveBulk removes a set of network policies specified by the given selector in a certain namespace.
func (np *networkPolicies) RemoveBulk(ctx context.Context, namespace string, selector resource.Selector) error {
	np.Logger.Debug(
		fmt.Sprintf("Removing network policies on namespace [%s] with the given selector: [%s]",
			namespace, selector.String(),
		),
	)

	err := np.API.NetworkingV1().NetworkPolicies(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		np.Logger.Debug(
			fmt.Sprintf("Removing network policies in namespace [%s] with selector: [%s] failed. Error: %s",
				namespace, selector.String(), err,
			),
		)
		return err
	}
	np.Logger.Debug(
		fmt.Sprintf("Removing network policies in namespace [%s] with selector: [%s] succeeded.",
			namespace, selector.String(),
		),
	)
	return nil
}

// Remove removes a network policy with the given name and living in the given namespace.
func (np *networkPolicies) Remove(ctx context.Context, name string, namespace string) error {
	np.Logger.Debug(fmt.Sprintf("Removing network policy with name [%s] in namespace [%s]", name, namespace))

	err := np.API.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		np.Logger.Debug(fmt.Sprintf("Removing network policy with name [%s] in namespace [%s] failed. Error: %s", name, namespace, err))
		return err
	}
	np.Logger.Debug(fmt.Sprintf("Removing network policy with name [%s] in namespace [%s] succeeded.", name, namespace))
	return nil
}

// Create creates a network policy.
func (np *networkPolicies) Create(ctx context.Context, input network.CreateNetworkPolicyInput) (resource.Resource, error) {
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

	np.Logger.Debug(fmt.Sprintf("Creating network policy with name [%s] in namespace [%s]", input.Name, input.Namespace))

	// Create network policy
	_, err := np.API.NetworkingV1().NetworkPolicies(input.Namespace).Create(ctx, createNetworkPolicy, metav1.CreateOptions{})
	if err != nil {
		np.Logger.Debug(fmt.Sprintf("Creating network policy with name [%s] in namespace [%s] failed. Error: %s", input.Name, input.Namespace, err))
		return nil, err
	}

	np.Logger.Debug(fmt.Sprintf("Creating network policy with name [%s] in namespace [%s] succeeded", input.Name, input.Namespace))
	return resource.NewResource(input.Name, input.Namespace, resource.NewSelector(input.Labels)), nil
}

// createEgressSpec creates all the egress rules needed by the networkPolicies.Create method.
func (np *networkPolicies) createEgressSpec(egressRule network.EgressRule,
	to []resource.Selector) []networkingv1.NetworkPolicyEgressRule {

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
func (np *networkPolicies) createIngressSpec(ingressRule network.IngressRule,
	from []resource.Selector) []networkingv1.NetworkPolicyIngressRule {
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

// NewNetworkPolicies initializes a new network.Policies using Kubernetes.
func NewNetworkPolicies(api kubernetes.Interface, logger ign.Logger) network.Policies {
	return &networkPolicies{
		API:    api,
		Logger: logger,
	}
}
