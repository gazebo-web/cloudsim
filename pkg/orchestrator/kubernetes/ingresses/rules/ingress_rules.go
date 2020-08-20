package rules

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ingressRules is an orchestrator.IngressRules implementation.
type ingressRules struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Get returns the rule definition of the given host from the given resource.
func (m *ingressRules) Get(resource orchestrator.Resource, host string) (orchestrator.Rule, error) {
	m.Logger.Debug(
		fmt.Sprintf(
			"Getting ingress rule with name [%s] in namespace [%s] and with the following selectors: [%s] ",
			resource.Name(), resource.Namespace(), resource.Selector().String(),
		),
	)
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(resource.Namespace()).Get(resource.Name(), metav1.GetOptions{})
	if err != nil {
		m.Logger.Debug(
			fmt.Sprintf(
				"Getting ingress rule with name [%s] failed. Error: [%s]",
				resource.Name(), err.Error(),
			),
		)
		return nil, err
	}
	var rule v1beta1.HTTPIngressRuleValue
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.Host == host {
			rule = *ingressRule.IngressRuleValue.HTTP
		}
	}
	paths := NewPaths(rule.Paths)
	out := NewRule(resource, host, paths)
	m.Logger.Debug(
		fmt.Sprintf(
			"Getting ingress rule with name [%s] succeded. Host: [%s]. Paths: [%+v]",
			resource.Name(), out.Host(), out.Paths(),
		),
	)
	return out, nil
}

// Upsert adds a set of paths to the given host's rule.
// If the paths already exist, it updates them.
func (m *ingressRules) Upsert(rule orchestrator.Rule, paths ...orchestrator.Path) error {
	m.Logger.Debug(fmt.Sprintf("Upserting rule from host [%s] ", rule.Host()))
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
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Error while updating rules from host [%s] ", rule.Host()))
		return err
	}
	m.Logger.Debug(fmt.Sprintf("Rule [%s] has been updated. Paths: [%+v]", rule.Host(), rule.Paths()))
	return nil
}

// Remove removes a set of paths from the given host's rule.
func (m *ingressRules) Remove(rule orchestrator.Rule, paths ...orchestrator.Path) error {
	m.Logger.Debug(fmt.Sprintf("Removing rule paths from host [%s] ", rule.Host()))
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
		m.Logger.Debug(fmt.Sprintf(
			"Error while removing rule paths from host [%s]. Error: %s",
			rule.Host(), orchestrator.ErrRuleNotFound),
		)
		return orchestrator.ErrRuleNotFound
	}

	rule.RemovePaths(paths)

	outputRules := rule.ToOutput()
	ingressRulesInput := outputRules.(v1beta1.IngressRule)
	ingress.Spec.Rules[position] = ingressRulesInput

	_, err = m.API.ExtensionsV1beta1().Ingresses(rule.Resource().Namespace()).Update(ingress)
	if err != nil {
		m.Logger.Debug(fmt.Sprintf(
			"Error while removing rule paths from host [%s]. Error: %s",
			rule.Host(), orchestrator.ErrRuleNotFound),
		)
		return err
	}
	m.Logger.Debug(fmt.Sprintf("Paths from rule host [%s] have been removed. Current paths: [%+v]", rule.Host(), rule.Paths()))
	return nil
}

// NewIngressRules initializes a new orchestrator.IngressRules implementation using Kubernetes.
func NewIngressRules(api kubernetes.Interface, logger ign.Logger) orchestrator.IngressRules {
	return &ingressRules{
		API:    api,
		Logger: logger,
	}
}
