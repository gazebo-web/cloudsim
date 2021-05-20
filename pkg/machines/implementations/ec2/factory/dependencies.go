package factory

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Dependencies is used to create an EC2 machines component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`
	// API is the EC2 API client used to interface with AWS EC2 in a single region.
	// If API is not provided, it will be initialized using Config values.
	API ec2iface.EC2API
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
