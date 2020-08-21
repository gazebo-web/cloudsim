package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// JobCheckPendingStatus is used to check that a certain simulation has pending status.
var JobCheckPendingStatus = &actions.Job{
	Name:     "check-pending-status",
	PreHooks: nil,
	// TODO: Use checkPendingStatus.
	Execute:         nil,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// JobCheckSimulationParenthood is used to check that a simulation is not of kind simulations.SimParent.
var JobCheckSimulationParenthood = &actions.Job{
	Name:     "check-simulation-parenthood",
	PreHooks: nil,
	// TODO: Use checkSimulationIsParent.
	Execute:         nil,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// JobCheckParentSimulationWithError is used to check if a parent simulation of a certain children simulation has an error.
var JobCheckParentSimulationWithError = &actions.Job{
	Name:     "check-parent-simulation-with-error",
	PreHooks: nil,
	// TODO: Use checkParentSimulationWithError.
	Execute:         nil,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// checkPendingStatus is the main process executed by JobCheckSimulationParenthood.
func checkPendingStatus(ctx actions.Context, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := ctx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Status() != simulations.StatusPending {
		return nil, simulations.ErrIncorrectStatus
	}
	return gid, nil
}

// checkSimulationIsParent is the main process executed by JobCheckPendingStatus.
func checkSimulationIsParent(ctx actions.Context, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := ctx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Kind() == simulations.SimParent {
		_, err := ctx.Services().Simulations().Reject(gid)
		if err != nil {
			return nil, err
		}
		return nil, simulations.ErrIncorrectKind
	}
	return gid, nil
}

// checkParentSimulationWithError is the main process executed by JobCheckParentSimulationWithError.
func checkParentSimulationWithError(ctx actions.Context, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := ctx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Kind() != simulations.SimChild {
		return gid, nil
	}
	parent, err := ctx.Services().Simulations().GetParent(gid)
	if err != nil {
		return nil, err
	}
	if simerr := parent.Error(); simerr != nil {
		return nil, simulations.ErrParentSimulationWithError
	}
	return gid, nil
}
