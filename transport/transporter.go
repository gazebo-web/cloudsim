package transport

// Transporter represents a set of methods that will open and close communication between two processes.
type Transporter interface {
	Connect() error
	IsConnected() bool
	Disconnect()
}
