package network

import (
	"fmt"
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

//
func (n networkPolicies) Create(input orchestrator.CreateNetworkPolicyInput) (orchestrator.Resource, error) {
	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: input.PodSelector.Map(),
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				// Dev note: Important -- the IP addresses listed here should be the IP of the Cloudsim pod.
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							IPBlock: &networkingv1.IPBlock{
								// We always allow traffic coming from the Cloudsim host.
								CIDR: fmt.Sprintf("%s/32", input.CIDR),
							},
						},
					},
				},
				// Allow traffic to websocket server
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: input.WebsocketPort,
							},
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				// Dev note: Important -- the IP addresses listed here should be the IP of the Cloudsim pod.
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							IPBlock: &networkingv1.IPBlock{
								// We always allow traffic targetted to the Cloudsim host
								CIDR: fmt.Sprintf("%s/32", input.CIDR),
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
		},
	}

	for _, to := range input.PeersTo {
		np.Spec.Egress = append(np.Spec.Egress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: to.Map(),
					},
				},
			},
		})
	}
	for _, from := range input.PeersFrom {
		np.Spec.Egress = append(np.Spec.Egress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: from.Map(),
					},
				},
			},
		})
	}

	_, err := n.API.NetworkingV1().NetworkPolicies(input.Namespace).Create(np)
	if err != nil {
		return nil, err
	}
	return orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels)), nil
}

// NewNetworkPolicies initializes a new orchestrator.NetworkPolicies using Kubernetes.
func NewNetworkPolicies(api kubernetes.Interface, logger ign.Logger) orchestrator.NetworkPolicies {
	return &networkPolicies{
		API:    api,
		Logger: logger,
	}
}
