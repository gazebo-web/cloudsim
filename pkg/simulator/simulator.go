package simulator

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Simulator groups a set of methods to perform different operations with simulations.
// Simulator should be implemented by the applications.
type Simulator interface {
	// Start triggers an action to start simulations.
	Start(ctx context.Context, groupID simulations.GroupID) error

	// Stop triggers an action to stop a simulation.
	Stop(ctx context.Context, groupID simulations.GroupID) error

	// Restart triggers an action to restart a simulation.
	Restart(ctx context.Context, groupID simulations.GroupID) error
}
