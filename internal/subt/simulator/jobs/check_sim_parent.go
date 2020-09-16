package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CheckSimulationParent is a job in charge of checking if a simulation is parent.
var CheckSimulationParent = jobs.CheckKind.Extend(actions.Job{
	Name:            "check-parent",
	PreHooks:        []actions.JobFunc{createCheckKindInput},
	PostHooks:       []actions.JobFunc{assertIsParent, returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

// assertIsParent evaluates if the value returned by CheckKind is true for the CheckSimulationParent job.
func assertIsParent(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isParent := value.(bool)
	if isParent {
		return nil, simulations.ErrIncorrectKind
	}
	return nil, nil
}

// createCheckKindInput is in charge of creating the input for the jobs.CheckKind generic job.
func createCheckKindInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	return jobs.CheckKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}, nil
}
