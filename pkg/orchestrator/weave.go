package orchestrator

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WeaveRemovePeer removes a node from the list of peers of the weave network.
// This step must be done manually before shutdown as per the weave documentation.
//   https://www.weave.works/docs/net/latest/operational-guide/tasks/#detecting-and-reclaiming-lost-ip-address-space
// Not doing so will make weave lose unrecoverable addresses to dead nodes.
func (kc *k8s) WeaveRemovePeer(ctx context.Context, node *apiv1.Node) error {
	// Get the weave pod name for the node
	nodeName := node.Name
	pods, err := kc.CoreV1().Pods(metav1.NamespaceSystem).List(metav1.ListOptions{
		LabelSelector: "name=weave-net",
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
		Limit:         1,
	})
	if err != nil || len(pods.Items) == 0 {
		msg := fmt.Sprintf("Failed to get weave pod for node %s.", nodeName)
		logger.Logger(ctx).Error(msg, err)
		return errors.New(msg)
	}
	weavePodName := pods.Items[0].Name

	// Remove the target node from the weave network
	command := []string{"sh", "-c", "'/home/weave/weave --local reset'"}
	if options, err := kc.PodExec(ctx, metav1.NamespaceSystem, weavePodName, "weave", command, nil); err != nil {
		msg := kc.PodCreateExecErrorMessage(
			fmt.Sprintf("Failed to remove node %s from weave network.", nodeName),
			options,
		)
		logger.Logger(ctx).Error(msg, err)
		return errors.New(msg)
	}

	return nil
}
