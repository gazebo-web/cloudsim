package gloo

import (
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

type rule struct {
	resource orchestrator.Resource
	host     string
	paths    []orchestrator.Path
	domains  []string
}

func (r *rule) Resource() orchestrator.Resource {
	return r.resource
}

func (r *rule) Host() string {
	return r.host
}

func (r *rule) Paths() []orchestrator.Path {
	return r.paths
}

func (r *rule) UpsertPaths(paths []orchestrator.Path) {
	for _, p := range paths {
		var updated bool
		for i, rulePath := range r.paths {
			if rulePath.Endpoint == p.Endpoint {
				updated = true
				r.paths[i] = p
				break
			}
		}
		if !updated {
			r.paths = append(r.paths, p)
		}
	}
}

func (r *rule) RemovePaths(paths []orchestrator.Path) {
	for _, p := range paths {
		for i, rulePath := range r.paths {
			if rulePath.Endpoint == p.Endpoint {
				pathsLen := len(r.paths)
				if pathsLen > 1 {
					r.paths[i] = r.paths[pathsLen-1]
				}
				r.paths = r.paths[:pathsLen-1]
				break
			}
		}
	}
}

func (r *rule) ToOutput() interface{} {
	return &gatewayapiv1.VirtualHost{
		Domains: r.domains,
		Routes:  generateRoutes(r.resource.Namespace(), r.paths),
	}
}

func generateRoutes(namespace string, paths []orchestrator.Path) []*gatewayapiv1.Route {
	routes := make([]*gatewayapiv1.Route, 0, len(paths))
	for _, p := range paths {
		routes = append(routes, generateRoute(namespace, p))
	}
	return routes
}

func generateRoute(namespace string, path orchestrator.Path) *gatewayapiv1.Route {
	return &gatewayapiv1.Route{
		Matchers: []*matchers.Matcher{
			generateMatcher(path.Address),
		},
		Action: generateRouteAction(namespace, path.Endpoint.Name),
	}
}

func generateMatcher(addr string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: addr,
		},
	}
}

func generateRouteAction(namespace string, endpointName string) *gatewayapiv1.Route_RouteAction {
	return &gatewayapiv1.Route_RouteAction{
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
}

func NewRule(resource orchestrator.Resource, host string, domains []string, paths ...orchestrator.Path) orchestrator.Rule {
	return &rule{
		resource: resource,
		domains:  domains,
		host:     host,
		paths:    paths,
	}
}
