package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/spdy"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/gazebo-web/gz-go/v7/validate"
	kubeapi "k8s.io/client-go/kubernetes"
)

// Dependencies is used to create a Kubernetes cluster component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`

	// API is the Kubernetes clientset.
	API kubeapi.Interface `validate:"required"`

	// SPDY is the SPDY executor initializer. It is required to run commands on pods.
	// If SPDY is not provided, a default SPDY executor will be created.
	SPDY spdy.Initializer
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
