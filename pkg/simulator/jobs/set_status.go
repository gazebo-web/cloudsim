package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SetSimulationStatusInput is the input for SetSimulationStatus job.
type SetSimulationStatusInput struct {
	// GroupID identifies the simulation that should change the status.
	GroupID simulations.GroupID
	// Status is the status that will be assigned to a certain simulation.
	Status simulations.Status
	// SetStatus is a function used to set a certain status to the given simulation.
	SetStatus func(GroupID simulations.GroupID, status simulations.Status) error
}

// SetSimulationStatusOutput is the output of the SetSimulationStatus job.
type SetSimulationStatusOutput struct {
	// GroupID identifies the simulation that has changed its status.
	GroupID simulations.GroupID
	// Status is the status that has been applied to the simulation identified by the GroupID.
	Status simulations.Status
}

// SetSimulationStatus is used to set a certain status to a simulation.
var SetSimulationStatus = &actions.Job{
	Execute: setSimulationStatus,
}

// setSimulationStatus is the execute function of the SetSimulationStatus job.
func setSimulationStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(SetSimulationStatusInput)

	err := input.SetStatus(input.GroupID, input.Status)
	if err != nil {
		return nil, err
	}

	return SetSimulationStatusOutput{
		GroupID: input.GroupID,
		Status:  input.Status,
	}, nil
}
