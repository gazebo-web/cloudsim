package factory

import (
	"github.com/gazebo-web/gz-go/v7"
	"github.com/gazebo-web/gz-go/v7/validate"
	kubeapi "k8s.io/client-go/kubernetes"
)

// Dependencies is used to create a Kubernetes cluster component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`

	// API is the Kubernetes clientset.
	// If API is not provided, an API instance will be created using the configuration defined in the Config object.
	API kubeapi.Interface
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
