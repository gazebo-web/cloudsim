package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CheckSimulationNoErrors checks that a group of simulations don't have errors.
var CheckSimulationNoErrors = jobs.CheckSimulationNoError.Extend(actions.Job{
	Name:       "check-sim-no-errors",
	PreHooks:   []actions.JobFunc{createCheckSimulationNoErrorInput},
	PostHooks:  []actions.JobFunc{checkNoErrorOutput, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// checkNoErrorOutput checks that the simulations provided to the execute function have no errors.
func checkNoErrorOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.CheckSimulationNoErrorOutput)
	if out.Error != nil {
		return nil, fmt.Errorf("error while checking if simulations have error status, base error: %w", out.Error)
	}
	return nil, nil
}

// createCheckSimulationNoErrorInput creates the input for the execute function.
func createCheckSimulationNoErrorInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get input
	s := value.(*state.StartSimulation)

	// Set store state
	store.SetState(s)

	// Prepare input
	var input jobs.CheckSimulationNoErrorInput

	// Get simulation from the GroupID
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Add simulation to the next job input.
	input = append(input, sim)

	// If the simulation isn't a parent, execute the job.
	if !sim.IsKind(simulations.SimParent) {
		return input, nil
	}

	// If simulation is parent, add the parent to check for errors as well.
	parent, err := s.Services().Simulations().GetParent(s.GroupID)
	if err != nil {
		return nil, err
	}
	input = append(input, parent)

	return input, nil
}
