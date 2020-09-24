package transport

// Messager represents a set of methods that message structs must implement to be used in a transport.
type Messager interface {
	// Get returns the message payload.
	GetPayload(out interface{}) error
}
