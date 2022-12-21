package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/validate"
	"github.com/gazebo-web/gz-go/v7"
	kubeapi "k8s.io/client-go/kubernetes"
)

// Dependencies are used to create a Kubernetes cluster component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`

	// API is the Kubernetes clientset.
	API kubeapi.Interface `validate:"required"`
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
