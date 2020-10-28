package gloo

import (
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestNewVirtualHosts(t *testing.T) {
	var vhs orchestrator.IngressRules
	var gw gateway.GatewayV1Interface
	var logger ign.Logger
	vhs = NewVirtualHosts(gw, logger)
	assert.IsType(t, &virtualHosts{}, vhs)
}

func TestVirtualHosts_Get(t *testing.T) {

}

func TestVirtualHosts_Upsert(t *testing.T) {

}

func TestVirtualHost_Remove(t *testing.T) {

}

func newTestVirtualService(name, namespace, upstream, regex string, domains []string) *gatewayv1.VirtualService {
	return &gatewayv1.VirtualService{Spec: v1.VirtualService{
		VirtualHost: &v1.VirtualHost{
			Domains: domains,
			Routes: []*v1.Route{
				{
					Matchers: []*matchers.Matcher{
						generateMatcher(regex),
					},
					Action: generateRouteAction(namespace, upstream),
					Name:   name,
				},
			},
		},
		DisplayName: name,
	}}
}