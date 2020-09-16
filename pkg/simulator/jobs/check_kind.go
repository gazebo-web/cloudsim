package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationKindInput is the input of the CheckSimulationKind job.
type CheckSimulationKindInput struct {
	Simulation simulations.Simulation
	Kind       simulations.Kind
}

// CheckSimulationKindOutput is the output of the CheckSimulationKind job.
type CheckSimulationKindOutput bool

// CheckSimulationKind is used to check that a certain simulation has a specific kind.
var CheckSimulationKind = &actions.Job{
	Execute: checkSimulationKind,
}

// checkSimulationKind is the execution of the CheckSimulationKind job.
func checkSimulationKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckSimulationKindInput)
	return input.Simulation.Kind() == input.Kind, nil
}
