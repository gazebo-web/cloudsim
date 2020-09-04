package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// UpdateSimulationStatusToWaitNodes is used to set a simulation status to wait nodes.
var UpdateSimulationStatusToWaitNodes = UpdateSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-wait-nodes",
	PreHooks:   []actions.JobFunc{updateSimulationStatusToWaitNodes},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

func updateSimulationStatusToWaitNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}
	return UpdateSimulationStatusInput{
		GroupID:        gid,
		Status:         simulations.StatusWaitingNodes,
		PreviousStatus: simulations.StatusWaitingInstances,
	}, nil
}
