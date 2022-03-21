package reflect

import (
	"errors"
	"reflect"
)

var (
	// ErrNotPointer is returned when a function receives a parameter that is not a pointer.
	ErrNotPointer = errors.New("out value is not pointer")
	// ErrNotCollection is returned when a function receives a parameter that is not a collection.
	ErrNotCollection = errors.New("out value is not collection")
	// ErrInvalidOutValue is returned when an out value cannot be set to the target value.
	ErrInvalidOutValue = errors.New("invalid out value type")
)

// SetValue sets the `out` parameter to the specified value.
// `out` must be a pointer.
func SetValue(out interface{}, value interface{}) (err error) {
	// Handle panics
	defer func() {
		if p := recover(); p != nil {
			err = ErrInvalidOutValue
		}
	}()

	// Get the pointer value
	p := reflect.ValueOf(out)
	if p.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	v := reflect.ValueOf(value)
	pv := p.Elem()
	pv.Set(v)

	return nil
}

// AppendToSlice appends values to a slice.
// `out` must be a pointer to a slice.
// `values` must be compatible with the slice type.
func AppendToSlice(out interface{}, values ...interface{}) error {
	// Get the pointer value
	p := reflect.ValueOf(out)
	if p.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	// Get the slice value
	v := p.Elem()
	if v.Kind() != reflect.Slice {
		return ErrNotCollection
	}

	// Append values to slice
	for _, value := range values {
		v.Set(reflect.Append(v, reflect.ValueOf(value)))
	}

	return nil
}

// SetMapValue sets a value in a map.
// `out` must be a map.
// `key` must be compatible with the map key type.
// `value` must be compatible with the map value type.
func SetMapValue(out interface{}, key interface{}, value interface{}) error {
	// Get the map value
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Map {
		return ErrNotCollection
	}

	// Set map value
	v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))

	return nil
}

// NewInstance returns a new instance of the same type as the input value.
// The returned value will contain the zero value of the type.
func NewInstance(value interface{}) interface{} {
	entity := reflect.ValueOf(value)

	if entity.Kind() == reflect.Ptr {
		entity = reflect.New(entity.Elem().Type())
	} else {
		entity = reflect.New(entity.Type()).Elem()
	}

	return entity.Interface()
}

// NewCollectionValueInstance receives a collection (slice, map) and returns a new instance of the collection's value
// type. If the collection value type is a pointer type, a pointer object to a new instance of the type value is
// returned.
func NewCollectionValueInstance(collection interface{}) (interface{}, error) {
	// Get the collection value
	s := reflect.TypeOf(collection)
	if s.Kind() != reflect.Slice && s.Kind() != reflect.Map {
		return nil, ErrNotCollection
	}

	// Get the collection's value type
	t := s.Elem()

	// Create the collection element instance
	var v reflect.Value
	// If the value is a pointer, assign a value of the correct type
	if t.Kind() == reflect.Ptr {
		// Pointer value instance value
		v = reflect.New(t.Elem())
	} else {
		v = reflect.New(t).Elem()
	}

	return v.Interface(), nil
}
