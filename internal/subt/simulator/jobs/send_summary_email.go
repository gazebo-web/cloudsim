package jobs

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// SendSummaryEmail is a job in charge of sending an email to participants with the simulation's statistics and score.
var SendSummaryEmail = &actions.Job{
	Name:            "send-summary-email",
	PreHooks:        []actions.JobFunc{setStopState},
	Execute:         sendSummaryEmail,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(&state.StopSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StopSimulation{}),
}

// sendSummaryEmail is the execute function of the SendSummaryEmail job.
func sendSummaryEmail(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	if sim.IsProcessed() {
		return nil, errors.New("simulation has been processed")
	}

	return s, nil
}
