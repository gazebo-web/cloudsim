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
	m.Logger.Debug(fmt.Sprintf(
		"Getting ingress rule with name [%s] in namespace [%s] and with the following selectors: [%s] ",
		resource.Name(), resource.Namespace(), resource.Selector().String(),
	))

	// Get ingress from cluster
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(resource.Namespace()).Get(resource.Name(), metav1.GetOptions{})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf(
			"Getting ingress rule with name [%s] failed. Error: [%s]",
			resource.Name(), err.Error(),
		))
		return nil, err
	}

	// Get rule that matches the given host
	var rule v1beta1.HTTPIngressRuleValue
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.Host == host {
			rule = *ingressRule.IngressRuleValue.HTTP
		}
	}

	// Prepare paths and create output
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

	// Get ingress from cluster
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(rule.Resource().Namespace()).Get(rule.Resource().Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Get rules from ingress
	updateRules := ingress.Spec.Rules

	// Find host
	position := -1
	for i, ingressRule := range updateRules {
		if ingressRule.Host == rule.Host() {
			position = i
		}
	}

	// Return error if host wasn't found
	if position == -1 {
		m.Logger.Debug(fmt.Sprintf(
			"Error while updating rule paths from host [%s]. Error: %s",
			rule.Host(), orchestrator.ErrRuleNotFound),
		)
		return orchestrator.ErrRuleNotFound
	}

	// Upsert paths into rule
	rule.UpsertPaths(paths)

	// Update ingress paths
	ingress.Spec.Rules[position] = rule.ToOutput().(v1beta1.IngressRule)

	// Update ingress in cluster
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

	// Get ingress from cluster
	ingress, err := m.API.ExtensionsV1beta1().Ingresses(rule.Resource().Namespace()).Get(rule.Resource().Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Get rules from ingress
	updateRules := ingress.Spec.Rules

	// Find host
	position := -1
	for i, ingressRule := range updateRules {
		if ingressRule.Host == rule.Host() {
			position = i
		}
	}

	// Return an error if the host wasn't found.
	if position == -1 {
		m.Logger.Debug(fmt.Sprintf(
			"Error while removing rule paths from host [%s]. Error: %s",
			rule.Host(), orchestrator.ErrRuleNotFound),
		)
		return orchestrator.ErrRuleNotFound
	}

	// Remove paths from rule
	rule.RemovePaths(paths)

	// Assign new rules to the ingress
	ingress.Spec.Rules[position] = rule.ToOutput().(v1beta1.IngressRule)

	// Update ingress
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
