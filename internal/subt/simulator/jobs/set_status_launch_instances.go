package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// UpdateSimulationStatusToLaunchInstances is used to set a simulation status to launch instances.
var UpdateSimulationStatusToLaunchInstances = UpdateSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-launch-instances",
	PreHooks:   []actions.JobFunc{updateSimulationStatusToLaunchInstances},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

func updateSimulationStatusToLaunchInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	return UpdateSimulationStatusInput{
		GroupID:        gid,
		Status:         simulations.StatusLaunchingInstances,
		PreviousStatus: simulations.StatusPending,
	}, nil
}
