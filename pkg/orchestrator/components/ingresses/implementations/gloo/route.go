package gloo

import (
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
)

// rule is an ingresses.Rule implementation using Gloo Virtual Host.
type rule struct {
	resource resource.Resource
	host     string
	paths    []ingresses.Path
	domains  []string
}

// Resource returns the resource associated with the current rule.
// It returns the reference to a virtual service.
func (r *rule) Resource() resource.Resource {
	return r.resource
}

// Host returns the Virtual Host domain used to identify this rule.
func (r *rule) Host() string {
	return r.host
}

// Paths returns a list of paths. This list abstracts a set of Gloo routes for a certain virtual host.
func (r *rule) Paths() []ingresses.Path {
	return r.paths
}

// UpsertPaths inserts and update the given routes into the current virtual host.
func (r *rule) UpsertPaths(paths []ingresses.Path) {
	r.paths = ingresses.UpsertPaths(r.paths, paths)
}

// RemovePaths removes the given routes from the current virtual host.
func (r *rule) RemovePaths(paths []ingresses.Path) {
	r.paths = ingresses.RemovePaths(r.paths, paths)
}

// ToOutput generates a Gloo representation of a Virtual Host from the current rule.
func (r *rule) ToOutput() interface{} {
	return &gatewayapiv1.VirtualHost{
		Domains: r.domains,
		Routes:  generateRoutes(r.resource.Namespace(), r.paths),
	}
}

// generateRoutes generates a set of routes from the given namespace a list of paths.
func generateRoutes(namespace string, paths []ingresses.Path) []*gatewayapiv1.Route {
	routes := make([]*gatewayapiv1.Route, len(paths))
	for i, p := range paths {
		routes[i] = generateRoute(namespace, p)
	}
	return routes
}

// generateRoute generates a gloo route from the given namespace and path.
func generateRoute(namespace string, path ingresses.Path) *gatewayapiv1.Route {
	return &gatewayapiv1.Route{
		Matchers: []*matchers.Matcher{
			GenerateRegexMatcher(path.Address),
		},
		Action: GenerateRouteAction(namespace, path.Endpoint.Name),
	}
}

// GenerateRegexMatcher generates a Regex matcher for the given value.
func GenerateRegexMatcher(value string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: value,
		},
	}
}

// GenerateRouteAction generates a RouteAction for the given pointing to the given upstream.
func GenerateRouteAction(namespace string, upstream string) *gatewayapiv1.Route_RouteAction {
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

// NewRule initializes a new ingresses.Rule implementation using a Gloo.
func NewRule(resource resource.Resource, host string, domains []string, paths ...ingresses.Path) ingresses.Rule {
	return &rule{
		resource: resource,
		host:     host,
		domains:  domains,
		paths:    paths,
	}
}
