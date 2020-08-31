package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToLaunchNodes is used to set a simulation status to launch nodes.
var UpdateSimulationStatusToLaunchNodes = &actions.Job{
	Name:            "set-simulation-status-launch-nodes",
	Execute:         updateSimulationStatusToLaunchNodes,
	RollbackHandler: rollbackUpdateSimulationStatusToPending,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToLaunchNodes is the main process executed by UpdateSimulationStatusToLaunchNodes.
func updateSimulationStatusToLaunchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusLaunchingNodes)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToPending is in charge of setting the status pending to the simulation
// changed by updateSimulationStatusToLaunchNodes.
func rollbackUpdateSimulationStatusToPending(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusPending)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
