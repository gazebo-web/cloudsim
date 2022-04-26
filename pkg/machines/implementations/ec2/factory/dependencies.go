package factory

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/pricing"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
)

// Dependencies is used to create an EC2 machines component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger ign.Logger `validate:"required"`
	// API is the EC2 API client used to interface with AWS EC2 in a single region.
	// If API is not provided, it will be initialized using Config values.
	API ec2iface.EC2API
	// PricingAPI is the Pricing API client used to interface with AWS Pricing API.
	// If PricingAPI is not provided, it will be initialized using Config values.
	PricingAPI *pricing.Pricing
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
