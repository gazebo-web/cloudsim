package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToLaunchPods is used to set a simulation status to launch pods.
var UpdateSimulationStatusToLaunchPods = &actions.Job{
	Name:            "set-simulation-status-launch-pods",
	Execute:         updateSimulationStatusToLaunchPods,
	RollbackHandler: rollbackUpdateSimulationStatusToNodesReady,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToLaunchPods is the main process executed by UpdateSimulationStatusToLaunchPods.
func updateSimulationStatusToLaunchPods(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusLaunchingPods)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToNodesReady is in charge of setting the status nodes ready to the simulation
// changed by updateSimulationStatusToLaunchInstances.
func rollbackUpdateSimulationStatusToNodesReady(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusNodesReady)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
