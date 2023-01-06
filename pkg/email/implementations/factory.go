package implementations

import (
	sesfactory "github.com/gazebo-web/cloudsim/pkg/email/implementations/ses/factory"
	factorymap "github.com/gazebo-web/cloudsim/pkg/factory/map"
)

const (
	// SES is the AWS Simple Email Service (SES) implementation factory identifier.
	SES = "ses"
)

// Factory provides a factory to create Storage implementations.
var Factory = factorymap.Map{
	SES: sesfactory.NewFunc,
}
