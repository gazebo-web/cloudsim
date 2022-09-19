package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	kubeapi "k8s.io/client-go/kubernetes"
)

// Dependencies is used to create a Kubernetes cluster component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`

	// API is the Kubernetes clientset.
	// If API is not provided, an API instance will be created using the configuration defined in the Config object.
	API kubeapi.Interface
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
