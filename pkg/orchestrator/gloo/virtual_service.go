package gloo

import (
	"fmt"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// virtualServices is an orchestrator.Ingresses implementation using Gloo.
// It's in charge of managing Gloo Virtual Services.
type virtualServices struct {
	Gateway gateway.GatewayV1Interface
	Logger  ign.Logger
}

// Get returns an orchestrator.Resource of type VirtualService.
func (v *virtualServices) Get(name string, namespace string) (orchestrator.Resource, error) {
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

// NewVirtualServices initializes a new orchestrator.Ingresses implementation using Gloo Virtual Services.
func NewVirtualServices(gw gateway.GatewayV1Interface, logger ign.Logger) orchestrator.Ingresses {
	return &virtualServices{
		Gateway: gw,
		Logger:  logger,
	}
}
