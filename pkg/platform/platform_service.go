package platform

import "context"

// IPlatformService represents a set of methods to perform on the platform.
type IPlatformService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Start starts the platform.
func (p Platform) Start(ctx context.Context) error {
	return nil
}

// Stop stops the platform.
func (p Platform) Stop(ctx context.Context) error {
	return nil
}
