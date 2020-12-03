package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationKindInput is the input of the CheckSimulationKind job.
type CheckSimulationKindInput struct {
	// The simulation that will be checked.
	Simulation simulations.Simulation
	// The kind that the simulation should match.
	Kind simulations.Kind
}

// CheckSimulationKindOutput is the output of the CheckSimulationKind job.
// It will be true if the simulation is of the expected kind.
type CheckSimulationKindOutput bool

// CheckSimulationKind is used to check that a certain simulation is of a specific kind.
// CheckSimulationKind can be used to check if a simulation is single, parent or child.
var CheckSimulationKind = &actions.Job{
	Execute: checkSimulationKind,
}

// checkSimulationKind is the execution of the CheckSimulationKind job.
func checkSimulationKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationKindInput)
	return CheckSimulationKindOutput(input.Simulation.IsKind(input.Kind)), nil
}
