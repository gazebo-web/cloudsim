package pods

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

// Pod is a Kubernetes pod. It extends the generic orchestrator.Resource interface.
type Pod interface {
	orchestrator.Resource
}

// pod is a Pod implementation that contains the basic information to identify a pod in a Kubernetes cluster.
type pod struct {
	name      string
	selector  string
	namespace string
}

// Name returns the name of the pod.
func (n pod) Name() string {
	return n.name
}

// Selector returns the selector of the pod.
func (n pod) Selector() string {
	return n.selector
}

// Namespace returns the namespace of the pod.
func (n pod) Namespace() string {
	return n.namespace
}

// NewPod returns a new Pod implementation using pod.
func NewPod(name string, namespace string, selector string) Pod {
	return &pod{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
