package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationNoErrorInput is the input of the CheckSimulationNoError job.
type CheckSimulationNoErrorInput []simulations.Simulation

// CheckSimulationNoErrorOutput is the output of the CheckSimulationNoError job.
type CheckSimulationNoErrorOutput bool

// CheckSimulationNoError is in charge of checking that the simulation has no error status.
var CheckSimulationNoError = &actions.Job{
	Execute: checkSimulationNoError,
}

// checkSimulationNoError is the execute function of the CheckSimulationNoError job.
func checkSimulationNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationNoErrorInput)

	for _, sim := range input {
		if sim.Error() != nil {
			return CheckSimulationNoErrorOutput(false), nil
		}
	}

	return CheckSimulationNoErrorOutput(true), nil
}
