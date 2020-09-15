package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckStatusInput is the input of the CheckStatus job.
type CheckStatusInput struct {
	Simulation simulations.Simulation
	Status     simulations.Status
}

// CheckStatusOutput is the output of the CheckStatus job.
type CheckStatusOutput simulations.Simulation

// CheckStatus is used to check that a certain simulation has a specific status.
var CheckStatus = &actions.Job{
	Name:    "check-status",
	Execute: checkStatus,
}

// checkStatus is the execution of the CheckStatus job.
func checkStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input := value.(CheckStatusInput)
	if input.Simulation.Status() != input.Status {
		return nil, simulations.ErrIncorrectStatus
	}
	return CheckStatusOutput(input.Simulation), nil
}
