package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// UpdateSimulationStatusToNodesReady is used to set a simulation status to nodes ready.
var UpdateSimulationStatusToNodesReady = &actions.Job{
	Name:            "set-simulation-status-nodes-ready",
	Execute:         updateSimulationStatusToNodesReady,
	RollbackHandler: rollbackUpdateSimulationStatusToWaitingNodes,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToNodesReady is the main process executed by UpdateSimulationStatusToNodesReady.
func updateSimulationStatusToNodesReady(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusNodesReady)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToWaitingNodes is in charge of setting the status to waiting nodes
// changed by updateSimulationStatusToNodesReady.
func rollbackUpdateSimulationStatusToWaitingNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusWaitingNodes)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}
