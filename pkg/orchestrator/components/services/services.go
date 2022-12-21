package services

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
)

// CreateServiceInput is used as an input of Services.Create method.
// It has all the information needed to create a new service.
type CreateServiceInput struct {
	// Name represents the name of the service.
	Name string

	// Type is the type of the service that's being created.
	// Kubernetes: ClusterIP or LoadBalancer.
	Type string

	// Namespace is the namespace where the service will live in.
	Namespace string

	// ServiceLabels are the unique set of key-value pairs that will define this service.
	ServiceLabels map[string]string

	// TargetLabels are the unique set of key-value pairs that the service will be pointed to.
	TargetLabels map[string]string

	// Ports describes the name and the port number that are going to be exposed by the created service.
	Ports map[string]int32
}

// Services groups a set of methods for managing services like Load Balancers.
// services are usually used to abstract a group of pods behind a single endpoint.
type Services interface {
	Create(ctx context.Context, input CreateServiceInput) (resource.Resource, error)
	Get(ctx context.Context, name string, namespace string) (resource.Resource, error)
	List(ctx context.Context, namespace string, selector resource.Selector) ([]resource.Resource, error)
	Remove(ctx context.Context, resource resource.Resource) error
}
