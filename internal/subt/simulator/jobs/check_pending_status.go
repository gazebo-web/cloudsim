package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CheckPendingStatus is used to check that a certain simulation has pending status.
var CheckPendingStatus = &actions.Job{
	Name:       "check-pending-status",
	Execute:    checkPendingStatus,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkPendingStatus is the main process executed by CheckPendingStatus.
func checkPendingStatus(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	data, ok := value.(*StartSimulationData)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(data.GroupID)
	if err != nil {
		return nil, err
	}

	if sim.Status() != simulations.StatusPending {
		return nil, simulations.ErrIncorrectStatus
	}
	return data, nil
}
