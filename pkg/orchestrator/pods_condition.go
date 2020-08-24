package orchestrator

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
)

// PodCondition is a function type that returns the pod condition or error by the given k8s Pod.
type PodCondition func(ctx context.Context, pod *apiv1.Pod) (bool, error)

// podRunningAndReady checks if a pod by name is running. This function is used
// for Wait polls.
func podRunningAndReady(ctx context.Context, pod *apiv1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case apiv1.PodFailed, apiv1.PodSucceeded:
		return false, conditions.ErrPodCompleted
	case apiv1.PodRunning:
		return podutil.IsPodReady(pod), nil
	}
	return false, nil
}
