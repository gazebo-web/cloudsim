package factory

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/gazebo-web/gz-go/v7/validate"
)

// Dependencies is used to create an S3 storage component.
type Dependencies struct {
	// Logger is used to store log information.
	Logger gz.Logger `validate:"required"`
	// API is the S3 API client used to interface with AWS S3.
	// If API is not provided, it will be initialized using Config values.
	API s3iface.S3API
}

// Validate validates that the dependencies values are valid.
func (d *Dependencies) Validate() error {
	return validate.DefaultStructValidator(d)
}
