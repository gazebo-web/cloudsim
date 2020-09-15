package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckKindInput is the input of the CheckKind job.
type CheckKindInput struct {
	Simulation simulations.Simulation
	Kind       simulations.Kind
}

// CheckKindOutput is the output of the CheckKind job.
type CheckKindOutput simulations.Simulation

// CheckKind is used to check that a certain simulation has a specific kind.
var CheckKind = &actions.Job{
	Name:    "check-kind",
	Execute: checkKind,
}

// checkKind is the execution of the CheckKind job.
func checkKind(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckKindInput)
	if input.Simulation.Kind() != input.Kind {
		return nil, simulations.ErrIncorrectKind
	}
	return CheckKindOutput(input.Simulation), nil
}
