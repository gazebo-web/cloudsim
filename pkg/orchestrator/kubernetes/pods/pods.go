package pods

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
)

// pods is a orchestrator.Pods implementation.
type pods struct {
	API    kubernetes.Interface
	SPDY   spdy.Initializer
	Logger ign.Logger
}

// Exec creates a new executor.
func (p *pods) Exec(pod orchestrator.Resource) orchestrator.Executor {
	p.Logger.Debug(fmt.Sprintf("Creating new executor for pod [%s]", pod.Name()))
	return newExecutor(p.API, pod, p.SPDY, p.Logger)
}

// Reader creates a new reader.
func (p *pods) Reader(pod orchestrator.Resource) orchestrator.Reader {
	p.Logger.Debug(fmt.Sprintf("Creating new reader for pod [%s]", pod.Name()))
	return newReader(p.API, pod, p.SPDY, p.Logger)
}

// WaitForCondition creates a new wait request that will be used to wait for a resource to match a certain condition.
// The wait request won't be triggered until the method Wait has been called.
func (p *pods) WaitForCondition(resource orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	opts := metav1.ListOptions{
		LabelSelector: resource.Selector().String(),
	}
	var podsNotReady []*apiv1.Pod
	job := func() (bool, error) {
		podsNotReady = nil
		po, err := p.API.CoreV1().Pods(resource.Namespace()).List(opts)
		if err != nil {
			return false, err
		}
		for _, i := range po.Items {
			if condition == orchestrator.ReadyCondition {
				ready, err := p.isPodReady(&i)
				if err != nil {
					return false, err
				}
				if !ready {
					pod := new(apiv1.Pod)
					*pod = i
					podsNotReady = append(podsNotReady, pod)
				}
			}
		}
		return len(podsNotReady) == 0, nil
	}
	return waiter.NewWaitRequest(job)
}

// isPodReady checks if the given Kubernetes Pod matches the ready condition.
func (p *pods) isPodReady(pod *apiv1.Pod) (bool, error) {
	if pod.Status.Phase == apiv1.PodFailed || pod.Status.Phase == apiv1.PodSucceeded {
		return false, conditions.ErrPodCompleted
	}
	return podutil.IsPodReady(pod), nil
}

// NewPods initializes a new orchestrator.Pods implementation for managing Kubernetes Pods.
func NewPods(api kubernetes.Interface, spdy spdy.Initializer, logger ign.Logger) orchestrator.Pods {
	return &pods{
		API:    api,
		SPDY:   spdy,
		Logger: logger,
	}
}
