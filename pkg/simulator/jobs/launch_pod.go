package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// LaunchPodsInput is the input of the LaunchPods job.
type LaunchPodsInput []orchestrator.CreatePodInput

// LaunchPodsOutput is the output of the LaunchPods job.
// This struct was set in place to let the post-hook handle errors.
type LaunchPodsOutput struct {
	Resources []orchestrator.Resource
	Error     error
}

// LaunchPods is generic to job to launch pods into a cluster.
var LaunchPods = &actions.Job{
	Execute: launchPods,
}

// launchPods is the main function executed by the LaunchPods job.
func launchPods(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input, ok := value.(LaunchPodsInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	var created []orchestrator.Resource
	var err error

	for _, in := range input {
		var res orchestrator.Resource
		res, err = s.Platform().Orchestrator().Pods().Create(in)
		if err != nil {
			return nil, err
		}
		created = append(created, res)
	}

	return LaunchPodsOutput{
		Resources: created,
		Error:     err,
	}, nil
}
