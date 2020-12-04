package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// GenerateCheckSimulationKindJob generates a specific check simulation kind job.
// The generated jobs are used to check if a simulation is single, parent or child.
func GenerateCheckSimulationKindJob(name string, kind simulations.Kind, inputType, outputType interface{}) *actions.Job {
	createCheckKindInput := generateCheckSimulationKindInputPreHook(kind)

	return jobs.CheckSimulationKind.Extend(actions.Job{
		Name:       name,
		PreHooks:   []actions.JobFunc{createCheckKindInput},
		PostHooks:  []actions.JobFunc{assertSimulationKind, returnState},
		InputType:  actions.GetJobDataType(inputType),
		OutputType: actions.GetJobDataType(outputType),
	})
}

// generateCheckSimulationKindInputPreHook generates a pre-hook to get the simulation from a certain group ID
// passed in the action store and prepares the proper dto for the generic job to check simulation kind.
func generateCheckSimulationKindInputPreHook(kind simulations.Kind) actions.JobFunc {
	return func(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		s := value.(*state.StartSimulation)

		store.SetState(s)

		sim, err := s.Services().Simulations().Get(s.GroupID)
		if err != nil {
			return nil, err
		}

		return jobs.CheckSimulationKindInput{
			Simulation: sim,
			Kind:       kind,
		}, nil
	}
}

// assertSimulationKind is the post-hook in charge of guarantee that the output of the CheckSimulationKind job operation
// is of the correct kind.
func assertSimulationKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isKind := value.(jobs.CheckSimulationKindOutput)
	if isKind {
		return nil, simulations.ErrIncorrectKind
	}
	return nil, nil
}
