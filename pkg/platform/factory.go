package platform

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	platformfactory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform/factory"
)

const (
	// Default is the default platform implementation.
	Default = "platform"
)

// Factory provides a factory to create Platform implementations.
var Factory = factorymap.Map{
	Default: platformfactory.NewFunc,
}
