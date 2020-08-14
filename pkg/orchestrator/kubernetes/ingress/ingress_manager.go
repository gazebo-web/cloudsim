package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// manager is an orchestrator.IngressManager implementation.
type manager struct {
	API kubernetes.Interface
}

// get returns an ingress from the given resource.
func (m manager) get(resource orchestrator.Resource) (*v1beta1.Ingress, error) {
	return m.API.ExtensionsV1beta1().Ingresses(resource.Namespace()).Get(resource.Name(), metav1.GetOptions{})
}

// update updates the given ingress in the namespace declared by the given resource.
func (m manager) update(resource orchestrator.Resource, ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return m.API.ExtensionsV1beta1().Ingresses(resource.Namespace()).Update(ingress)
}

// Rules returns a new orchestrator.Ruler implementation configured with the given ingress.
func (m *manager) Rules(ingress orchestrator.Resource) orchestrator.Ruler {
	return &ruler{
		resource: ingress,
		manager:  m,
	}
}

// NewManager initializes a new orchestrator.IngressManager implementation using Kubernetes.
func NewManager(api kubernetes.Interface) orchestrator.IngressManager {
	return &manager{
		API: api,
	}
}
