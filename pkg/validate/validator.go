package validate

import "gopkg.in/go-playground/validator.v9"

// Validator allows defining validation behavior for an object.
type Validator interface {
	// Validate validates that the object is valid.
	Validate() error
}

// Validate validates an object if it implements the Validator interface.
// Objects that do not implement the Validator interface are not validated.
func Validate(target interface{}) error {
	if defaulter, ok := target.(Validator); ok {
		if err := defaulter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// DefaultStructValidator validates that the passed structure is valid.
func DefaultStructValidator(target interface{}) error {
	return validator.New().Struct(target)
}
