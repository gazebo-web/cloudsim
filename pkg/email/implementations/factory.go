package implementations

import (
	sesfactory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/email/implementations/ses/factory"
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
)

const (
	// SES is the AWS Simple Email Service (SES) implementation factory identifier.
	SES = "ses"
)

// Factory provides a factory to create Storage implementations.
var Factory = factorymap.Map{
	SES: sesfactory.NewFunc,
}
