package factory

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Dependencies is used to create an SES storage component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`
	// API is the SES API client used to interface with AWS SES.
	// If API is not provided, it will be initialized using Config values.
	API sesiface.SESAPI
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
