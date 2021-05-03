package cycler

import (
    "github.com/pkg/errors"
    "reflect"
)

var (
    // ErrNotSlice is returned when NewCyclerFromSlice was not passed a input of type slice.
    ErrNotSlice = errors.New("values is not a slice")
    // ErrNoValues is returned when a cycler was initialized with no values.
    ErrNoValues = errors.New("no configured values")
)

// Cycler provides a way to cycle through a set of elements.
type Cycler interface {
    // Get returns the current value in the cycler. It does not rotate the cycler.
    Get() interface{}
    // Next rotates the cycler and returns the next value.
    Next() interface{}
    // Seek rotates the cycler until it finds a specific value.
    // If the values is not found the cycler will remain in the same index.
    Seek(target interface{}) interface{}
    // Len returns the number of values in the cycler.
    Len() int
}

// cycler is a Cycler implementation with internal state. The internal state allows it to keep track of its position
// in the set of elements, which comes in handy when used to share limited external resources between different threads.
// This implementation IS NOT thread safe. Synchronization between multiple threads accessing a single cycler must be
// handled separately.
type cycler struct {
    values []interface{}
    index int
}

// Get returns the current value in the cycler. It does not rotate the cycler.
func (c *cycler) Get() interface{} {
    return c.values[c.index]
}

// Next rotates the cycler and returns the value at the new position.
func (c *cycler) Next() interface{} {
    // Rotate the index
    c.index = (c.index + 1) % len(c.values)

    return c.Get()
}

// Seek rotates the cycler until it finds a specific value.
// If the values is not found the cycler will remain in the same index.
// This implementation is currently designed to support elementary values and structs containing elementary types and
// other structs. Pointers, and pointer based types (e.g. slices) will not be found unless they point to the exact
// same value in the cycler.
func (c *cycler) Seek(target interface{}) interface{} {
    var i int
    var value interface{}
    for i, value = 0, c.Get(); i < c.Len(); i, value = i+1, c.Next() {
        if value == target {
            return value
        }
    }

    return value
}

// Len returns the number of values in the cycler.
func (c *cycler) Len() int {
    return len(c.values)
}

// NewCyclerFromSlice creates a Cycler from a slice. The cycler will contain all the elements in the slice.
// This will return an error if `values` is not a slice or if the slice does not contain any elements.
func NewCyclerFromSlice(values interface{}) (Cycler, error) {
    // Verify that the input value is a slice
    v := reflect.ValueOf(values)
    if v.Kind() != reflect.Slice {
        return nil, ErrNotSlice
    }

    // Create an []interface{} slice from the input slice values
    slice := make([]interface{}, v.Len())
    for i:=0; i < v.Len(); i++ {
        slice[i] = v.Index(i).Interface()
    }

    return NewCycler(slice...)
}

// NewCycler creates a Cycler from passed values.
// This will return an error if no elements are passed.
func NewCycler(values ...interface{}) (Cycler, error) {
    if len(values) == 0 {
        return nil, ErrNoValues
    }

    return &cycler{
        values: values,
    }, nil
}