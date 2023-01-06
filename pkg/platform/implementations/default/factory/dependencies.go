package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/validate"
	"github.com/gazebo-web/gz-go/v7"
)

// Dependencies is used to create an EC2 machines component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
