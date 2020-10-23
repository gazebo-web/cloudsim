package gloo

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// ingresses is a gloo implementation of orchestrator.Ingresses.
type ingresses struct {
	API     gloo.GlooV1Interface
	Gateway gateway.GatewayV1Interface
}

// Get returns a Gloo Ingress from with a certain name in the given namespace.
func (i *ingresses) Get(name string, namespace string) (orchestrator.Resource, error) {
	panic("implement me")
}

// NewIngresses initializes a new orchestrator.Ingresses implementation using Gloo.
func NewIngresses(api gloo.GlooV1Interface, gw gateway.GatewayV1Interface) orchestrator.Ingresses {
	return &ingresses{
		API:     api,
		Gateway: gw,
	}
}
