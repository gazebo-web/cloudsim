package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// manager is an orchestrator.IngressManager implementation.
type manager struct {
	API kubernetes.Interface
}

// Get returns an Ingress with the given
func (m *manager) Get(name string, namespace string) (orchestrator.Resource, error) {
	_, err := m.API.ExtensionsV1beta1().Ingresses(name).Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return NewIngress(name, namespace), nil
}

// NewManager initializes a new orchestrator.IngressManager implementation using Kubernetes.
func NewManager(api kubernetes.Interface) orchestrator.IngressManager {
	return &manager{
		API: api,
	}
}
