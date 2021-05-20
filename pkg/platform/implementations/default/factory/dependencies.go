package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Dependencies is used to create an EC2 machines component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
