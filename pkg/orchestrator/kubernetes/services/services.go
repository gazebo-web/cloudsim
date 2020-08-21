package services

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/types"
	"gitlab.com/ignitionrobotics/web/ign-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// services is a orchestrator.Services implementation.
type services struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Create creates a new service defined by the given input.
func (s *services) Create(input orchestrator.CreateServiceInput) error {
	s.Logger.Debug(fmt.Sprintf("Creating new Service. Input: %+v", input))
	var ports []corev1.ServicePort
	for key, value := range input.Ports {
		ports = append(ports, corev1.ServicePort{Name: key, Port: value})
	}

	// Prepare the resource
	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
			Labels:    input.ServiceLabels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceType(input.Type),
			Selector: input.TargetLabels,
			Ports:    ports,
		},
	}

	// Launch the resource
	_, err := s.API.CoreV1().Services(input.Namespace).Create(newService)
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Creating new Service %s failed. Error: %+v", input.Name, err))
		return err
	}
	s.Logger.Debug(fmt.Sprintf("Creating new Service %s succeded.", input.Name))
	return nil
}

func (s *services) Get(name, namespace string) (orchestrator.Resource, error) {
	s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s].", name, namespace))
	output, err := s.API.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s] failed. Error: %s", name, namespace, err))
		return nil, err
	}
	s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s] succeded.", name, namespace))
	return types.NewResource(name, namespace, types.NewSelector(output.Labels)), nil
}

func (s *services) GetAllBySelector(namespace string, selector orchestrator.Selector) ([]orchestrator.Resource, error) {
	s.Logger.Debug(fmt.Sprintf("Getting all services that match the following selectors: [%s]", selector.String()))
	list, err := s.API.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		s.Logger.Debug(fmt.Sprintf(
			"Getting all services matching selector: [%s] failed. Error: %s",
			selector.String(), err.Error()),
		)
		return nil, err
	}
	var output []orchestrator.Resource
	for _, srv := range list.Items {
		selector := types.NewSelector(srv.Labels)
		output = append(output, types.NewResource(srv.Name, srv.Namespace, selector))
	}
	s.Logger.Debug(fmt.Sprintf(
		"Getting all services matching selector: [%s] succeeded. Output: %+v",
		selector.String(), output),
	)
	return output, nil
}

func (s *services) Remove(resource orchestrator.Resource) error {
	s.Logger.Debug(fmt.Sprintf(
		"Removing service with name [%s] in namespace [%s] that match the following selectors: [%s].",
		resource.Name(), resource.Namespace(), resource.Selector().String()),
	)
	err := s.API.CoreV1().Services(resource.Namespace()).Delete(resource.Name(), &metav1.DeleteOptions{})
	if err != nil {
		s.Logger.Debug(fmt.Sprintf(
			"Removing service with name [%s] in namespace [%s] that match the following selectors: [%s] failed. Error: %s",
			resource.Name(), resource.Namespace(), resource.Selector().String(), err.Error()),
		)
		return err
	}
	s.Logger.Debug(fmt.Sprintf(
		"Removing service with name [%s] in namespace [%s] that match the following selectors: [%s] succeeded.",
		resource.Name(), resource.Namespace(), resource.Selector().String()),
	)
	return nil
}

// NewServices initializes a new orchestrator.Services implementation using services.
func NewServices(api kubernetes.Interface, logger ign.Logger) orchestrator.Services {
	return &services{
		API:    api,
		Logger: logger,
	}
}
