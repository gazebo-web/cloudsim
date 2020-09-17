package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// CheckSimulationNoErrorInput is the input of the CheckSimulationNoError job.
type CheckSimulationNoErrorInput struct {
	GroupID simulations.GroupID
}

// CheckSimulationNoErrorOutput is the output of the CheckSimulationNoError job.
type CheckSimulationNoErrorOutput bool

// CheckSimulationNoError is in charge of checking that the simulation has no errors.
var CheckSimulationNoError = &actions.Job{
	Execute: checkSimulationNoError,
}

// checkSimulationNoError is the execute function of the CheckSimulationNoError job.
func checkSimulationNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationNoErrorInput)

	s := store.State().(state.Services)

	sim, err := s.Services().Simulations().Get(input.GroupID)
	if err != nil {
		return nil, err
	}

	return CheckSimulationNoErrorOutput(sim.Error() == nil), nil
}
