package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ingresses is an orchestrator.Ingresses implementation.
type ingresses struct {
	API kubernetes.Interface
}

// Get returns an Resource with the given
func (m *ingresses) Get(name string, namespace string) (orchestrator.Resource, error) {
	_, err := m.API.ExtensionsV1beta1().Ingresses(name).Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return NewIngress(name, namespace), nil
}

// NewIngresses initializes a new orchestrator.Ingresses implementation using Kubernetes.
func NewIngresses(api kubernetes.Interface) orchestrator.Ingresses {
	return &ingresses{
		API: api,
	}
}
