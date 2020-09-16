package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// LaunchPodInput is the input of the LaunchPod job.
type LaunchPodInput orchestrator.CreatePodInput

// LaunchPodOutput is the output of the LaunchPod job.
// This struct was set in place to let the post-hook handle errors.
type LaunchPodOutput struct {
	Resource orchestrator.Resource
	Error    error
}

// LaunchPod is generic to job to launch pods into a cluster.
var LaunchPod = &actions.Job{
	Execute: launchPod,
}

// launchPod is the main function executed by the LaunchPod job.
func launchPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.Platform)

	// Parse input
	input, ok := value.(LaunchPodInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	createPodInput := orchestrator.CreatePodInput(input)

	// Create pod
	res, err := s.Platform().Orchestrator().Pods().Create(createPodInput)
	if err != nil {
		return nil, err
	}

	return LaunchPodOutput{
		Resource: res,
		Error:    err,
	}, nil
}
