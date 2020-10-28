package gloo

import (
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// rule is an orchestrator.Rule implementation using Gloo Virtual Host.
type rule struct {
	resource orchestrator.Resource
	host     string
	paths    []orchestrator.Path
	domains  []string
}

// Resource returns the resource associated with the current rule.
// It returns the reference to a virtual service.
func (r *rule) Resource() orchestrator.Resource {
	return r.resource
}

// Host returns the Virtual Host domain used to identify this rule.
func (r *rule) Host() string {
	return r.host
}

// Paths returns a list of paths. This list abstracts a set of Gloo routes for a certain virtual host.
func (r *rule) Paths() []orchestrator.Path {
	return r.paths
}

// UpsertPaths inserts and update the given routes into the current virtual host.
func (r *rule) UpsertPaths(paths []orchestrator.Path) {
	r.paths = orchestrator.UpsertPaths(r.paths, paths)
}

// RemovePaths removes the given routes from the current virtual host.
func (r *rule) RemovePaths(paths []orchestrator.Path) {
	r.paths = orchestrator.RemovePaths(r.paths, paths)
}

// ToOutput generates a Gloo representation of a Virtual Host from the current rule.
func (r *rule) ToOutput() interface{} {
	return &gatewayapiv1.VirtualHost{
		Domains: r.domains,
		Routes:  generateRoutes(r.resource.Namespace(), r.paths),
	}
}

// generateRoutes generates a set of routes from the given namespace a list of paths.
func generateRoutes(namespace string, paths []orchestrator.Path) []*gatewayapiv1.Route {
	routes := make([]*gatewayapiv1.Route, 0, len(paths))
	for _, p := range paths {
		routes = append(routes, generateRoute(namespace, p))
	}
	return routes
}

// generateRoute generates a gloo route from the given namespace and path.
func generateRoute(namespace string, path orchestrator.Path) *gatewayapiv1.Route {
	return &gatewayapiv1.Route{
		Matchers: []*matchers.Matcher{
			generateMatcher(path.Address),
		},
		Action: generateRouteAction(namespace, path.Endpoint.Name),
	}
}

// generateMatcher generates a Regex matcher for the given value.
func generateMatcher(value string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: value,
		},
	}
}

// generateRouteAction generates a RouteAction for the given pointing to the given upstream.
func generateRouteAction(namespace string, upstream string) *gatewayapiv1.Route_RouteAction {
	return &gatewayapiv1.Route_RouteAction{
		RouteAction: &glooapiv1.RouteAction{
			Destination: &glooapiv1.RouteAction_Single{
				Single: &glooapiv1.Destination{
					DestinationType: &glooapiv1.Destination_Upstream{
						Upstream: &core.ResourceRef{
							Name:      upstream,
							Namespace: namespace,
						},
					},
				},
			},
		},
	}
}

// NewRule initializes a new orchestrator.Rule implementation using a Gloo.
func NewRule(resource orchestrator.Resource, host string, domains []string, paths ...orchestrator.Path) orchestrator.Rule {
	return &rule{
		resource: resource,
		host:     host,
		domains:  domains,
		paths:    paths,
	}
}
