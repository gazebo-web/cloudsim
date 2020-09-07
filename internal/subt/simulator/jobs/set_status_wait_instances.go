package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// UpdateSimulationStatusToWaitInstances is used to set a simulation status to waiting instances.
var UpdateSimulationStatusToWaitInstances = UpdateSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-wait-instances",
	PreHooks:   []actions.JobFunc{updateSimulationStatusToWaitInstances},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

func updateSimulationStatusToWaitInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	return UpdateSimulationStatusInput{
		GroupID:        gid,
		Status:         simulations.StatusWaitingInstances,
		PreviousStatus: simulations.StatusLaunchingInstances,
	}, nil
}
