package pods

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/client/conditions"
	"sync"
	"testing"
	"time"
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

/* TODO: Uncomment this test when addressing the following task:
	https://app.asana.com/0/851925973517080/1188870406911377
func TestManager_Executor(t *testing.T) {
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

	ex := m.Exec(NewPod("test", "default", "test=app"))

	assert.NotNil(t, ex)
	assert.NoError(t, ex.Cmd([]string{"ping", "-c 10", "1.1.1.1"}))

	assert.Equal(t, 1, fake.Calls)
}
*/

func TestManager_WaitForPodsToBeReady(t *testing.T) {
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Spec: apiv1.PodSpec{},
		Status: apiv1.PodStatus{
			Conditions: []apiv1.PodCondition{
				{
					Type:   apiv1.PodReady,
					Status: apiv1.ConditionUnknown,
				},
			},
			Phase: apiv1.PodPending,
		},
	}

	client := fake.NewSimpleClientset(pod)
	fake := spdy.NewSPDYFakeInitializer()
	m := NewManager(client, fake)
	r := m.Condition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	pod.Status.Phase = apiv1.PodRunning
	_, err = client.CoreV1().Pods("default").Update(pod)
	assert.NoError(t, err)

	pod.Status.Conditions = []apiv1.PodCondition{
		{Type: apiv1.PodReady, Status: apiv1.ConditionTrue},
	}
	_, err = client.CoreV1().Pods("default").UpdateStatus(pod)
	assert.NoError(t, err)

	wg.Wait()
	assert.NoError(t, err)
}

func TestManager_WaitForPodsErrWhenPodStateSucceeded(t *testing.T) {
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Spec: apiv1.PodSpec{},
		Status: apiv1.PodStatus{
			Conditions: []apiv1.PodCondition{},
			Phase:      apiv1.PodSucceeded,
		},
	}

	client := fake.NewSimpleClientset(pod)
	fake := spdy.NewSPDYFakeInitializer()
	m := NewManager(client, fake)
	r := m.Condition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, conditions.ErrPodCompleted, err)
}

func TestManager_WaitForPodsErrWhenPodStateFailed(t *testing.T) {
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Spec: apiv1.PodSpec{},
		Status: apiv1.PodStatus{
			Conditions: []apiv1.PodCondition{},
			Phase:      apiv1.PodFailed,
		},
	}

	client := fake.NewSimpleClientset(pod)
	fake := spdy.NewSPDYFakeInitializer()
	m := NewManager(client, fake)
	r := m.Condition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, conditions.ErrPodCompleted, err)
}
