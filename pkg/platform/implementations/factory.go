package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	platformfactory "github.com/gazebo-web/cloudsim/v4/pkg/platform/implementations/default/factory"
)

const (
	// Default is the default platform implementation.
	Default = "platform"
)

// Factory provides a factory to create Platform implementations.
var Factory = factorymap.Map{
	Default: platformfactory.NewFunc,
}
