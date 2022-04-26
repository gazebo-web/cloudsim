package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"sync"
	"testing"
	"time"
)

func TestNewNodeNodes(t *testing.T) {
	client := fake.NewSimpleClientset()
	nm := NewNodes(client, ign.NewLoggerNoRollbar("TestNodes", ign.VerbosityDebug))
	assert.NotNil(t, nm)
}

func TestConditionSetAsExpected(t *testing.T) {
	m := &kubernetesNodes{}
	assert.True(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeReady,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeReady,
					Status: apiv1.ConditionFalse,
				},
			},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeReady,
					Status: apiv1.ConditionUnknown,
				},
			},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeDiskPressure,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeMemoryPressure,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}, resource.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeNetworkUnavailable,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}, resource.ReadyCondition))
}

func TestWait_WaitForNodesToBeReady(t *testing.T) {
	node := apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeReady,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodes(cli, ign.NewLoggerNoRollbar("TestNodes", ign.VerbosityDebug))
	selector := resource.NewSelector(map[string]string{"test": "app"})
	res := resource.NewResource("test", "default", selector)
	r := nm.WaitForCondition(res, resource.ReadyCondition)

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

func TestWait_ErrWhenNodesArentReady(t *testing.T) {
	node := apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"test": "app",
			},
		},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeNetworkUnavailable,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodes(cli, ign.NewLoggerNoRollbar("TestNodes", ign.VerbosityDebug))

	selector := resource.NewSelector(map[string]string{"test": "app"})
	res := resource.NewResource("test", "default", selector)
	r := nm.WaitForCondition(res, resource.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()
	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, waiter.ErrRequestTimeout, err)
}
