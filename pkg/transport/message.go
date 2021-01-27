package transport

// Message represents a set of methods that message structs must implement to be used in a transport.
type Message interface {
	// Get returns the message payload.
	GetPayload(out interface{}) error
}
