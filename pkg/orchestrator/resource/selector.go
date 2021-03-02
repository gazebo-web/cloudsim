package resource

import "fmt"

// Selector is used to identify a certain resource.
type Selector interface {
	// String returns the selector represented in string format.
	String() string
	// Map returns the underlying selector's map.
	Map() map[string]string
	// Extend extends the underlying base map with the extension selector.
	// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
	Extend(extension Selector) Selector
	// Set sets the given value to the given key. If the key already exists, it will be overwritten.
	Set(key string, value string)
}

// selector is a group of key-pair values that identify a resource.
type selector map[string]string

// Set sets the given value to the given key. If the key already exists, it will be overwritten.
func (s selector) Set(key string, value string) {
	s[key] = value
}

// Extend extends the underlying base map with the extension selector.
// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
func (s selector) Extend(extension Selector) Selector {
	for k, v := range extension.Map() {
		s[k] = v
	}
	return s
}

// Map returns the selector in map format.
func (s selector) Map() map[string]string {
	return s
}

// String returns the selector in string format.
func (s selector) String() string {
	var out string
	count := len(s)
	for key, value := range s {
		out += fmt.Sprintf("%s=%s", key, value)
		count--
		if count > 1 {
			out += ","
		}
	}
	return out
}

// NewSelector initializes a new Selector from the given map.
// If `nil` is passed as input, an empty selector will be returned.
func NewSelector(input map[string]string) Selector {
	if input == nil {
		input = map[string]string{}
	}
	return selector(input)
}
