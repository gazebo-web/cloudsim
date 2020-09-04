package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// UpdateSimulationStatusToInstancesReady is used to set a simulation status to instances ready.
var UpdateSimulationStatusToInstancesReady = UpdateSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-instances-ready",
	PreHooks:   []actions.JobFunc{updateSimulationStatusToInstancesReady},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

func updateSimulationStatusToInstancesReady(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	return UpdateSimulationStatusInput{
		GroupID:        gid,
		Status:         simulations.StatusInstancesReady,
		PreviousStatus: simulations.StatusWaitingInstances,
	}, nil
}
