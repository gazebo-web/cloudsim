package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
)

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
func (r ruler) Upsert(rule orchestrator.Rule, paths ...orchestrator.Path) error {
	ing, err := r.manager.get(r.resource)
	if err != nil {
		return err
	}
	updateRules := ing.Spec.Rules
	position := -1
	for i, ingressRule := range updateRules {
		if ingressRule.Host == rule.Host() {
			position = i
		}
	}
	if position == -1 {
		return orchestrator.ErrRuleNotFound
	}

	rule.UpsertPaths(paths)

	outputRules := rule.ToOutput()
	ingressRules := outputRules.(v1beta1.IngressRule)
	ing.Spec.Rules[position] = ingressRules

	_, err = r.manager.update(r.resource, ing)
	return err
}

// Remove removes a set of paths from the given host's rule.
func (r ruler) Remove(host string, paths ...orchestrator.Path) error {
	panic("implement me")
}
