package containers

// Container represents a generic container unit used by different containerization platforms.
type Container interface {
	// Start starts the container.
	Start()
	// Stop stops the container.
	Stop()
	// Remove removes the container.
	Remove()
	// ID returns the container's id.
	ID() string
	// Name returns the container's name.
	Name() string
	// Image returns the container's image.
	Image() string
}
