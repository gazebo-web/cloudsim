package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// ChargeUser is a job that charges a user after a simulation has finished for the time it has been running.
var ChargeUser = &actions.Job{
	Name:       "charge-user",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    chargeUser,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// setRate is the execute function of the ChargeUser job that will charge a user for the time a simulation has been running.
func chargeUser(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	if !s.SubTServices().Billing().IsEnabled() {
		return s, nil
	}

	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	user, em := s.SubTServices().Users().GetUserFromUsername(sim.GetCreator())
	if em != nil {
		return nil, em.BaseError
	}

	err = s.SubTServices().Billing().SubtractCredits(context.Background(), user, sim)
	if err != nil {
		return nil, err
	}

	return s, nil
}
