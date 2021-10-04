package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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
	assert.IsType(t, &kubernetesPods{}, m)
	pm := m.(*kubernetesPods)
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
					Status: apiv1.ConditionTrue,
				},
			},
			Phase: apiv1.PodRunning,
		},
	}

	client := fake.NewSimpleClientset(pod)
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	m := NewPods(client, f, logger)
	selector := resource.NewSelector(map[string]string{"test": "app"})
	res := resource.NewResource("test", "default", selector)
	r := m.WaitForCondition(res, resource.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

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
	selector := resource.NewSelector(map[string]string{"test": "app"})
	res := resource.NewResource("test", "default", selector)
	r := m.WaitForCondition(res, resource.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	wg.Wait()
	assert.Error(t, err)
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

	selector := resource.NewSelector(map[string]string{"test": "app"})
	res := resource.NewResource("test", "default", selector)
	r := m.WaitForCondition(res, resource.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	wg.Wait()
	assert.Error(t, err)
}

func TestPods_CreateSuccess(t *testing.T) {
	client := fake.NewSimpleClientset()
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	p := NewPods(client, f, logger)

	res, err := p.Create(pods.CreatePodInput{
		Name:                          "test",
		Namespace:                     "default",
		Labels:                        map[string]string{"app": "test"},
		RestartPolicy:                 pods.RestartPolicyNever,
		TerminationGracePeriodSeconds: time.Second * 5,
		NodeSelector:                  nil,
		Containers: []pods.Container{
			{
				Name:                     "test",
				Image:                    "ignition/test",
				Args:                     nil,
				Privileged:               nil,
				AllowPrivilegeEscalation: nil,
				Ports:                    nil,
				Volumes:                  nil,
				EnvVars:                  nil,
			},
		},
		Volumes:     nil,
		Nameservers: nil,
	})

	assert.NoError(t, err)

	createdPod, err := client.CoreV1().Pods(res.Namespace()).Get(res.Name(), metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, res.Name(), createdPod.Name)
	assert.Equal(t, res.Namespace(), createdPod.Namespace)
	assert.Equal(t, res.Selector().Map(), createdPod.GetLabels())
	assert.Equal(t, apiv1.RestartPolicyNever, createdPod.Spec.RestartPolicy)
	assert.Len(t, createdPod.Spec.Containers, 1)
	assert.Equal(t, "ignition/test", createdPod.Spec.Containers[0].Image)
}

func TestPods_CreateFailsWhenPodAlreadyExists(t *testing.T) {
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
	p := NewPods(client, f, logger)

	_, err := p.Create(pods.CreatePodInput{
		Name:                          "test",
		Namespace:                     "default",
		Labels:                        map[string]string{"app": "test"},
		RestartPolicy:                 pods.RestartPolicyNever,
		TerminationGracePeriodSeconds: time.Second * 5,
		NodeSelector:                  nil,
		Containers: []pods.Container{
			{
				Name:                     "test",
				Image:                    "ignition/test",
				Args:                     nil,
				Privileged:               nil,
				AllowPrivilegeEscalation: nil,
				Ports:                    nil,
				Volumes:                  nil,
				EnvVars:                  nil,
			},
		},
		Volumes:     nil,
		Nameservers: nil,
	})

	assert.Error(t, err)
}

func TestPods_DeleteFailsWhenPodDoesNotExist(t *testing.T) {
	client := fake.NewSimpleClientset()
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	p := NewPods(client, f, logger)

	_, err := p.Delete(
		resource.NewResource("test", "default", resource.NewSelector(map[string]string{})),
	)

	assert.Error(t, err)
}

func TestPods_DeleteSuccess(t *testing.T) {
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
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
	p := NewPods(client, f, logger)

	_, err := p.Delete(
		resource.NewResource("test", "default", resource.NewSelector(map[string]string{
			"app": "test",
		})),
	)

	assert.NoError(t, err)

	_, err = client.CoreV1().Pods("default").Get("test", metav1.GetOptions{})
	assert.Error(t, err)
}

func TestPods_GetFails(t *testing.T) {
	client := fake.NewSimpleClientset()
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	p := NewPods(client, f, logger)

	_, err := p.Get("test", "default")

	assert.Error(t, err)
}

func TestPods_GetSuccess(t *testing.T) {
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
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
	p := NewPods(client, f, logger)

	_, err := p.Get("test", "default")

	assert.NoError(t, err)
}

func TestPods_List(t *testing.T) {
	client := fake.NewSimpleClientset(
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-1",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-2",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-3",
				Namespace: "cloudsim",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
	)
	f := spdy.NewSPDYFakeInitializer()
	logger := ign.NewLoggerNoRollbar("TestPods", ign.VerbosityDebug)
	p := NewPods(client, f, logger)

	// Getting pods in a certain namespace
	list, err := p.List("default", resource.NewSelector(map[string]string{
		"app": "test",
	}))
	require.NoError(t, err)
	assert.Len(t, list, 2)

	// Getting elements from another namespace should only return the elements from that namespace.
	list, err = p.List("cloudsim", resource.NewSelector(map[string]string{
		"app": "test",
	}))
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// A wrong defined selector should return an empty response.
	list, err = p.List("default", resource.NewSelector(map[string]string{
		"app": "undefined",
	}))
	require.NoError(t, err)
	assert.Len(t, list, 0)

	// An empty selector should return all pods in the given namespace.
	list, err = p.List("default", resource.NewSelector(nil))
	require.NoError(t, err)
	assert.Len(t, list, 2)

	// A nil selector should return all pods in the given namespace.
	list, err = p.List("default", nil)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}
