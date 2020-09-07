package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

const dataPreviousStatusKey = "previous-status"

// UpdateSimulationStatusInput is the input used by the UpdateSimulationStatus execute function.
// The specific applications should use this input to pass data to the main process from prehooks.
type UpdateSimulationStatusInput struct {
	GroupID        simulations.GroupID
	Status         simulations.Status
	PreviousStatus simulations.Status
}

// UpdateSimulationStatus is generic to job to update the status of a certain simulation.
var UpdateSimulationStatus = &actions.Job{
	Name:            "set-simulation-status",
	Execute:         updateSimulationStatus,
	RollbackHandler: rollbackUpdateSimulationStatus,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatus is the main process executed by UpdateSimulationStatus.
func updateSimulationStatus(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	input, ok := value.(UpdateSimulationStatusInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(input.GroupID, input.Status)

	if dataErr := deployment.SetJobData(tx, nil, dataPreviousStatusKey, input.PreviousStatus); dataErr != nil {
		return nil, dataErr
	}

	if err != nil {
		return nil, err
	}

	return input.GroupID, nil
}

func rollbackUpdateSimulationStatus(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	input, dataErr := deployment.GetJobData(tx, nil, dataPreviousStatusKey)
	if dataErr != nil {
		return nil, dataErr
	}

	status, ok := input.(simulations.Status)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)

	updateErr := simCtx.Services().Simulations().UpdateStatus(gid, status)
	if updateErr != nil {
		return nil, updateErr
	}

	return nil, nil
}
