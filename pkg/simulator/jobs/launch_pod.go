package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// LaunchPodInput is the input of the LaunchPod job.
type LaunchPodInput orchestrator.CreatePodInput

// LaunchPodOutput is the output of the LaunchPod job.
type LaunchPodOutput struct {
	Resource orchestrator.Resource
	Error    error
}

// LaunchPod is generic to job to launch pods into a cluster.
var LaunchPod = &actions.Job{
	Execute: launchPod,
}

func launchPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Create ctx
	simCtx := context.NewContext(ctx)

	// Parse input
	input, ok := value.(LaunchPodInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	createPodInput := orchestrator.CreatePodInput(input)

	// Create pod
	res, err := simCtx.Platform().Orchestrator().Pods().Create(createPodInput)
	if err != nil {
		return nil, err
	}

	return LaunchPodOutput{
		Resource: res,
		Error:    err,
	}, nil
}
