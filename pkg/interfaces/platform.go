package interfaces

import "context"

// IPlatform defines the set of methods of a Platform.
type IPlatform interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RequestLaunch(ctx context.Context, groupID string)
	RequestTermination(ctx context.Context, groupID string)
}