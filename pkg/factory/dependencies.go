package factory

import (
	"github.com/gazebo-web/gz-go/v7/validate"
	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"reflect"
)

// Dependencies is a container for a set of dependencies. A dependency is an object or value that is not initialized
// by the factory and cannot be obtained from the config.
// It is recommended that a received Dependencies object is marshalled onto a concrete struct before being used. Read
// the ToStruct documentation for an in-depth explanation.
type Dependencies map[string]interface{}

// Set sets a dependency.
func (d Dependencies) Set(key string, value interface{}) {
	d[key] = value
}

// Get gets a single dependency.
// It is recommended to use the `ToStruct` method instead, as it allows getting and validating dependencies in a single
// operation.
// `out` should be a pointer.
func (d Dependencies) Get(key string, out interface{}) (err error) {
	// Look for the value in the dependency map
	value, ok := d[key]
	if !ok {
		return ErrMissingDependency
	}

	// Assign the dependency value
	// Handle panics by returning an error if the value is not compatible with the out type
	defer func() {
		if r := recover(); r != nil {
			err = ErrInvalidDependency
		}
	}()
	v := reflect.ValueOf(out)
	v.Elem().Set(reflect.ValueOf(value))

	return nil
}

// marshal marshals the dependency map into a struct.
func (d Dependencies) marshal(out interface{}) error {
	return mapstructure.Decode(d, out)
}

// validate validates a marshalled map.
func (d Dependencies) validate(out interface{}) error {
	if err := validate.Validate(out); err != nil {
		return errors.Wrap(ErrInvalidDependency, err.Error())
	}
	return nil
}

// ToStruct marshals this dependency map into a struct.
// A factory cannot know the required dependencies needed to create a specific object type. In order to allow an object
// type creation function to verify that it is being passed all the dependencies it requires, it can define a struct
// that specifies its required dependencies, and call this method on the Dependencies object it receives.
// `out` must be a pointer to a struct. If `out` implements the validate.Validator interface, the populated struct will
// be automatically validated.
//
// Example
// ```
// // ObjectDependencies defines the dependencies required to create Object using a factory creation function.
//
//	type ObjectDependencies struct {
//	    gz.Logger `validate:"required"`
//	}
//
// // Defining this method will automatically validate that the received dependencies match the expected values
//
//	func (od *ObjectDependencies) Validate() error {
//	    return validate.DefaultStructValidator(od)
//	}
//
// [...]
//
// // Factory creation function
//
//	func([...], dependencies factory.Dependencies) [...] {
//	    objectDependencies := &ObjectDependencies{}
//	    // The dependencies.ToStruct call will populate objectDependencies.
//	    // It will also validate that all required fields are in place because ObjectDependencies has a Validate method
//	    // defined.
//	    if err := dependencies.ToStruct(objectDependencies); err != nil {
//	         [...] // Handle error
//	    }
//	    [...]
//	}
//
// ```
func (d Dependencies) ToStruct(out interface{}) error {
	if err := d.marshal(out); err != nil {
		return err
	}
	if err := d.validate(out); err != nil {
		return err
	}

	return nil
}

// DeepCopy returns a deep copy of a Dependencies object.
// It should be called before modifying dependencies to avoid interfering with dependencies in other factories.
func (d Dependencies) DeepCopy() (Dependencies, error) {
	// Copy the value
	out, err := copystructure.Copy(d)
	if err != nil {
		return nil, err
	}

	// Cast to Dependencies
	dc, ok := out.(Dependencies)
	if !ok {
		return nil, errors.New("deep copied value is not of type Dependencies")
	}

	return dc, nil
}

// NewDependencies returns a new Dependencies object.
func NewDependencies() Dependencies {
	return make(Dependencies, 0)
}
