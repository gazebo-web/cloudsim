package orchestrator

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func (kc Kubernetes) PodWaitForReadyCondition(ctx context.Context, c kubernetes.Interface, namespace string, groupIDLabel string, timeout time.Duration) error {
	opts := metav1.ListOptions{LabelSelector: groupIDLabel}
	return kc.PodWaitToMatchCondition(ctx, namespace, opts, "Ready", timeout, podRunningAndReady)
}

func (kc Kubernetes) PodWaitToMatchCondition(ctx context.Context, namespace string, opts metav1.ListOptions, condStr string, timeout time.Duration, condition PodCondition) error {
	logger.Logger(ctx).Info(fmt.Sprintf("Waiting up to %v for matching pods' status to be %s", timeout, condStr))
	for start := time.Now(); time.Since(start) < timeout; tools.Sleep(pollFrequency) {
		pods, err := kc.CoreV1().Pods(namespace).List(opts)
		if err != nil {
			return err
		}
		var conditionNotMatch []string
		for _, pod := range pods.Items {
			done, err := condition(ctx, &pod)
			if done && err != nil {
				return fmt.Errorf("unexpected error: %v", err)
			}
			if !done {
				conditionNotMatch = append(conditionNotMatch, pod.Name)
			}
		}
		if len(conditionNotMatch) <= 0 {
			logger.Logger(ctx).Info(fmt.Sprintf("Success. Pods match condition '%s'", condStr))
			return err
		}
		logger.Logger(ctx).Debug(fmt.Sprintf("%d pods are not %s: %v", len(conditionNotMatch), condStr, conditionNotMatch))
	}
	return errors.Errorf("gave up waiting for matching pods to be '%s' after %v", condStr, timeout)
}
