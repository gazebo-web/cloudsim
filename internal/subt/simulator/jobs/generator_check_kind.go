package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// GenerateCheckSimulationNotOfKindJob is a job generator for checking that a certain simulation is not of a certain kind.
func GenerateCheckSimulationNotOfKindJob(name string, kind simulations.Kind, inputType, outputType interface{}) *actions.Job {
	createCheckKindInput := generateCheckSimulationNotOfKindInputPreHook(kind)

	return jobs.CheckSimulationKind.Extend(actions.Job{
		Name:       name,
		PreHooks:   []actions.JobFunc{setStartState, createCheckKindInput},
		PostHooks:  []actions.JobFunc{assertSimulationNotOfKind, returnState},
		InputType:  actions.GetJobDataType(inputType),
		OutputType: actions.GetJobDataType(outputType),
	})
}

// generateCheckSimulationNotOfKindInputPreHook generates a pre-hook to get the simulation from a certain group ID
// passed in the action store and prepares the proper dto for the generic job to check simulation is not of a certain kind.
func generateCheckSimulationNotOfKindInputPreHook(kind simulations.Kind) actions.JobFunc {
	return func(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		s := value.(*state.StartSimulation)

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

// assertSimulationNotOfKind is the post-hook in charge of guarantee that the output of the jobs.CheckSimulationKindOutput job operation
// is not of a certain kind.
func assertSimulationNotOfKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isKind := value.(jobs.CheckSimulationKindOutput)
	if isKind {
		return nil, simulations.ErrIncorrectKind
	}
	return nil, nil
}
