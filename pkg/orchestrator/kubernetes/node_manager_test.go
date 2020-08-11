package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"sync"
	"testing"
	"time"
)

func TestNewNodeManager(t *testing.T) {
	client := fake.NewSimpleClientset()
	nm := NewNodeManager(client)
	assert.NotNil(t, nm)
}

func TestCondition_CreatesWaitRequest(t *testing.T) {
	cli := fake.NewSimpleClientset()
	nm := NewNodeManager(cli)

	n := NewNode("test", "default", "test=app")
	w := nm.Condition(n, orchestrator.ReadyCondition)

	assert.IsType(t, &nodeConditionWaitRequest{}, w)

	wr, ok := w.(*nodeConditionWaitRequest)

	assert.True(t, ok)
	assert.NotNil(t, wr.node)
	assert.Equal(t, orchestrator.ReadyCondition, wr.condition)
}

func TestConditionSetAsExpected(t *testing.T) {
	r := &nodeConditionWaitRequest{}
	assert.True(t, r.isConditionSetAsExpected(apiv1.Node{
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

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
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

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{},
		},
	}, orchestrator.ReadyCondition))

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
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

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
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

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NodeSpec{},
		Status: apiv1.NodeStatus{
			Conditions: []apiv1.NodeCondition{
				{
					Type:   apiv1.NodeKubeletConfigOk,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}, orchestrator.ReadyCondition))

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
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

	assert.False(t, r.isConditionSetAsExpected(apiv1.Node{
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
					Type:   apiv1.NodeKubeletConfigOk,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodeManager(cli)

	n := NewNode("test", "default", "test=app")
	r := nm.Condition(n, orchestrator.ReadyCondition)

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
					Type:   apiv1.NodeKubeletConfigOk,
					Status: apiv1.ConditionTrue,
				},
			},
		},
	}
	cli := fake.NewSimpleClientset(&node)
	nm := NewNodeManager(cli)

	n := NewNode("test", "default", "test=app")
	r := nm.Condition(n, orchestrator.ReadyCondition)

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
