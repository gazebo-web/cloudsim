package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CheckSimulationParenthood is used to check that a simulation is not of kind simulations.SimParent.
var CheckSimulationParenthood = &actions.Job{
	Name:       "check-simulation-parenthood",
	Execute:    checkSimulationIsParent,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkSimulationIsParent is the main process executed by CheckPendingStatus.
func checkSimulationIsParent(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

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
