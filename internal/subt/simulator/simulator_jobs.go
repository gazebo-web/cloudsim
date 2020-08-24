package simulator

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	JobCheckPendingStatus,
	JobCheckSimulationParenthood,
	JobCheckParentSimulationWithError,
	JobUpdateSimulationStatusToLaunchNodes,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{}

// JobsRestartSimulation groups the jobs needed to restart a simulation.
var JobsRestartSimulation = actions.Jobs{}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckPendingStatus is used to check that a certain simulation has pending status.
var JobCheckPendingStatus = &actions.Job{
	Name:       "check-pending-status",
	Execute:    checkPendingStatus,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
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
	Name:       "check-simulation-parenthood",
	Execute:    checkSimulationIsParent,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
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
	Name:       "check-parent-simulation-with-error",
	Execute:    checkParentSimulationWithError,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
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

// JobUpdateSimulationStatusToLaunchNodes is used to set a simulation status to launch nodes.
var JobUpdateSimulationStatusToLaunchNodes = &actions.Job{
	Name:       "set-simulation-status-launch-nodes",
	Execute:    updateSimulationStatusToLaunchNodes,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToLaunchNodes is the main process executed by JobUpdateSimulationStatusToLaunchNodes.
func updateSimulationStatusToLaunchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusLaunchingNodes)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobLaunchNodes is used to launch the required nodes to run a simulation.
var JobLaunchNodes = &actions.Job{
	Name:       "launch-nodes",
	Execute:    launchNodes,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(cloud.CreateMachinesOutput{}),
}

// launchNodes is the main process executed by JobLaunchNodes.
func launchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	output, err := simCtx.Platform().Machines().Create(sim.ToCreateMachinesInput())
	deployment.SetJobData(tx, &deployment.CurrentJob, actions.GetJobDataType(cloud.CreateMachinesOutput{}), output)

	if err != nil {
		return nil, err
	}

	err = simCtx.Services().Simulations().Update(gid, sim)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------
