package factory

import "github.com/gazebo-web/cloudsim/pkg/validate"

// Config is used to create an S3 storage component.
type Config struct {
	// Region is the region the S3 component will operate in.
	Region string `validate:"required"`
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}
