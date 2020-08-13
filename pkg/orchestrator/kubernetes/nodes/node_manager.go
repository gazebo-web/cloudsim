package nodes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// manager is a orchestrator.NodeManager implementation.
type manager struct {
	API kubernetes.Interface
}

// Condition returns a waiter.Waiter request to wait until a node reaches the given condition.
func (m *manager) Condition(node orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	var nodesNotReady []*apiv1.Node
	opts := metav1.ListOptions{
		LabelSelector: node.Selector(),
	}
	job := func() (bool, error) {
		nodesNotReady = nil
		nodes, err := m.API.CoreV1().Nodes().List(opts)
		if err != nil {
			return false, err
		}
		for _, n := range nodes.Items {
			if !m.isConditionSetAsExpected(n, condition) {
				var node = new(apiv1.Node)
				*node = n
				nodesNotReady = append(nodesNotReady, node)
			}
		}
		return len(nodesNotReady) == 0, nil
	}
	return waiter.NewWaitRequest(job)
}

// isConditionSetAsExpected checks if the given Kubernetes Node matches the expected condition.
func (m *manager) isConditionSetAsExpected(node apiv1.Node, expected orchestrator.Condition) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == apiv1.NodeConditionType(expected.Type) &&
			cond.Status == apiv1.ConditionStatus(expected.Status) {
			return true
		}
	}
	return false
}

// NewManager returns a orchestrator.NodeManager implementation with the given kubernetes.Interface API.
func NewManager(api kubernetes.Interface) orchestrator.NodeManager {
	return &manager{
		API: api,
	}
}
