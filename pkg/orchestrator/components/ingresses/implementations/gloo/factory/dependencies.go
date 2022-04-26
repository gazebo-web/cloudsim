package factory

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
)

// Dependencies is used to create a Gloo cluster component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`
	// Gloo is the Gloo clientset.
	Gloo gloo.GlooV1Interface
	// GlooGateway is the Gloo Gateway clientset.
	GlooGateway gateway.GatewayV1Interface
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
