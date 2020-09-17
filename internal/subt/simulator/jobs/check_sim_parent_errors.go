package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

var CheckSimulationNoErrors = jobs.CheckSimulationNoError.Extend(actions.Job{
	Name:            "check-parent-sim-no-errors",
	PreHooks:        []actions.JobFunc{createCheckSimulationNoErrorInput},
	PostHooks:       []actions.JobFunc{checkNoErrorOutput, returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

func checkNoErrorOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	hasErr := value.(jobs.CheckSimulationNoErrorOutput)
	if hasErr {
		return nil, simulations.ErrSimulationWithError
	}
	return nil, nil
}

func createCheckSimulationNoErrorInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)
	store.SetState(s)

	var input jobs.CheckSimulationNoErrorInput

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	input = append(input, sim)

	if sim.Kind() != simulations.SimChild {
		return input, nil
	}

	parent, err := s.Services().Simulations().GetParent(s.GroupID)
	if err != nil {
		return nil, err
	}

	input = append(input, parent)

	return input, nil
}
