package kubernetes

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"time"
)

type nodeManager struct {
	API kubernetes.Interface
}

type nodeConditionWaitRequest struct {
	nodeManager *nodeManager
	node        orchestrator.Resource
	condition   orchestrator.Condition
}

func (r nodeConditionWaitRequest) Wait(timeout time.Duration, pollFrequency time.Duration) error {
	var nodesNotReady []*apiv1.Node
	opts := metav1.ListOptions{
		LabelSelector: r.node.Selector(),
	}
	err := wait.PollImmediate(pollFrequency, timeout, func() (bool, error) {
		nodesNotReady = nil
		nodes, err := r.nodeManager.API.CoreV1().Nodes().List(opts)
		if err != nil {
			return false, err
		}
		for _, n := range nodes.Items {
			if !r.isConditionSetAsExpected(n, orchestrator.ReadyCondition) {
				var node = new(apiv1.Node)
				*node = n
				nodesNotReady = append(nodesNotReady, node)
			}
		}
		return len(nodesNotReady) == 0, nil
	})
	if err != nil && err != wait.ErrWaitTimeout {
		return err
	}
	if len(nodesNotReady) > 0 {
		return fmt.Errorf("nodes are not ready: %#v", nodesNotReady)
	}
	return nil
}

// isConditionSetAsExpected checks if the given condition is met by a Node.
func (r nodeConditionWaitRequest) isConditionSetAsExpected(node apiv1.Node, expected orchestrator.Condition) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == apiv1.NodeConditionType(expected.Type) &&
			cond.Status == apiv1.ConditionStatus(expected.Status) {
			return true
		}
	}
	return false
}

func (n *nodeManager) Condition(node orchestrator.Resource, condition orchestrator.Condition) orchestrator.Waiter {
	return &nodeConditionWaitRequest{
		nodeManager: n,
		node:        node,
		condition:   condition,
	}
}

func NewNodeManager(api kubernetes.Interface) orchestrator.NodeManager {
	return &nodeManager{
		API: api,
	}
}
