package kubernetes

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// kubernetesServices is a services.Services implementation.
type kubernetesServices struct {
	API    kubernetes.Interface
	Logger ign.Logger
}

// Create creates a new service defined by the given input.
func (s *kubernetesServices) Create(ctx context.Context, input services.CreateServiceInput) (resource.Resource, error) {
	s.Logger.Debug(fmt.Sprintf("Creating new Service. Input: %+v", input))

	// Create service port from input
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
	_, err := s.API.CoreV1().Services(input.Namespace).Create(ctx, newService, metav1.CreateOptions{})
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Creating new Service %s failed. Error: %+v", input.Name, err))
		return nil, err
	}

	s.Logger.Debug(fmt.Sprintf("Creating new Service %s succeeded.", input.Name))

	selector := resource.NewSelector(input.ServiceLabels)
	res := resource.NewResource(input.Name, input.Namespace, selector)

	return res, nil
}

func (s *kubernetesServices) Get(ctx context.Context, name string, namespace string) (resource.Resource, error) {
	s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s].", name, namespace))

	output, err := s.API.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s] failed. Error: %s", name, namespace, err))
		return nil, err
	}

	s.Logger.Debug(fmt.Sprintf("Getting service with name [%s] in namespace [%s] succeeded.", name, namespace))
	return resource.NewResource(name, namespace, resource.NewSelector(output.Labels)), nil
}

func (s *kubernetesServices) List(ctx context.Context, namespace string, selector resource.Selector) ([]resource.Resource, error) {
	s.Logger.Debug(fmt.Sprintf("Getting all services that match the following selectors: [%s]", selector.String()))

	list, err := s.API.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})

	if err != nil {
		s.Logger.Debug(fmt.Sprintf(
			"Getting all services matching selector: [%s] failed. Error: %s",
			selector.String(), err.Error()),
		)
		return nil, err
	}

	var output []resource.Resource

	for _, srv := range list.Items {
		selector := resource.NewSelector(srv.Labels)
		output = append(output, resource.NewResource(srv.Name, srv.Namespace, selector))
	}

	s.Logger.Debug(fmt.Sprintf(
		"Getting all services matching selector: [%s] succeeded. Output: %+v",
		selector.String(), output),
	)
	return output, nil
}

func (s *kubernetesServices) Remove(ctx context.Context, resource resource.Resource) error {
	s.Logger.Debug(fmt.Sprintf(
		"Removing service with name [%s] in namespace [%s].",
		resource.Name(), resource.Namespace()),
	)

	err := s.API.CoreV1().Services(resource.Namespace()).Delete(ctx, resource.Name(), metav1.DeleteOptions{})

	if err != nil {
		s.Logger.Debug(fmt.Sprintf(
			"Removing service with name [%s] in namespace [%s] failed. Error: %s",
			resource.Name(), resource.Namespace(), err.Error()),
		)
		return err
	}

	s.Logger.Debug(fmt.Sprintf(
		"Removing service with name [%s] in namespace [%s] succeeded.",
		resource.Name(), resource.Namespace()),
	)
	return nil
}

// NewServices initializes a new services.Services implementation using kubernetesServices.
func NewServices(api kubernetes.Interface, logger ign.Logger) services.Services {
	return &kubernetesServices{
		API:    api,
		Logger: logger,
	}
}
