package containers

import "fmt"

// Container represents a generic container unit used by different containerization platforms.
type Container interface {
	// Start starts the container.
	Start() error
	// Stop stops the container.
	Stop() error
	// Remove removes the container.
	Remove() error
	// ID returns the container's id.
	ID() string
	// Name returns the container's name.
	Name() string
	// Image returns the container's image.
	Image() string
	// EnvVars returns the container's environment variables.
	EnvVars() EnvVars
}

// EnvVars represents a group of environment variables.
type EnvVars map[string]string

// ToMap returns the underlying map.
func (env EnvVars) ToMap() map[string]string {
	if len(env) == 0 {
		return nil
	}
	return env
}

// ToSlice converts the underlying map into a slice.
func (env EnvVars) ToSlice() []string {
	var result []string
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}
