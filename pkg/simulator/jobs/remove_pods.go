package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemovePodsInput is the input of the RemovePods job.
type RemovePodsInput []orchestrator.Resource

// RemovePodsOutput is the output of the RemovePods job.
// This struct was set in place to let the post-hook handle errors.
type RemovePodsOutput struct {
	Resources []orchestrator.Resource
	Error     error
}

// LaunchPods is a generic job to remove pods from a cluster.
var RemovePods = &actions.Job{
	Execute: removePods,
}

// launchPods is the main function executed by the LaunchPods job.
func removePods(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input, ok := value.(RemovePodsInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	if len(input) == 0 {
		return LaunchPodsOutput{
			Resources: []orchestrator.Resource{},
			Error:     nil,
		}, nil
	}

	var deleted []orchestrator.Resource
	var err error

	for _, in := range input {
		var res orchestrator.Resource
		res, err = s.Platform().Orchestrator().Pods().Delete(in)
		if err != nil {
			return nil, err
		}
		deleted = append(deleted, res)
	}

	return RemovePodsOutput{
		Resources: deleted,
		Error:     err,
	}, nil
}
