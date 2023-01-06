package factory

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/gazebo-web/cloudsim/pkg/validate"
	"github.com/gazebo-web/gz-go/v7"
)

// Dependencies is used to create an SES storage component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`
	// API is the SES API client used to interface with AWS SES.
	// If API is not provided, it will be initialized using Config values.
	API sesiface.SESAPI
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
