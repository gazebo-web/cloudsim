package gloo

import (
	"errors"
	"fmt"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualServices is an orchestrator.Ingresses implementation using Gloo.
// It's in charge of managing Gloo Virtual Services.
type VirtualServices struct {
	Client  gloo.GlooV1Interface
	Gateway gateway.GatewayV1Interface
	Logger  ign.Logger
}

// Get returns an orchestrator.Resource of type VirtualService.
func (v *VirtualServices) Get(name string, namespace string) (orchestrator.Resource, error) {
	v.Logger.Debug(fmt.Sprintf("Getting virtual service with name [%s] in namespace [%s]", name, namespace))
	vs, err := v.Gateway.VirtualServices(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		v.Logger.Debug(fmt.Sprintf("Getting virtual service with name [%s] in namespace [%s] failed. Error: %s.",
			name, namespace, err))
		return nil, err
	}

	v.Logger.Debug(fmt.Sprintf("Getting virtual service with name [%s] in namespace [%s] succeeded.", name, namespace))

	s := orchestrator.NewSelector(vs.Labels)
	return orchestrator.NewResource(vs.Name, vs.Namespace, s), nil
}

// GetUpstream is used to get the matching upstream for a certain service in the cluster identified by the given selector.
func (v *VirtualServices) GetUpstream(namespace string, selector orchestrator.Selector) (orchestrator.Resource, error) {
	v.Logger.Debug(
		fmt.Sprintf("Getting upstream on namespace [%s] pointing to the given labels [%s]",
			namespace, selector.Map()),
	)

	list, err := v.Client.Upstreams(namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		v.Logger.Debug(
			fmt.Sprintf("Failed to get upstream on namespace [%s] pointing to the given labels [%s]. Error: %s",
				namespace, selector.Map()),
		)
		return nil, err
	}

	if len(list.Items) < 1 {
		return nil, errors.New("did not find a Gloo upstream for target Kubernetes service")
	} else if len(list.Items) > 1 {
		return nil, errors.New("found too many Gloo upstreams for target Kubernetes service")
	}

	s := orchestrator.NewSelector(list.Items[0].Labels)
	res := orchestrator.NewResource(list.Items[0].Name, namespace, s)
	return res, nil
}

// NewVirtualServices initializes a new orchestrator.Ingresses implementation using Gloo Virtual Services.
func NewVirtualServices(gw gateway.GatewayV1Interface, logger ign.Logger, client gloo.GlooV1Interface) orchestrator.Ingresses {
	return &VirtualServices{
		Gateway: gw,
		Logger:  logger,
		Client:  client,
	}
}
