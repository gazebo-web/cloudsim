package simulator

import (
	"context"
	"errors"
	"github.com/gazebo-web/cloudsim/pkg/platform"
	"github.com/gazebo-web/cloudsim/pkg/simulations"
)

var (
	// ErrInvalidInput is returned when an invalid input is provided.
	ErrInvalidInput = errors.New("invalid input")
)

// Simulator groups a set of methods to perform different operations with simulations.
// Simulator should be implemented by the applications.
type Simulator interface {
	// Start triggers an action to start simulations.
	Start(ctx context.Context, platform platform.Platform, groupID simulations.GroupID) error

	// Stop triggers an action to stop a simulation.
	Stop(ctx context.Context, platform platform.Platform, groupID simulations.GroupID) error
}
