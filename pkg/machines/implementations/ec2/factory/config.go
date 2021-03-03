package factory

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"

// Config is used to create an EC2 machines component.
type Config struct {
	// Region is the region the EC2 component will operate in.
	Region string `validate:"required"`
}

// Validate validates that the config values are valid.
func (d *Config) Validate() error {
	return validate.DefaultStructValidator(d)
}
