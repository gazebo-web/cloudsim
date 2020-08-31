package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToWaitNodes is used to set a simulation status to wait nodes.
var UpdateSimulationStatusToWaitNodes = &actions.Job{
	Name:            "set-simulation-status-wait-instances",
	Execute:         updateSimulationStatusToWaitNodes,
	RollbackHandler: rollbackUpdateSimulationStatusToInstancesReady,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToWaitNodes is the main process executed by UpdateSimulationStatusToWaitNodes.
func updateSimulationStatusToWaitNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusWaitingNodes)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToInstancesReady is in charge of setting the status to instances ready
// changed by updateSimulationStatusToWaitNodes.
func rollbackUpdateSimulationStatusToInstancesReady(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusInstancesReady)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
