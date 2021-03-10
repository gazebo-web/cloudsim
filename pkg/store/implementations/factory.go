package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/store"
)

const (
	// Store is the default Store implementation factory identifier.
	Store = "store"
)

// Factory provides a factory to create Cluster implementations.
var Factory = factorymap.Map{
	Store: store.NewFunc,
}
