package ingresses

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ingresses is an orchestrator.Ingresses implementation.
type ingresses struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Get returns an ingress with the given name.
func (m *ingresses) Get(name string, namespace string) (orchestrator.Resource, error) {
	m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s]", name, namespace))

	out, err := m.API.ExtensionsV1beta1().Ingresses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s] failed.", name, namespace))
		return nil, err
	}

	m.Logger.Debug(fmt.Sprintf("Getting ingress with name [%s] in namespace [%s] succeeded.", name, namespace))

	selector := orchestrator.NewSelector(out.Labels)
	return orchestrator.NewResource(name, namespace, selector), nil
}

// NewIngresses initializes a new orchestrator.Ingresses implementation using Kubernetes.
func NewIngresses(api kubernetes.Interface, logger ign.Logger) orchestrator.Ingresses {
	return &ingresses{
		API:    api,
		Logger: logger,
	}
}
