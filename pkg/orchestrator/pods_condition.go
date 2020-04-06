package orchestrator

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
)

type PodCondition func(ctx context.Context, pod *apiv1.Pod) (bool, error)

func podRunningAndReady(ctx context.Context, pod *apiv1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case apiv1.PodFailed, apiv1.PodSucceeded:
		return false, conditions.ErrPodCompleted
	case apiv1.PodRunning:
		return podutil.IsPodReady(pod), nil
	}
	return false, nil
}
