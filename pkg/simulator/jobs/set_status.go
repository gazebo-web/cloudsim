package jobs

import (
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/simulations"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// SetSimulationStatusInput is the input for SetSimulationStatus job.
type SetSimulationStatusInput struct {
	// GroupID identifies the simulation that should change the status.
	GroupID simulations.GroupID
	// Status is the status that will be assigned to a certain simulation.
	Status simulations.Status
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

	s := store.State().(state.ServicesGetter)

	err := s.Services().Simulations().UpdateStatus(input.GroupID, input.Status)
	if err != nil {
		return nil, err
	}

	return SetSimulationStatusOutput{
		GroupID: input.GroupID,
		Status:  input.Status,
	}, nil
}
