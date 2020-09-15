package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CheckPendingStatus is used to check that a certain simulation has pending status.
var CheckPendingStatus = jobs.CheckStatus.Extend(actions.Job{
	Name:            "check-pending-status",
	PreHooks:        []actions.JobFunc{createCheckStatusInput},
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

func createCheckStatusInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	store.SetState(s)

	return jobs.CheckStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}, nil
}
