package defaults

// Defaulter allows defining default values for an object.
// Note that because the methods in this interface directly modify an object's properties, it should only be
// implemented by objects that support pointer receivers.
type Defaulter interface {
	// SetDefaults sets the object's default values. This method must modify the object.
	SetDefaults() error
}

// SetDefaults sets default values for an object if it implements the Defaulter interface.
// Objects that do not implement the defaulter interface are not modified.
func SetDefaults(target interface{}) error {
	if defaulter, ok := target.(Defaulter); ok {
		if err := defaulter.SetDefaults(); err != nil {
			return err
		}
	}

	return nil
}
