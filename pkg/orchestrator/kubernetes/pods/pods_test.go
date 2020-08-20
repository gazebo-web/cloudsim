package pods

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/client/conditions"
	"sync"
	"testing"
	"time"
)

func TestNewPods(t *testing.T) {
	client := fake.NewSimpleClientset()
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	m := NewPods(client, f, logger)
	assert.NotNil(t, m)
	assert.IsType(t, &pods{}, m)
	pm := m.(*pods)
	assert.NotNil(t, pm.API)
}

/* TODO: Uncomment this test when addressing the following task:
	https://app.asana.com/0/851925973517080/1188870406911377
func TestPods_Executor(t *testing.T) {
	resource := apiv1.Pod{
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

	client := fake.NewSimpleClientset(&resource)
	fake := spdy.NewSPDYFakeInitializer()
	m := NewPods(client, fake)

	ex := m.Exec(NewPod("test", "default", "test=app"))

	assert.NotNil(t, ex)
	assert.NoError(t, ex.Cmd([]string{"ping", "-c 10", "1.1.1.1"}))

	assert.Equal(t, 1, fake.Calls)
}
*/

func TestPods_WaitForPodsToBeReady(t *testing.T) {
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
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	m := NewPods(client, f, logger)
	r := m.WaitForCondition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

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

func TestPods_WaitForPodsErrWhenPodStateSucceeded(t *testing.T) {
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
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	m := NewPods(client, f, logger)
	r := m.WaitForCondition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

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

func TestPods_WaitForPodsErrWhenPodStateFailed(t *testing.T) {
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
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	m := NewPods(client, f, logger)
	r := m.WaitForCondition(NewPod("test", "default", "test=app"), orchestrator.ReadyCondition)

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
