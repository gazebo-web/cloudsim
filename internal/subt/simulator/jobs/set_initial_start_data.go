package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

// SetInitialStartData is used to set the basic data needed to start a simulation.
// It's usually the first job in the Start Simulation Action.
var SetInitialStartData = &actions.Job{
	Name:       "set-initial-data",
	Execute:    setInitialStartData,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

func setInitialStartData(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	data := ctx.Store().Get().(*StartSimulationData)

	data.GroupID = gid

	return gid, nil
}
