package simulator

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

var JobsStartSimulation = actions.Jobs{
	JobCheckPendingStatus,
	JobCheckSimulationParenthood,
	JobCheckParentSimulationWithError,
}

var JobsStopSimulation = actions.Jobs{}

var JobsRestartSimulation = actions.Jobs{}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckPendingStatus is used to check that a certain simulation has pending status.
var JobCheckPendingStatus = &actions.Job{
	Name:            "check-pending-status",
	PreHooks:        nil,
	Execute:         checkPendingStatus,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// checkPendingStatus is the main process executed by JobCheckSimulationParenthood.
func checkPendingStatus(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Status() != simulations.StatusPending {
		return nil, simulations.ErrIncorrectStatus
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckSimulationParenthood is used to check that a simulation is not of kind simulations.SimParent.
var JobCheckSimulationParenthood = &actions.Job{
	Name:            "check-simulation-parenthood",
	PreHooks:        nil,
	Execute:         checkSimulationIsParent,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// checkSimulationIsParent is the main process executed by JobCheckPendingStatus.
func checkSimulationIsParent(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Kind() == simulations.SimParent {
		_, err := simCtx.Services().Simulations().Reject(gid)
		if err != nil {
			return nil, err
		}
		return nil, simulations.ErrIncorrectKind
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckParentSimulationWithError is used to check if a parent simulation of a certain children simulation has an error.
var JobCheckParentSimulationWithError = &actions.Job{
	Name:            "check-parent-simulation-with-error",
	PreHooks:        nil,
	Execute:         checkParentSimulationWithError,
	PostHooks:       nil,
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// checkParentSimulationWithError is the main process executed by JobCheckParentSimulationWithError.
func checkParentSimulationWithError(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}
	if sim.Kind() != simulations.SimChild {
		return gid, nil
	}
	parent, err := simCtx.Services().Simulations().GetParent(gid)
	if err != nil {
		return nil, err
	}
	if simerr := parent.Error(); simerr != nil {
		return nil, simulations.ErrParentSimulationWithError
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------
