package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// ConfigCheckSimulationNotOfKindJob is used to configure the job returned by the GenerateCheckSimulationNotOfKindJob function.
type ConfigCheckSimulationNotOfKindJob struct {
	// Name is the name of the job.
	Name string
	// Kind is the Kind that should be compared.
	Kind simulations.Kind
	// PreHooks are the pre-hooks that should go before the PreparationPreHook.
	PreHooks []actions.JobFunc
	// PreparationPreHook is a single pre-hook in charge of setting the input for the CheckSimulationKind generic job.
	PreparationPreHook actions.JobFunc
	// InputType is the job's input type.
	InputType interface{}
	// OutputType is the job's output type.
	OutputType interface{}
}

// GenerateCheckSimulationNotOfKindJob is a job generator for checking that a certain simulation is not of a certain kind.
func GenerateCheckSimulationNotOfKindJob(config ConfigCheckSimulationNotOfKindJob) *actions.Job {
	return jobs.CheckSimulationKind.Extend(actions.Job{
		Name:       config.Name,
		PreHooks:   append(config.PreHooks, config.PreparationPreHook),
		PostHooks:  []actions.JobFunc{assertSimulationNotOfKind, returnState},
		InputType:  actions.GetJobDataType(config.InputType),
		OutputType: actions.GetJobDataType(config.OutputType),
	})
}

// generateCheckStartSimulationNotOfKindInputPreHook generates a pre-hook to get the simulation from a certain group ID
// passed in the action store and prepares the proper dto for the generic job to check simulation is not of a certain kind.
func generateCheckStartSimulationNotOfKindInputPreHook(kind simulations.Kind) actions.JobFunc {
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

// generateCheckStopSimulationNotOfKindInputPreHook generates a pre-hook to get the simulation from a certain group ID
// passed in the action store and prepares the proper dto for the generic job to check simulation is not of a certain kind.
func generateCheckStopSimulationNotOfKindInputPreHook(kind simulations.Kind) actions.JobFunc {
	return func(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		s := value.(*state.StopSimulation)

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

// assertSimulationNotOfKind is the post-hook in charge of guaranteeing that the output of the jobs.CheckSimulationKindOutput job operation
// is not of a certain kind.
func assertSimulationNotOfKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	isKind := value.(jobs.CheckSimulationKindOutput)
	if isKind {
		return nil, simulations.ErrIncorrectKind
	}
	return nil, nil
}
