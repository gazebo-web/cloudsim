package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// UpdateSimulationStatusToLaunchPods is used to set a simulation status to launch pods.
var UpdateSimulationStatusToLaunchPods = UpdateSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-launch-pods",
	PreHooks:   []actions.JobFunc{updateSimulationStatusToLaunchPods},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

func updateSimulationStatusToLaunchPods(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	return UpdateSimulationStatusInput{
		GroupID:        gid,
		Status:         simulations.StatusLaunchingPods,
		PreviousStatus: simulations.StatusWaitingNodes,
	}, nil
}
