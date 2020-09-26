package gloo

import (
	"context"
	"errors"
	"fmt"
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getVirtualService retrieves a Gloo Virtual Service from the cluster.
// Virtual Services are used to setup routing rules in Gloo.
func getVirtualService(ctx context.Context, gloo Clientset, glooNamespace string,
	virtualServiceName string) (*gatewayv1.VirtualService, error) {
	return gloo.Gateway().VirtualServices(glooNamespace).Get(virtualServiceName, metav1.GetOptions{})
}

// updateVirtualService updates a Gloo Virtual Service resource in the cluster.
func updateVirtualService(ctx context.Context, gloo Clientset, glooNamespace string,
	virtualService *gatewayv1.VirtualService) (*gatewayv1.VirtualService, error) {

	virtualService, err := gloo.Gateway().VirtualServices(glooNamespace).Update(virtualService)
	if err != nil {
		// Get the Virtual Service name
		vsName := "nil"
		if virtualService != nil {
			vsName = virtualService.Name
		}
		errMsg := fmt.Sprintf("failed to update virtual service [%s] in namespace [%s]", vsName, glooNamespace)
		ign.LoggerFromContext(ctx).Error(errMsg, err)
		return nil, errors.New(errMsg)
	}

	return virtualService, nil
}

// UpsertVirtualServiceRoute inserts or updates a set of Gloo Virtual Service routes.
// The `virtualService` parameter is used to select the virtual service to add paths to.
func UpsertVirtualServiceRoute(ctx context.Context, gloo Clientset, glooNamespace string,
	virtualServiceName string, routes ...*gatewayapiv1.Route) (*gatewayv1.VirtualService, error) {

	// Get the virtual service from the cluster
	virtualService, err := getVirtualService(ctx, gloo, glooNamespace, virtualServiceName)
	if err != nil {
		return nil, err
	}

	// The Virtual Host of a Virtual Service contains routing rules
	virtualHost := virtualService.Spec.VirtualHost

	for _, route := range routes {
		// Try to find and update the route
		updated := false
		for i, vhRoute := range virtualHost.Routes {
			if vhRoute.Name == route.Name {
				updated = true
				if route != nil {
					virtualHost.Routes[i] = route
				}
				break
			}
		}
		// No route was updated, create a new one
		if !updated && route != nil {
			virtualHost.Routes = append(virtualHost.Routes, route)
		}
	}

	// Apply updated rule
	return updateVirtualService(ctx, gloo, glooNamespace, virtualService)
}

// RemoveVirtualServiceRoute removes a set of Gloo Virtual Service routes.
// The `virtualService` parameter is used to select the virtual service from which to remove paths.
func RemoveVirtualServiceRoute(ctx context.Context, gloo Clientset, glooNamespace string,
	virtualServiceName string, routes ...*gatewayapiv1.Route) (*gatewayv1.VirtualService, error) {

	// Get the virtual service from the cluster
	virtualService, err := getVirtualService(ctx, gloo, glooNamespace, virtualServiceName)
	if err != nil {
		return nil, err
	}

	// The Virtual Host contains routing rules
	virtualHost := virtualService.Spec.VirtualHost

	// Remove paths
	// To remove routes, they are sent to the back of the slice and the slice is shrunk
	for _, route := range routes {
		routesLen := len(virtualHost.Routes)

		// No more routes to remove
		if routesLen == 0 {
			break
		}

		for i, vhRoute := range virtualHost.Routes {
			if vhRoute.Name == route.Name {
				if routesLen > 1 {
					virtualHost.Routes[i] = virtualHost.Routes[routesLen-1]
				}
				virtualHost.Routes = virtualHost.Routes[:routesLen-1]
				break
			}
		}
	}

	// Apply updated rule
	return updateVirtualService(ctx, gloo, glooNamespace, virtualService)
}

// CreateVirtualHostRouteExactMatcher returns an exact route matcher.
// The route of the request must match the exact string for this matcher to pass.
func CreateVirtualHostRouteExactMatcher(matcher string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Exact{
			Exact: matcher,
		},
	}
}

// CreateVirtualHostRoutePrefixMatcher returns a prefix route matcher.
// The route of the request must be prefixed with the matcher string for this matcher to pass.
func CreateVirtualHostRoutePrefixMatcher(matcher string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Prefix{
			Prefix: matcher,
		},
	}
}

// CreateVirtualHostRouteRegexMatcher returns a regex route matcher.
// The route of the request must match with the passed regex for this matcher to pass.
func CreateVirtualHostRouteRegexMatcher(matcher string) *matchers.Matcher {
	return &matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: matcher,
		},
	}
}

// CreateVirtualHostRouteAction returns a route action that points to a single upstream.
// `namespace` is the upstream namespace. Unless Gloo was installed directly in a namespace, this will typically
// be the Gloo namespace.
func CreateVirtualHostRouteAction(namespace string, upstreamName string) *gatewayapiv1.Route_RouteAction {
	return &gatewayapiv1.Route_RouteAction{
		RouteAction: &glooapiv1.RouteAction{
			Destination: &glooapiv1.RouteAction_Single{
				Single: &glooapiv1.Destination{
					DestinationType: &glooapiv1.Destination_Upstream{
						Upstream: &core.ResourceRef{
							Name:      upstreamName,
							Namespace: namespace,
						},
					},
				},
			},
		},
	}
}

// CreateVirtualHostRoute returns a Virtual Host route that binds a set of matchers with an action.
// TODO: Support other Action types.
func CreateVirtualHostRoute(name string, matchers []*matchers.Matcher,
	action *gatewayapiv1.Route_RouteAction) *gatewayapiv1.Route {

	return &gatewayapiv1.Route{
		Name:     name,
		Matchers: matchers,
		Action:   action,
	}
}
