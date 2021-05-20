package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
)

// Config is used to create an EC2 machines component.
type Config struct {
	// Region is the region the EC2 component will operate in.
	Region string `validate:"required"`
	// Limit is the maximum number of EC2 instances that this component has available.
	// If set to -1, it will not limit the number of instances. Note that the component will still be subject to EC2
	// instance availability.
	Limit *int64
	// WorkerGroupName is the label value set on all machines created by this component. It is used to identify
	// machines created by this component.
	WorkerGroupName string
	// Zones contains the set of availability zones the machines component will launch simulation instances in.
	Zones []ec2.Zone `validate:"required"`
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}