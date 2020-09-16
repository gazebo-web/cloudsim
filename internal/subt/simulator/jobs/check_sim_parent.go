package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CheckSimulationIsParent is a job in charge of checking if a simulation is parent.
var CheckSimulationIsParent = jobs.CheckSimulationKind.Extend(actions.Job{
	Name:            "check-parent",
	PreHooks:        []actions.JobFunc{createCheckKindInput},
	PostHooks:       []actions.JobFunc{assertIsParent, returnState},
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// assertIsParent evaluates if the value returned by CheckKind is true for the CheckSimulationIsParent job.
func assertIsParent(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isParent := value.(jobs.CheckSimulationKindOutput)
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

	return jobs.CheckSimulationKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}, nil
}
