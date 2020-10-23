package gloo

import (
	"fmt"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// upstreams is a gloo.Upstream implementation of orchestrator.Ingresses.
type upstreams struct {
	API     gloo.GlooV1Interface
	Gateway gateway.GatewayV1Interface
	Logger  ign.Logger
}

// Get returns an upstream from with a certain name in the given namespace.
func (u *upstreams) Get(name string, namespace string) (orchestrator.Resource, error) {
	u.Logger.Debug(fmt.Sprintf("Getting upstream with name [%s] in namespace [%s]", name, namespace))

	out, err := u.API.Upstreams(name).Get(namespace, metav1.GetOptions{})
	if err != nil {
		u.Logger.Debug(fmt.Sprintf("Getting upstream with name [%s] in namespace [%s] failed.", name, namespace))
		return nil, err
	}

	u.Logger.Debug(fmt.Sprintf("Getting upstream with name [%s] in namespace [%s] succeeded.", name, namespace))

	selector := orchestrator.NewSelector(out.Labels)
	return orchestrator.NewResource(name, namespace, selector), nil
}

// NewIngresses initializes a new orchestrator.Ingresses implementation using Gloo.
func NewIngresses(api gloo.GlooV1Interface, gw gateway.GatewayV1Interface, logger ign.Logger) orchestrator.Ingresses {
	return &upstreams{
		API:     api,
		Gateway: gw,
		Logger:  logger,
	}
}
