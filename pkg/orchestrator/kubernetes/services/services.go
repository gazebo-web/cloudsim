package services

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
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
	s.Logger.Debug(fmt.Sprintf("Creating new Service %s succeded."), input.Name)
	return nil
}

// NewServices initializes a new orchestrator.Services implementation using services.
func NewServices(api kubernetes.Interface, logger ign.Logger) orchestrator.Services {
	return &services{
		API:    api,
		Logger: logger,
	}
}
