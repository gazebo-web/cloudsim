package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// LaunchGazeboServerPod is used to launch the gazebo server pods for a simulation.
var LaunchGazeboServerPod = &actions.Job{
	Name:       "launch-gazebo-server-pod",
	Execute:    launchGazeboServerPod,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

func launchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	return gid, nil
}
