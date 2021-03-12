package nodes

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// nodes is a orchestrator.Nodes implementation.
type nodes struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// WaitForCondition creates a new wait request that will be used to wait for a resource to match a certain condition.
// The wait request won't be triggered until the method Wait has been called.
func (m *nodes) WaitForCondition(resource orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	m.Logger.Debug(fmt.Sprintf("Creating wait for condition [%+v] request on nodes matching the following selector: [%s]",
		condition, resource.Selector(),
	))

	// Prepare options
	opts := metav1.ListOptions{
		LabelSelector: resource.Selector().String(),
	}

	// Create job
	job := func() (bool, error) {
		var nodesNotReady []*apiv1.Node
		nodes, err := m.API.CoreV1().Nodes().List(opts)
		if err != nil {
			return false, err
		}
		if len(nodes.Items) == 0 {
			return false, nil
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

	m.Logger.Debug(fmt.Sprintf(
		"Wait for condition [%+v] request on nodes matching the following selector: [%s] was created.",
		condition, resource.Selector(),
	))

	// Return new wait request with the created job
	return waiter.NewWaitRequest(job)
}

// isConditionSetAsExpected checks if the given Kubernetes Resource matches the expected condition.
func (m *nodes) isConditionSetAsExpected(node apiv1.Node, expected orchestrator.Condition) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == apiv1.NodeConditionType(expected.Type) &&
			cond.Status == apiv1.ConditionStatus(expected.Status) {
			return true
		}
	}
	return false
}

// NewNodes returns a orchestrator.Nodes implementation with the given kubernetes.Interface API.
func NewNodes(api kubernetes.Interface, logger ign.Logger) orchestrator.Nodes {
	return &nodes{
		API:    api,
		Logger: logger,
	}
}
