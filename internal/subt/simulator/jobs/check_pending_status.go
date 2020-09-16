package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CheckPendingStatus is used to check that a certain simulation has pending status.
var CheckPendingStatus = jobs.CheckSimulationStatus.Extend(actions.Job{
	Name:            "check-pending-status",
	PreHooks:        []actions.JobFunc{createCheckPendingStatusInput},
	PostHooks:       []actions.JobFunc{assertPendingStatus, returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

// createCheckPendingStatusInput prepares the input for the check status job.
func createCheckPendingStatusInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	store.SetState(s)

	return jobs.CheckSimulationStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}, nil
}

// assertPendingStatus validates that the simulation has the pending status.
func assertPendingStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isPending := value.(jobs.CheckSimulationStatusOutput)
	if !isPending {
		return nil, simulations.ErrIncorrectStatus
	}
	return value, nil
}
