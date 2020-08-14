package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
)

// rule is an orchestrator.Rule implementation.
type rule struct {
	host  string
	paths []orchestrator.Path
}

// Host returns the rule's host.
func (r rule) Host() string {
	return r.host
}

// Paths returns an array of paths.
func (r rule) Paths() []orchestrator.Path {
	return r.paths
}

// NewRule initializes a new orchestrator.Rule.
func NewRule(host string, paths []orchestrator.Path) orchestrator.Rule {
	return &rule{
		host:  host,
		paths: paths,
	}
}

// ruler is an orchestrator.Ruler implementation.
type ruler struct {
	resource orchestrator.Resource
	manager  *manager
}

// Get returns the rule definition of the given host.
func (r ruler) Get(host string) (orchestrator.Rule, error) {
	i, err := r.manager.get(r.resource)
	if err != nil {
		return nil, err
	}
	var rule *v1beta1.HTTPIngressRuleValue
	for _, ingressRule := range i.Spec.Rules {
		if ingressRule.Host == host {
			rule = ingressRule.IngressRuleValue.HTTP
		}
	}
	paths := NewPaths(rule.Paths)
	return NewRule(host, paths), nil
}

// Upsert adds a set of paths to the given host's rule.
func (r ruler) Upsert(host string, paths ...orchestrator.Path) error {
	panic("implement me")
}

// Remove removes a set of paths from the given host's rule.
func (r ruler) Remove(host string, paths ...orchestrator.Path) error {
	panic("implement me")
}
