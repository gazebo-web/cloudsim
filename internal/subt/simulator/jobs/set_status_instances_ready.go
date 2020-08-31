package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToInstancesReady is used to set a simulation status to instances ready.
var UpdateSimulationStatusToInstancesReady = &actions.Job{
	Name:            "set-simulation-status-instances-ready",
	Execute:         updateSimulationStatusToInstancesReady,
	RollbackHandler: rollbackUpdateSimulationStatusToWaiting,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToInstancesReady is the main process executed by UpdateSimulationStatusToInstancesReady.
func updateSimulationStatusToInstancesReady(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusInstancesReady)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToWaiting is in charge of setting the status to wait instances
// changed by updateSimulationStatusToWaitInstances.
func rollbackUpdateSimulationStatusToWaiting(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusWaitingInstances)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
