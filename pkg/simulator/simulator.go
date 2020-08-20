package simulator

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Simulator groups a set of methods that returns actions to perform
// different operations with simulations.
// Simulator should be implemented by the applications.
type Simulator interface {
	// Start returns the action to start a simulation.
	Start(ctx context.Context, groupID simulations.GroupID) actions.Action

	// Stop returns the action to stop a simulation.
	Stop(ctx context.Context, groupID simulations.GroupID) actions.Action

	// Restart returns the action to restart a simulation.
	Restart(ctx context.Context, groupID simulations.GroupID) actions.Action
}
