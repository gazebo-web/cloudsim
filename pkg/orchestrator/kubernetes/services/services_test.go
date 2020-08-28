package services

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestCreateService(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))
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

func TestCreateServiceFailsWhenServiceIsAlreadyCreated(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))
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

	err = s.Create(orchestrator.CreateServiceInput{
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
	assert.Error(t, err)
}

func TestGetServiceFailsWhenServiceDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))
	_, err := s.Get("test", "default")
	assert.Error(t, err)
}

func TestGetServiceSuccessWhenServiceExists(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

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

	result, err := s.Get("service-test", "default")
	assert.NoError(t, err)
	assert.Equal(t, "service-test", result.Name())
}

func TestGetAllServicesSuccess(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

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

	err = s.Create(orchestrator.CreateServiceInput{
		Name:      "service-test2",
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

	result, err := s.GetAllBySelector("default", orchestrator.NewSelector(map[string]string{"service": "test"}))
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetAllServicesFailsWhenUsingWrongLabels(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

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

	err = s.Create(orchestrator.CreateServiceInput{
		Name:      "service-test2",
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

	result, err := s.GetAllBySelector("default", orchestrator.NewSelector(map[string]string{"another": "test"}))
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestGetAllServicesFailsWhenNoServicesDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

	result, err := s.GetAllBySelector("default", orchestrator.NewSelector(map[string]string{"some": "test"}))
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestRemoveServiceSuccessWhenServiceExists(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

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

	res, err := s.Get("service-test", "default")
	assert.NoError(t, err)

	err = s.Remove(res)
	assert.NoError(t, err)
}

func TestRemoveServiceFailsWhenServiceDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	s := NewServices(client, ign.NewLoggerNoRollbar("TestService", ign.VerbosityDebug))

	res := orchestrator.NewResource("test", "default", orchestrator.NewSelector(nil))

	err := s.Remove(res)
	assert.Error(t, err)
}
