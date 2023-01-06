package jobs

import (
	"fmt"
	"github.com/gazebo-web/cloudsim/v4/pkg/actions"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulations"
	"github.com/jinzhu/gorm"
)

// CheckSimulationNoErrorInput is the input of the CheckSimulationNoError job.
type CheckSimulationNoErrorInput []simulations.Simulation

// CheckSimulationNoErrorOutput is the output of the CheckSimulationNoError job.
type CheckSimulationNoErrorOutput struct {
	Error error
}

// CheckSimulationNoError is in charge of checking that the simulation has no error status.
var CheckSimulationNoError = &actions.Job{
	Execute: checkSimulationNoError,
}

// checkSimulationNoError is the execute function of the CheckSimulationNoError job.
func checkSimulationNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationNoErrorInput)

	for _, sim := range input {
		if sim.GetError() != nil {
			return CheckSimulationNoErrorOutput{
				Error: fmt.Errorf("simulation [%s] with error status [%s]", sim.GetGroupID(), *sim.GetError()),
			}, nil
		}
	}

	return CheckSimulationNoErrorOutput{
		Error: nil,
	}, nil
}
