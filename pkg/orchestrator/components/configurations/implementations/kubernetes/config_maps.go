package kubernetes

import (
	"context"
	"fmt"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
	"github.com/gazebo-web/gz-go/v7"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// configMaps is a configurations.Configurations Kubernetes implementation.
type configMaps struct {
	API    kubernetes.Interface
	Logger gz.Logger
}

// Create creates a config map.
func (cm *configMaps) Create(ctx context.Context, input configurations.CreateConfigurationInput) (resource.Resource, error) {
	// Prepare input for Kubernetes
	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Data:       input.Data,
		BinaryData: input.BinaryData,
	}

	cm.Logger.Debug(
		fmt.Sprintf(
			"Creating config map with name [%s] in namespace [%s]",
			input.Name,
			input.Namespace,
		),
	)

	// Create config map
	_, err := cm.API.CoreV1().ConfigMaps(input.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		cm.Logger.Debug(
			fmt.Sprintf(
				"Creating config map with name [%s] in namespace [%s] failed. Error: %s",
				input.Name,
				input.Namespace,
				err,
			),
		)
		return nil, err
	}

	cm.Logger.Debug(
		fmt.Sprintf(
			"Creating config map with name [%s] in namespace [%s] succeeded",
			input.Name,
			input.Namespace,
		),
	)
	return resource.NewResource(input.Name, input.Namespace, resource.NewSelector(input.Labels)), nil
}

// Delete removes a config map with the given name in the given namespace.
func (cm *configMaps) Delete(ctx context.Context, resource resource.Resource) (resource.Resource, error) {
	cm.Logger.Debug(
		fmt.Sprintf("Deleting pod with name [%s] in namespace [%s]", resource.Name(), resource.Namespace()),
	)

	err := cm.API.CoreV1().ConfigMaps(resource.Namespace()).Delete(ctx, resource.Name(), metav1.DeleteOptions{})
	if err != nil {
		cm.Logger.Debug(fmt.Sprintf(
			"Deleting pod with name [%s] in namespace [%s] failed. Error: %+v.",
			resource.Name(), resource.Namespace(), err,
		))
		return nil, err
	}

	cm.Logger.Debug(fmt.Sprintf(
		"Deleting pod with name [%s] in namespace [%s] succeeded.",
		resource.Name(), resource.Namespace(),
	))

	return resource, nil
}

// NewConfigMaps initializes a new configurations.Configurations Kubernetes implementation.
func NewConfigMaps(api kubernetes.Interface, logger gz.Logger) configurations.Configurations {
	return &configMaps{
		API:    api,
		Logger: logger,
	}
}
