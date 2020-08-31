package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToWaitInstances is used to set a simulation status to waiting instances.
var UpdateSimulationStatusToWaitInstances = &actions.Job{
	Name:            "set-simulation-status-wait-instances",
	Execute:         updateSimulationStatusToWaitInstances,
	RollbackHandler: rollbackUpdateSimulationStatusToLaunching,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToWaitInstances is the main process executed by UpdateSimulationStatusToWaitInstances.
func updateSimulationStatusToWaitInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusWaitingInstances)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToLaunching is in charge of setting the status to launching instances
// changed by updateSimulationStatusToWaitInstances.
func rollbackUpdateSimulationStatusToLaunching(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusLaunchingInstances)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
