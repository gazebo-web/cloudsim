package rules

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingress"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestManager_GetRuleReturnsIngressRule(t *testing.T) {
	ing := newTestIngress()

	client := fake.NewSimpleClientset(&ing)
	rulesManager := NewManager(client)

	rule, err := rulesManager.Get(ingress.NewIngress("test", "default"), "test.com")
	assert.NoError(t, err)
	assert.Equal(t, "test.com", rule.Host())
	assert.Len(t, rule.Paths(), 1)
	assert.Equal(t, "test", rule.Paths()[0].Regex)
	assert.Equal(t, "test-service", rule.Paths()[0].Endpoint.Name)
	assert.Equal(t, int32(3333), rule.Paths()[0].Endpoint.Port)
}

func TestManager_GetRuleReturnsErrorWhenIngressDoesntExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewManager(client)

	_, err := m.Get(ingress.NewIngress("test", "default"), "test.com")
	assert.Error(t, err)
}

func TestManager_UpsertRulesReturnsErrorIfIngressDoesntExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewManager(client)
	path := orchestrator.Path{
		Regex: "some-regex",
		Endpoint: orchestrator.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	err := m.Upsert(NewRule(ingress.NewIngress("test", "default"), "test.org", []orchestrator.Path{}), path)
	assert.Error(t, err)
}

func TestManager_UpsertRulesReturnsErrorIfRuleDoesntExist(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	m := NewManager(client)
	path := orchestrator.Path{
		Regex: "some-regex",
		Endpoint: orchestrator.Endpoint{
			Name: "http",
			Port: 80,
		},
	}
	err := m.Upsert(NewRule(ingress.NewIngress("test", "default"), "test.org", []orchestrator.Path{}), path)
	assert.Error(t, err)
	assert.Equal(t, orchestrator.ErrRuleNotFound, err)
}

func TestManager_UpsertRulesReturnsNoError(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)
	m := NewManager(client)

	resource := ingress.NewIngress("test", "default")

	r, err := m.Get(resource, "test.com")
	assert.NoError(t, err)

	path := orchestrator.Path{
		Regex: "some-regex",
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
