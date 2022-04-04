package configurations

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
)

// CreateConfigurationInput is the input for creating a new configuration.
type CreateConfigurationInput struct {
	// Name is the name of the configuration.
	Name string

	// Namespace is the namespace where this configuration will be created.
	Namespace string

	// Labels is the group of key-value pairs that will identify this configuration.
	Labels map[string]string

	// Data contains the set of UTF-8 encoded data that will be stored in this configuration.
	Data map[string]string

	// Data contains the set of binary data that will be stored in this configuration.
	BinaryData map[string][]byte
}

// Configurations groups a set of methods to manage cluster configurations.
// Configurations are typically used to store configuration information for pods outside the pods themselves.
// Depending on the Cluster implementation, Configurations can be mounted directly in a pod when running it.
type Configurations interface {
	// Create creates a new configuration.
	Create(input CreateConfigurationInput) (resource.Resource, error)
	// Delete deletes a configuration.
	Delete(resource resource.Resource) (resource.Resource, error)
}
