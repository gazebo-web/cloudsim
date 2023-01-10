package jobs

import (
	"context"
	"github.com/gazebo-web/cloudsim/v4/pkg/actions"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulator"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulator/state"
	"github.com/jinzhu/gorm"
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

		// If assertion fails but LaunchPodsInput is nil, consider as no pods need to be launched.
		if input == nil {
			return LaunchPodsOutput{
				Resources: []resource.Resource{},
				Error:     nil,
			}, nil
		}

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
		res, err = s.Platform().Orchestrator().Pods().Create(context.Background(), in)
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
