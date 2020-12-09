package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationStatusInput is the input of the CheckSimulationStatus job.
type CheckSimulationStatusInput struct {
	// Simulation is the simulation that will be checked.
	Simulation simulations.Simulation
	// Status is the status that the Simulation should match.
	Status simulations.Status
}

// CheckSimulationStatusOutput is the output of the CheckSimulationStatus job.
type CheckSimulationStatusOutput bool

// CheckSimulationStatus is used to check that a certain simulation has a specific status.
// It returns true if the simulation matches the given status.
var CheckSimulationStatus = &actions.Job{
	Execute: checkSimulationStatus,
}

// checkSimulationStatus is the execute function of the CheckSimulationStatus job.
func checkSimulationStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationStatusInput)
	output := CheckSimulationStatusOutput(input.Simulation.HasStatus(input.Status))
	return output, nil
}
