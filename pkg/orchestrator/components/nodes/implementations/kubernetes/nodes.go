package kubernetes

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// kubernetesNodes is a nodes.Nodes implementation.
type kubernetesNodes struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// WaitForCondition creates a new wait request that will be used to wait for a resource to match a certain condition.
// The wait request won't be triggered until the method Wait has been called.
func (m *kubernetesNodes) WaitForCondition(ctx context.Context, node resource.Resource, condition resource.Condition) waiter.Waiter {
	m.Logger.Debug(fmt.Sprintf("Creating wait for condition [%+v] request on nodes matching the following selector: [%s]",
		condition, node.Selector(),
	))

	// Prepare options
	opts := metav1.ListOptions{
		LabelSelector: node.Selector().String(),
	}

	// Create job
	job := func() (bool, error) {
		var nodesNotReady []*apiv1.Node
		nodeList, err := m.API.CoreV1().Nodes().List(ctx, opts)
		if err != nil {
			m.Logger.Debug("[WaitForCondition] Failed to get nodes from orchestrator: ", err)
			return false, nil
		}
		if len(nodeList.Items) == 0 {
			return false, nodes.ErrMissingNodes
		}
		for _, n := range nodeList.Items {
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
		condition, node.Selector(),
	))

	// Return new wait request with the created job
	return waiter.NewWaitRequest(job)
}

// isConditionSetAsExpected checks if the given Kubernetes Resource matches the expected condition.
func (m *kubernetesNodes) isConditionSetAsExpected(node apiv1.Node, expected resource.Condition) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == apiv1.NodeConditionType(expected.Type) &&
			cond.Status == apiv1.ConditionStatus(expected.Status) {
			return true
		}
	}
	return false
}

// NewNodes returns a nodes.Nodes implementation with the given kubernetes.Interface API.
func NewNodes(api kubernetes.Interface, logger ign.Logger) nodes.Nodes {
	return &kubernetesNodes{
		API:    api,
		Logger: logger,
	}
}
