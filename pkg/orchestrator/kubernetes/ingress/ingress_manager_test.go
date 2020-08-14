package ingress

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewManager(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewManager(client)
	assert.IsType(t, &manager{}, m)
}

func TestManager_GetIngress(t *testing.T) {
	ing := newTestIngress()
	client := fake.NewSimpleClientset(&ing)

	m := NewManager(client)

	mgr, ok := m.(*manager)
	assert.True(t, ok)

	result, err := mgr.get(NewIngress("test", "default"))
	assert.NoError(t, err)
	assert.Equal(t, ing, *result)
}

func TestManager_UpdateIngress(t *testing.T) {
	baseIng := newTestIngress()
	client := fake.NewSimpleClientset(&baseIng)

	m := NewManager(client)
	mgr, ok := m.(*manager)
	assert.True(t, ok)

	updateIng := baseIng
	updateIng.Spec.Rules = []v1beta1.IngressRule{
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
									IntVal: 4444,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := mgr.update(NewIngress("test", "default"), &updateIng)
	assert.NoError(t, err)
	assert.Len(t, result.Spec.Rules[0].HTTP.Paths, 1)
	assert.Equal(t, int32(4444), result.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort.IntVal)
}

func TestManager_GetRuleReturnsIngress(t *testing.T) {
	ing := newTestIngress()

	client := fake.NewSimpleClientset(&ing)
	m := NewManager(client)
	r := m.Rules(NewIngress("test", "default"))

	rule, err := r.Get("test.com")
	assert.NoError(t, err)
	assert.Equal(t, "test.com", rule.Host())
	assert.Len(t, rule.Paths(), 1)
	assert.Equal(t, "test", rule.Paths()[0].Regex)
	assert.Equal(t, "test-service", rule.Paths()[0].Endpoint.Name)
	assert.Equal(t, int32(3333), rule.Paths()[0].Endpoint.Port)
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
