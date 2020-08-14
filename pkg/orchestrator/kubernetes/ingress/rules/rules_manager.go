package rules

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// manager is an orchestrator.IngressRulesManager implementation.
type manager struct {
	API kubernetes.Interface
}

// Get returns the rule definition of the given host from the given resource.
func (m manager) Get(resource orchestrator.Resource, host string) (orchestrator.Rule, error) {
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(resource.Namespace()).Get(resource.Name(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	var rule *v1beta1.HTTPIngressRuleValue
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.Host == host {
			rule = ingressRule.IngressRuleValue.HTTP
		}
	}
	paths := NewPaths(rule.Paths)
	return NewRule(resource, host, paths), nil
}

// Upsert adds a set of paths to the given host's rule.
// If the paths already exist, it updates them.
func (m manager) Upsert(rule orchestrator.Rule, paths ...orchestrator.Path) error {
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(rule.Resource().Namespace()).Get(rule.Resource().Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	updateRules := ingress.Spec.Rules
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
	ingress.Spec.Rules[position] = ingressRules

	_, err = m.API.ExtensionsV1beta1().Ingresses(rule.Resource().Namespace()).Update(ingress)
	return err
}

// Remove removes a set of paths from the given host's rule.
func (m manager) Remove(host string, paths ...orchestrator.Path) error {
	panic("implement me")
}

func NewManager(api kubernetes.Interface) orchestrator.IngressRulesManager {
	return &manager{API: api}
}
