package kubernetes

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// kubernetesIngresses is an ingresses.Ingresses implementation.
type kubernetesIngresses struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Get returns an ingress with the given name.
func (m *kubernetesIngresses) Get(name string, namespace string) (resource.Resource, error) {
	m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s]", name, namespace))

	out, err := m.API.ExtensionsV1beta1().Ingresses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s] failed.", name, namespace))
		return nil, err
	}

	m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s] succeeded.", name, namespace))

	selector := resource.NewSelector(out.Labels)
	return resource.NewResource(name, namespace, selector), nil
}

// NewIngresses initializes a new ingresses.Ingresses implementation using Kubernetes.
func NewIngresses(api kubernetes.Interface, logger ign.Logger) ingresses.Ingresses {
	return &kubernetesIngresses{
		API:    api,
		Logger: logger,
	}
}
