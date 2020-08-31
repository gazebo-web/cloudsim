package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CheckParentSimulationWithError is used to check if a parent simulation of a certain children simulation has an error.
var CheckParentSimulationWithError = &actions.Job{
	Name:       "check-parent-simulation-with-error",
	Execute:    checkParentSimulationWithError,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkParentSimulationWithError is the main process executed by CheckParentSimulationWithError.
func checkParentSimulationWithError(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
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
