package gloo

import (
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"testing"
)

func TestGenerateRouteAction(t *testing.T) {
	const endpointName string = "test"
	const namespace string = "testing"

	r := &gatewayapiv1.Route_RouteAction{
		RouteAction: &glooapiv1.RouteAction{
			Destination: &glooapiv1.RouteAction_Single{
				Single: &glooapiv1.Destination{
					DestinationType: &glooapiv1.Destination_Upstream{
						Upstream: &core.ResourceRef{
							Name:      endpointName,
							Namespace: namespace,
						},
					},
				},
			},
		},
	}

	out := generateRouteAction(namespace, endpointName)

	assert.Equal(t, r, out)
}

func TestGenerateMatchers(t *testing.T) {
	const addr string = "/test"
	m := &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: addr,
		},
	}

	out := generateMatcher(addr)

	assert.Equal(t, m, out)
}

func TestGenerateRoute(t *testing.T) {
	const addr string = "test"
	const namespace string = "default"
	const endpoint string = "/"

	p := NewPath("test-route", generateMatcher(addr), generateRouteAction(namespace, endpoint))
	out := generateRoute(namespace, p)

	r := &gatewayapiv1.Route{
		Matchers: []*matchers.Matcher{
			generateMatcher(addr),
		},
		Action: generateRouteAction(namespace, endpoint),
	}

	assert.Equal(t, r, out)
}

func TestGenerateRoutes(t *testing.T) {
	const addr string = "test"
	const namespace string = "default"
	const endpoint string = "/"

	p := NewPath("test-route", generateMatcher(addr), generateRouteAction(namespace, endpoint))
	out := generateRoutes(namespace, []orchestrator.Path{p, p, p})

	exp := []*gatewayapiv1.Route{generateRoute(namespace, p), generateRoute(namespace, p), generateRoute(namespace, p)}

	assert.Equal(t, exp, out)
}

func TestRoute_ToOutput(t *testing.T) {
	const ns string = "default"
	const addr string = "test"
	const endpoint string = "/"
	var dom = []string{"openrobotics.org"}

	res := orchestrator.NewResource("test", ns, nil)

	p := NewPath("test-route", generateMatcher(addr), generateRouteAction(ns, endpoint))
	r := NewRule(res, "somehost", []string{"openrobotics.org"}, p)

	vh := &gatewayapiv1.VirtualHost{
		Domains: dom,
		Routes:  generateRoutes(ns, []orchestrator.Path{p}),
	}

	assert.IsType(t, vh, r.ToOutput())
	assert.Equal(t, vh, r.ToOutput())
}
