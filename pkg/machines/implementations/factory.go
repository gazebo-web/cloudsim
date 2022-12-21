package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/pkg/factory/map"
	ec2factory "github.com/gazebo-web/cloudsim/pkg/machines/implementations/ec2/factory"
)

const (
	// EC2 is the EC2 implementation factory identifier.
	EC2 = "ec2"
)

// Factory provides a factory to create Machines implementations.
var Factory = factorymap.Map{
	EC2: ec2factory.NewFunc,
}
