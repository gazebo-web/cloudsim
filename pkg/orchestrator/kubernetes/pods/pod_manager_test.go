package pods

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewManager(t *testing.T) {
	client := fake.NewSimpleClientset()
	fake := spdy.NewSPDYFakeInitializer()
	m := NewManager(client, fake)
	assert.NotNil(t, m)
	assert.IsType(t, &manager{}, m)
	pm := m.(*manager)
	assert.NotNil(t, pm.API)
}

func TestNewManager_Executor(t *testing.T) {
	pod := apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Spec:   apiv1.PodSpec{},
		Status: apiv1.PodStatus{},
	}

	client := fake.NewSimpleClientset(&pod)
	fake := spdy.NewSPDYFakeInitializer()
	m := NewManager(client, fake)

	ex := m.Exec(NewPod("test", "default", "app=test"))

	assert.NotNil(t, ex)
	assert.NoError(t, ex.Cmd([]string{"ping", "-c 10", "1.1.1.1"}))

	assert.Equal(t, 1, fake.Calls)
}
