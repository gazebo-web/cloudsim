package nodes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go"
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
	m := &nodes{}
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
	}, orchestrator.ReadyCondition))

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
	}, orchestrator.ReadyCondition))

	assert.False(t, m.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{},
		},
	}, orchestrator.ReadyCondition))

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
	}, orchestrator.ReadyCondition))

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
	}, orchestrator.ReadyCondition))

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
	}, orchestrator.ReadyCondition))

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
	}, orchestrator.ReadyCondition))
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
					Type:   apiv1.NodePIDPressure,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodes(cli, ign.NewLoggerNoRollbar("TestNodes", ign.VerbosityDebug))
	selector := orchestrator.NewSelector(map[string]string{"test": "app"})
	res := orchestrator.NewResource("test", "default", selector)
	r := nm.WaitForCondition(res, orchestrator.ReadyCondition)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = r.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()

	node.Status.Conditions = append(node.Status.Conditions, apiv1.NodeCondition{Type: apiv1.NodeReady, Status: apiv1.ConditionTrue})

	_, err = cli.CoreV1().Nodes().UpdateStatus(&node)
	assert.NoError(t, err)

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
					Type:   apiv1.NodePIDPressure,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodes(cli, ign.NewLoggerNoRollbar("TestNodes", ign.VerbosityDebug))

	selector := orchestrator.NewSelector(map[string]string{"test": "app"})
	res := orchestrator.NewResource("test", "default", selector)
	r := nm.WaitForCondition(res, orchestrator.ReadyCondition)

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
