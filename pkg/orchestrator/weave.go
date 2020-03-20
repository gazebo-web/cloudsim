package orchestrator

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kc Kubernetes) WeaveRemovePeer(ctx context.Context, node *apiv1.Node) error {
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