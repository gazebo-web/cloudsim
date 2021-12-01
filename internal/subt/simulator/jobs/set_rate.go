package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// SetRate is a job in charge of setting Rate field for a certain simulation in order to charge the user when the simulation
// is marked as completed.
var SetRate = &actions.Job{
	Name:       "set-rate",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    setRate,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
}

// setRate is the execute function of the SetRate job that will set the rate a which a simulation should be charged.
func setRate(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if !s.SubTServices().Billing().IsEnabled() {
		return s, nil
	}

	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	if !s.SubTServices().Users().IsSystemAdmin(sim.GetCreator()) {
		return s, nil
	}

	rate, err := s.Platform().Machines().CalculateCost(s.CreateMachinesInput)
	if err != nil {
		return nil, err
	}

	sim.SetRate(rate)

	err = s.SubTServices().Simulations().Update(s.GroupID, sim)
	if err != nil {
		return nil, err
	}

	return s, nil
}
