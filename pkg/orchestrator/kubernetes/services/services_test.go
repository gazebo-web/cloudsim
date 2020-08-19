package services

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestCreateNewService(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client)
	err := s.Create(orchestrator.CreateServiceInput{
		Name:      "service-test",
		Type:      "ClusterIP",
		Namespace: "default",
		ServiceLabels: map[string]string{
			"service": "test",
		},
		TargetLabels: map[string]string{
			"app": "test",
		},
		Ports: map[string]int32{
			"http": 80,
		},
	})
	assert.NoError(t, err)
	result, err := client.CoreV1().Services("default").Get("service-test", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, corev1.ServiceTypeClusterIP, result.Spec.Type)
}
