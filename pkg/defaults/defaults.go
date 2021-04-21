package defaults

import (
	"errors"
	"fmt"
	"github.com/creasty/defaults"
)

// Defaulter allows defining default values for an object.
// Note that because the methods in this interface directly modify an object's properties, it should only be
// implemented by objects that support pointer receivers.
type Defaulter interface {
	// SetDefaults sets the object's default values. This method must modify the object.
	SetDefaults() error
}

// SetValues sets default values for an object if it implements the Defaulter interface.
// Objects that do not implement the defaulter interface are not modified.
func SetValues(target interface{}) (err error) {
	if defaulter, ok := target.(Defaulter); ok {
		// Handle potential panics
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprint("defaults.SetValues panicked when attempting to set default values: ", r)
				err = errors.New(errMsg)
			}
		}()
		// Set defaults
		if err := defaulter.SetDefaults(); err != nil {
			return err
		}
	}

	return nil
}

// SetStructValues sets struct field default values defined in `default` tags.
func SetStructValues(target interface{}) error {
	return defaults.Set(target)
}
