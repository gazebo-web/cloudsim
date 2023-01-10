package factory

import "github.com/gazebo-web/cloudsim/v4/pkg/validate"

// Config is used to create an SES storage component.
type Config struct {
	// Region is the region the SES component will operate in.
	Region string `validate:"required"`
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}
