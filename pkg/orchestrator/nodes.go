package orchestrator

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// NodeWaitForReady function waits (blocks) until the identified nodes are ready, or until there
// is a timeout or error.
func (kc Kubernetes) NodeWaitForReady(ctx context.Context, namespace string, groupIDLabel string, timeout time.Duration) error {
	opts := metav1.ListOptions{LabelSelector: groupIDLabel}
	return kc.NodeWaitToMatchCondition(ctx, namespace, opts, timeout)
}

// NodeWaitToMatchCondition finds match Nodes based on the input ListOptions.
// Waits and checks if all matched nodes are in the given PodCondition
func (kc Kubernetes) NodeWaitToMatchCondition(ctx context.Context, namespace string, opts metav1.ListOptions, timeout time.Duration) error {
	logger.Logger(ctx).Info(fmt.Sprintf("Waiting up to %v for match nodes to be ready", timeout))

	maxAllowedNotReadyNodes := 0
	var notReady []*apiv1.Node
	err := wait.PollImmediate(pollFrequency, timeout, func() (bool, error) {
		notReady = nil
		// It should be OK to list unschedulable Nodes here.
		nodes, err := kc.CoreV1().Nodes().List(opts)
		if err != nil {
			return false, err
		}
		logger.Logger(ctx).Debug(fmt.Sprintf("Found nodes %v", nodes.Items))
		for i := range nodes.Items {
			node := &nodes.Items[i]
			if !isNodeConditionSetAsExpected(ctx, node, apiv1.NodeReady, true) {
				notReady = append(notReady, node)
			}
		}
		return len(notReady) <= maxAllowedNotReadyNodes, nil
	})

	if err != nil && err != wait.ErrWaitTimeout {
		return err
	}

	if len(notReady) > maxAllowedNotReadyNodes {
		msg := ""
		for _, node := range notReady {
			msg = fmt.Sprintf("%s, %s", msg, node.Name)
		}
		return errors.Errorf("Nodes not ready: %#v", msg)
	}
	return nil
}

// isNodeConditionSetAsExpected checks if the given condition is met by a node.
func isNodeConditionSetAsExpected(ctx context.Context, node *apiv1.Node, conditionType apiv1.NodeConditionType, wantTrue bool) bool {
	// Check the node readiness condition (logging all).
	for _, cond := range node.Status.Conditions {
		// Ensure that the condition type and the status matches as desired.
		if cond.Type == conditionType {
			if cond.Type == apiv1.NodeReady {
				// For NodeReady we should check if Taints are gone as well
				// See taint.MatchTaint
				if wantTrue {
					if cond.Status == apiv1.ConditionTrue {
						return true
					}
					msg := fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
						conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message)
					logger.Logger(ctx).Debug(msg)
					return false

				}
				// TODO: check if the Node is tainted once we enable NC notReady/unreachable taints by default
				if cond.Status != apiv1.ConditionTrue {
					return true
				}
				logger.Logger(ctx).Debug(fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
					conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message))
				return false

			}
			if (wantTrue && (cond.Status == apiv1.ConditionTrue)) || (!wantTrue && (cond.Status != apiv1.ConditionTrue)) {
				return true
			}
			logger.Logger(ctx).Debug(fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
				conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message))
			return false
		}
	}

	logger.Logger(ctx).Debug(fmt.Sprintf("Couldn't find condition %v on node %v", conditionType, node.Name))
	return false
}
