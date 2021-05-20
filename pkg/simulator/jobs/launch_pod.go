package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// LaunchPodsInput is the input of the LaunchPods job.
type LaunchPodsInput []pods.CreatePodInput

// LaunchPodsOutput is the output of the LaunchPods job.
// This struct was set in place to let the post-hook handle errors.
type LaunchPodsOutput struct {
	Resources []resource.Resource
	Error     error
}

// LaunchPods is a generic job to launch pods on a cluster.
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

	if len(input) == 0 {
		return LaunchPodsOutput{
			Resources: []resource.Resource{},
			Error:     nil,
		}, nil
	}

	var created []resource.Resource
	var err error

	for _, in := range input {
		var res resource.Resource
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
