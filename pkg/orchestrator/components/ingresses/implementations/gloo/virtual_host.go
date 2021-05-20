package gloo

import (
	"fmt"
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	glooapiv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
)

// virtualHosts is a ingresses.IngressRules implementation to manage Gloo Virtual Hosts.
type virtualHosts struct {
	Gateway gateway.GatewayV1Interface
	Logger  ign.Logger
	lock    sync.RWMutex
}

// Get returns a representation of a virtual host using an ingresses.Rule.
func (v *virtualHosts) Get(resource resource.Resource, host string) (ingresses.Rule, error) {
	// Get the virtual service
	vs, err := v.Gateway.VirtualServices(resource.Namespace()).Get(resource.Name(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Get the domains from the underlying virtual host.
	domains := vs.Spec.GetVirtualHost().GetDomains()

	// Find if the given host is part of the domain list
	var domain string
	for _, d := range domains {
		if d == host {
			domain = d
		}
	}

	// Fail if there's no route that matches the host.
	if len(domain) == 0 {
		return nil, ingresses.ErrRuleNotFound
	}

	// Create the path that will represent the routes from a certain virtual host
	paths := NewPaths(vs.Spec.VirtualHost)

	// Create the rule that will represent the virtual host.
	out := NewRule(resource, host, vs.Spec.VirtualHost.Domains, paths...)

	return out, nil
}

// Upsert inserts routes on `paths` in the virtual host represented by `rule`.
func (v *virtualHosts) Upsert(rule ingresses.Rule, paths ...ingresses.Path) error {
	v.Logger.Debug(fmt.Sprintf("Upserting rule for virtual host [%s] ", rule.Host()))

	vs, err := v.Gateway.VirtualServices(rule.Resource().Namespace()).Get(rule.Resource().Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	rule.UpsertPaths(paths)

	vs.Spec.VirtualHost = rule.ToOutput().(*gatewayapiv1.VirtualHost)

	v.lock.Lock()
	defer v.lock.Unlock()

	_, err = v.Gateway.VirtualServices(rule.Resource().Namespace()).Update(vs)
	if err != nil {
		v.Logger.Debug(fmt.Sprintf("Error while updating routes from virtual host [%s]. Error: %s.",
			rule.Host(), err))
		return err
	}

	v.Logger.Debug(fmt.Sprintf("Virtual host [%s] has been updated. Routes: [%+v]", rule.Host(), rule.Paths()))

	return nil
}

// Remove removes the routes listed on `paths` from the virtual host represented by `rule`.
func (v *virtualHosts) Remove(rule ingresses.Rule, paths ...ingresses.Path) error {
	v.Logger.Debug(fmt.Sprintf("Removing routes from virtual host [%s] ", rule.Host()))

	// Get ingress from cluster
	vs, err := v.Gateway.VirtualServices(rule.Resource().Namespace()).Get(rule.Resource().Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(rule.Paths()) == 0 {
		v.Logger.Debug(fmt.Sprintf("Error while removing routes from virtual host [%s]. Error: %s",
			rule.Host(), ingresses.ErrRuleEmpty))
		return ingresses.ErrRuleEmpty
	}

	v.lock.Lock()
	defer v.lock.Unlock()

	// Remove paths from rule
	rule.RemovePaths(paths)

	// Assign new rules to the ingress
	vs.Spec.VirtualHost = rule.ToOutput().(*gatewayapiv1.VirtualHost)

	// Update ingress
	_, err = v.Gateway.VirtualServices(rule.Resource().Namespace()).Update(vs)
	if err != nil {
		v.Logger.Debug(fmt.Sprintf("Error while removing routes from virtual host [%s]. Error: %s",
			rule.Host(), ingresses.ErrRuleNotFound))
		return err
	}

	v.Logger.Debug(fmt.Sprintf("Rotues from virtual host [%s] have been removed. Current routes: [%+v]",
		rule.Host(), rule.Paths()))
	return nil
}

// NewPaths extracts the list of routes from the given virtual host.
// The list of routes will be represented by a slice of ingresses.Path.
func NewPaths(vh *gatewayapiv1.VirtualHost) []ingresses.Path {
	routes := make([]ingresses.Path, len(vh.Routes))
	for i, r := range vh.Routes {
		m := r.Matchers[0]
		ra := r.Action.(*gatewayapiv1.Route_RouteAction)
		routes[i] = NewPath(r.Name, m, ra)
	}
	return routes
}

// NewPath converts a certain matcher and route action into a generic ingresses.Path.
func NewPath(routeName string, matcher *matchers.Matcher, action *gatewayapiv1.Route_RouteAction) ingresses.Path {
	dest := action.RouteAction.Destination.(*glooapiv1.RouteAction_Single)
	up := dest.Single.DestinationType.(*glooapiv1.Destination_Upstream)
	return ingresses.Path{
		UID:     routeName,
		Address: matcher.GetRegex(),
		Endpoint: ingresses.Endpoint{
			Name: up.Upstream.Name,
		},
	}
}

// NewVirtualHosts initializes a new ingresses.IngressRules implementation using Gloo.
func NewVirtualHosts(gw gateway.GatewayV1Interface, logger ign.Logger) ingresses.IngressRules {
	return &virtualHosts{
		Gateway: gw,
		Logger:  logger,
	}
}
