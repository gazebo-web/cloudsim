package pods

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
)

// manager is a orchestrator.PodManager implementation.
type manager struct {
	API  kubernetes.Interface
	SPDY spdy.Initializer
}

// Exec creates a new executor.
func (m *manager) Exec(pod orchestrator.Resource) orchestrator.Executor {
	return newExecutor(m.API, pod, m.SPDY)
}

// Reader creates a new reader.
func (m *manager) Reader(pod orchestrator.Resource) orchestrator.Reader {
	return newReader(m.API, pod, m.SPDY)
}

// WaitForCondition creates a new wait request.
func (m *manager) WaitForCondition(pod orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	opts := metav1.ListOptions{
		LabelSelector: pod.Selector(),
	}
	var podsNotReady []*apiv1.Pod
	job := func() (bool, error) {
		podsNotReady = nil
		pods, err := m.API.CoreV1().Pods(pod.Namespace()).List(opts)
		if err != nil {
			return false, err
		}
		for _, p := range pods.Items {
			if condition == orchestrator.ReadyCondition {
				ready, err := m.isPodReady(&p)
				if err != nil {
					return false, err
				}
				if !ready {
					pod := new(apiv1.Pod)
					*pod = p
					podsNotReady = append(podsNotReady, pod)
				}
			}
		}
		return len(podsNotReady) == 0, nil
	}
	return waiter.NewWaitRequest(job)
}

// isPodReady checks if the given Kubernetes Pod matches the ready condition.
func (m *manager) isPodReady(pod *apiv1.Pod) (bool, error) {
	if pod.Status.Phase == apiv1.PodFailed || pod.Status.Phase == apiv1.PodSucceeded {
		return false, conditions.ErrPodCompleted
	}
	return podutil.IsPodReady(pod), nil
}

// NewManager initializes a new manager.
func NewManager(api kubernetes.Interface, spdy spdy.Initializer) orchestrator.PodManager {
	return &manager{
		API:  api,
		SPDY: spdy,
	}
}
