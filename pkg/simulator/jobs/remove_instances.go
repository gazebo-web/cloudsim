package jobs

import (
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/machines"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// RemoveInstancesInput is the input of the RemoveInstances job.
// It's used to pass the list of instances to remove.
type RemoveInstancesInput []machines.TerminateMachinesInput

// RemoveInstances is a generic job to remove instances.
var RemoveInstances = &actions.Job{
	Execute: removeInstances,
}

func removeInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get the store
	s := store.State().(state.PlatformGetter)

	// Parse the input
	input := value.(RemoveInstancesInput)

	// Trigger the machine termination.
	for _, in := range input {
		_ = s.Platform().Machines().Terminate(in)
	}

	return nil, nil
}
