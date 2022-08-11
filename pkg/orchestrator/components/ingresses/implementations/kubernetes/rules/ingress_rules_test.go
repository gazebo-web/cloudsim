package rules

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRule_GetRuleReturnsIngressRule(t *testing.T) {
	ing := newTestIngress()

	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	ir := NewIngressRules(client, logger)

	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)
	rule, err := ir.Get(context.TODO(), res, "test.com")
	assert.NoError(t, err)
	assert.Equal(t, "test.com", rule.Host())
	assert.Len(t, rule.Paths(), 1)
	assert.Equal(t, "test", rule.Paths()[0].Address)
	assert.Equal(t, "test-service", rule.Paths()[0].Endpoint.Name)
	assert.Equal(t, int32(3333), rule.Paths()[0].Endpoint.Port)
}

func TestRule_GetRuleReturnsErrorWhenIngressDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngressRules(client, ign.NewLoggerNoRollbar("TestRule", ign.VerbosityDebug))
	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)
	_, err := m.Get(context.TODO(), res, "test.com")
	assert.Error(t, err)
}

func TestRule_UpsertRulesReturnsErrorIfIngressDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)
	path := ingresses.Path{
		Address: "some-regex",
		Endpoint: ingresses.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)
	rule := NewRule(res, "test.org", []ingresses.Path{})
	err := m.Upsert(context.TODO(), rule, path)
	assert.Error(t, err)
}

func TestRule_UpsertRulesReturnsErrorIfRuleDoesNotExist(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)
	path := ingresses.Path{
		Address: "some-regex",
		Endpoint: ingresses.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)
	rule := NewRule(res, "test.org", []ingresses.Path{})
	err := m.Upsert(context.TODO(), rule, path)
	assert.Error(t, err)
	assert.Equal(t, ingresses.ErrRuleNotFound, err)
}

func TestRule_UpsertRulesReturnsNoError(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)

	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)

	r, err := m.Get(context.TODO(), res, "test.com")
	assert.NoError(t, err)

	path := ingresses.Path{
		Address: "some-regex",
		Endpoint: ingresses.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	err = m.Upsert(context.TODO(), r, path)
	assert.NoError(t, err)

	result, err := m.Get(context.TODO(), res, "test.com")
	assert.NoError(t, err)
	assert.Len(t, result.Paths(), 2)
}

func TestRule_RemovePathsReturnsNoError(t *testing.T) {
	ing := newTestIngress()

	k8sIngressPathToRemove := networkingv1.HTTPIngressPath{
		Path: "delete-me",
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: "test",
				Port: networkingv1.ServiceBackendPort{
					Number: 1234,
				},
			},
		},
	}

	ing.Spec.Rules[0].HTTP.Paths = append(ing.Spec.Rules[0].HTTP.Paths, k8sIngressPathToRemove)

	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)

	selector := resource.NewSelector(nil)
	res := resource.NewResource("test", "default", selector)
	r, err := m.Get(context.TODO(), res, "test.com")
	assert.NoError(t, err)
	assert.Len(t, r.Paths(), 2)

	pathsToRemove := NewPaths([]networkingv1.HTTPIngressPath{k8sIngressPathToRemove})

	err = m.Remove(context.TODO(), r, pathsToRemove...)
	assert.NoError(t, err)

	r, err = m.Get(context.TODO(), res, "test.com")
	assert.NoError(t, err)
	assert.Len(t, r.Paths(), 1)
}

func newTestIngress() networkingv1.Ingress {
	return networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "test.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "test",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 3333,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Status: networkingv1.IngressStatus{},
	}
}
