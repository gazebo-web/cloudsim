package factorymap

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
)

// Map is the default Factory implementation.
// It can be initialized using the `Register` method, by literal map initialization, or both.
type Map map[string]factory.NewFunc

// Verify that map implements factory.Factory
var _ factory.Factory = (Map)(nil)

// New sets the `out` parameter to a new instance of the specific object type.
// `out` should be a pointer to a value able to contain the expected type.
func (fm Map) New(cfg *factory.Config, dependencies factory.Dependencies, out interface{}) error {
	if cfg == nil {
		return factory.ErrorWithContext(factory.ErrNilConfig)
	}

	var fn factory.NewFunc
	var ok bool
	if fn, ok = fm[cfg.Type]; !ok {
		return factory.ErrorWithContext(factory.ErrFactoryTypeDoesNotExist)
	}

	// Create the object
	if err := fn(cfg.Config, dependencies, out); err != nil {
		return err
	}

	return nil
}

// Register registers a new object type the factory can create.
// Attempting to register an object type that already exists will result in an ErrFactoryTypeAlreadyExists error.
func (fm Map) Register(objectType string, fn factory.NewFunc) error {
	// Attempting to overwrite an existing type will result in an error
	if _, exists := fm[objectType]; exists {
		return factory.ErrFactoryTypeAlreadyExists
	}
	fm[objectType] = fn

	return nil
}

// NewMap creates and returns a new factory Map.
func NewMap() factory.Factory {
	return make(Map, 0)
}
