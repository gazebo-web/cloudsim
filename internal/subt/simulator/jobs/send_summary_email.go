package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SendSummaryEmail is a job in charge of sending an email to participants with the simulation's statistics and score.
var SendSummaryEmail = &actions.Job{
	Name:       "send-summary-email",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    sendSummaryEmail,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// sendSummaryEmail is the execute function of the SendSummaryEmail job.
func sendSummaryEmail(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	if sim.IsProcessed() {
		return nil, simulations.ErrSimulationProcessed
	}

	user, em := s.Services().Users().GetUserFromUsername(sim.GetCreator())
	if em != nil {
		return nil, em.BaseError
	}

	var recipients []string

	recipients = append(recipients, s.Platform().Store().Ignition().DefaultRecipients()...)

	if user.Email != nil {
		recipients = append(recipients, *user.Email)
	}

	owner := sim.GetOwner()
	if owner == nil {
		// TODO: Send summary to recipients
		return s, nil
	}

	org, em := s.Services().Users().GetOrganization(*owner)
	if em != nil {
		return nil, em.BaseError
	}

	if org.Email != nil {
		recipients = append(recipients, *org.Email)
	}

	// TODO: Send summary to recipients

	return s, nil
}
