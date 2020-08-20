package rules

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRule_GetRuleReturnsIngressRule(t *testing.T) {
	ing := newTestIngress()

	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	ir := NewIngressRules(client, logger)

	rule, err := ir.Get(ingresses.NewIngress("test", "default"), "test.com")
	assert.NoError(t, err)
	assert.Equal(t, "test.com", rule.Host())
	assert.Len(t, rule.Paths(), 1)
	assert.Equal(t, "test", rule.Paths()[0].Address)
	assert.Equal(t, "test-service", rule.Paths()[0].Endpoint.Name)
	assert.Equal(t, int32(3333), rule.Paths()[0].Endpoint.Port)
}

func TestRule_GetRuleReturnsErrorWhenIngressDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngressRules(client, nil)

	_, err := m.Get(ingresses.NewIngress("test", "default"), "test.com")
	assert.Error(t, err)
}

func TestRule_UpsertRulesReturnsErrorIfIngressDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)
	path := orchestrator.Path{
		Address: "some-regex",
		Endpoint: orchestrator.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	resource := ingresses.NewIngress("test", "default")
	rule := NewRule(resource, "test.org", []orchestrator.Path{})
	err := m.Upsert(rule, path)
	assert.Error(t, err)
}

func TestRule_UpsertRulesReturnsErrorIfRuleDoesNotExist(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)
	path := orchestrator.Path{
		Address: "some-regex",
		Endpoint: orchestrator.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	resource := ingresses.NewIngress("test", "default")
	rule := NewRule(resource, "test.org", []orchestrator.Path{})
	err := m.Upsert(rule, path)
	assert.Error(t, err)
	assert.Equal(t, orchestrator.ErrRuleNotFound, err)
}

func TestRule_UpsertRulesReturnsNoError(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)

	resource := ingresses.NewIngress("test", "default")

	r, err := m.Get(resource, "test.com")
	assert.NoError(t, err)

	path := orchestrator.Path{
		Address: "some-regex",
		Endpoint: orchestrator.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	err = m.Upsert(r, path)
	assert.NoError(t, err)

	result, err := m.Get(resource, "test.com")
	assert.NoError(t, err)
	assert.Len(t, result.Paths(), 2)
}

func TestRule_RemovePathsReturnsNoError(t *testing.T) {
	ing := newTestIngress()

	k8sIngressPathToRemove := v1beta1.HTTPIngressPath{
		Path:    "delete-me",
		Backend: v1beta1.IngressBackend{},
	}

	ing.Spec.Rules[0].HTTP.Paths = append(ing.Spec.Rules[0].HTTP.Paths, k8sIngressPathToRemove)

	client := fake.NewSimpleClientset(&ing)
	logger := ign.NewLoggerNoRollbar("TestRules", ign.VerbosityDebug)
	m := NewIngressRules(client, logger)

	resource := ingresses.NewIngress("test", "default")
	r, err := m.Get(resource, "test.com")
	assert.NoError(t, err)
	assert.Len(t, r.Paths(), 2)

	pathsToRemove := NewPaths([]v1beta1.HTTPIngressPath{k8sIngressPathToRemove})

	err = m.Remove(r, pathsToRemove...)
	assert.NoError(t, err)

	r, err = m.Get(resource, "test.com")
	assert.NoError(t, err)
	assert.Len(t, r.Paths(), 1)
}

func newTestIngress() v1beta1.Ingress {
	return v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: "test.com",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "test",
									Backend: v1beta1.IngressBackend{
										ServiceName: "test-service",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3333,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Status: v1beta1.IngressStatus{},
	}
}
