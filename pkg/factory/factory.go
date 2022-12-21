package factory

import (
	errorsutils "github.com/gazebo-web/cloudsim/pkg/utils/errors"
	"github.com/gazebo-web/cloudsim/pkg/utils/reflect"
	"github.com/gazebo-web/cloudsim/pkg/validate"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	// ErrFactoryTypeAlreadyExists is returned when a creator function is registered for a type that already exists.
	ErrFactoryTypeAlreadyExists = errors.New("a creator function for the type already exists")
	// ErrFactoryTypeDoesNotExist is returned when a new object is requrested for an unregistered type.
	ErrFactoryTypeDoesNotExist = errors.New("factory type not found")
	// ErrNilConfig is returned when a nil config is passed to a factory.
	ErrNilConfig = errors.New("factory config is nil")
	// ErrMissingDependency is returned when a creator function fails to find a required dependency.
	ErrMissingDependency = errors.New("creator function did not receive required dependency")
	// ErrInvalidDependency is returned when a creator function finds an invalid dependency (e.g. wrong type, not
	// initialized, etc.).
	ErrInvalidDependency = errors.New("creator function received invalid dependency")
)

// NewFunc is a factory creation function used by a factory to instance objects.
// Creation functions are in charge of processing the provided config, gathering dependencies, initializing any
// additional subcomponents or dependencies and returning a new object of the expected type.
type NewFunc func(config interface{}, dependencies Dependencies, out interface{}) error

// Factory instances objects of a specific type with a given configuration.
// A factory works by mapping an object type to a factory creation function (NewFunc).
//
// # Example
//
// The `example` package has an `Exampler` interface, and three different implementations of it.
// A package `Factory` is defined to allow other packages to dynamically request `Exampler` implementations.
//
// In order to allow the package factory to create `Exampler` implementations, each implementation must define and
// register a factory creation function (NewFunc) in the package factory.
//
// Define a factory creation function for each implementation that will be registered in the factory.
//
// example/implementations/type1/example.go
// ```
//
//	func newType1(config interface{}, dependencies factory.Dependencies, out interface{}) error {
//	    // It is recommended to define a Config struct for your factory creation function.
//	    var typeConfig *Config
//	    if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
//	        return err
//	    }
//
//	    // Parse dependencies
//	    // It is recommended to define a Config struct for your factory creation function.
//	    var typeDependencies Dependencies
//	    if err := dependencies.ToStruct(&typeDependencies); err != nil {
//	        return err
//	    }
//
//	    // Initialize the implementation
//	    t1 := [...]
//
//	    // Set the the output value to the newly created instance
//		   factory.SetValue(out, t1)
//	}
//
// ```
//
// Create a package factory and register all implementations.
//
// example/factory.go
// ```
// package example
//
// const (
//
//	// Implementation types
//	Type1 = "type1"
//	Type2 = "type2"
//	Type3 = "type3"
//
// )
//
//	var Factory := factorymap.Map{
//	    Type1: type1New,
//	    Type2: type2New,
//	    Type3: type3New,
//	}
//
// ```
//
// The package factory can now be used by other packages that need `Exampler`s.
//
// consumer/consumer.go
// ```
//
//	func Example() error {
//	    // Prepare the factory config
//	    // This config is loaded manually here, but it can also be loaded from a file
//	    config := &factory.Config{
//	         Type: example.Type1,
//	         Config: [...],
//	    }
//
//	    // Prepare factory dependencies
//	    dependencies := Dependencies{ [...] },
//
//	    // Create a new `Exampler` object
//	    var exampler example.Exampler
//	    err := example.Factory.New(config, dependencies, &exampler)
//	    if err != nil {
//	        return err
//	    }
//	}
//
// ```
type Factory interface {
	// New sets the `out` parameter to a new instance of the specific object type.
	// `out` should be a pointer to a value able to contain the expected type.
	New(cfg *Config, dependencies Dependencies, out interface{}) error
	// Register registers a new object type the factory can create.
	Register(objectType string, fn NewFunc) error
}

// SetValue sets the `out` parameter to the specified value.
// It should be called by factory creation functions to set `out`'s output value.
// `out` must be a pointer to a value able to contain the expected type.
func SetValue(out interface{}, value interface{}) error {
	switch value.(type) {
	case map[string]interface{}, ConfigValues:
		return mapstructure.Decode(value, out)
	default:
		return reflect.SetValue(out, value)
	}
}

// SetValueAndValidate sets the `out` parameter to the specified value.
// The `out` value is validated if it implements validate.Validator.
func SetValueAndValidate(out interface{}, value interface{}) error {
	// Set value
	if err := SetValue(out, value); err != nil {
		return err
	}
	// Validate
	if err := validate.Validate(out); err != nil {
		return err
	}

	return nil
}

// ErrorWithContext wraps an error with information about the function that generated the error.
func ErrorWithContext(err error) error {
	errMsg := "factory function failed to create value"
	return errorsutils.WithFunctionContext(err, errMsg, 2)
}
