package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

type LaunchPodOutput struct {
	Resource orchestrator.Resource
	Error    error
}

// UpdateSimulationStatus is generic to job to update the status of a certain simulation.
var LaunchPod = &actions.Job{
	Name:            "launch-pod",
	Execute:         launchPod,
	RollbackHandler: rollbackUpdateSimulationStatus,
	InputType:       actions.GetJobDataType(orchestrator.CreatePodInput{}),
	OutputType:      actions.GetJobDataType(LaunchPodOutput{}),
}

func launchPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Create ctx
	simCtx := context.NewContext(ctx)

	// Parse input
	createPodInput, ok := value.(orchestrator.CreatePodInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Create pod
	res, err := simCtx.Platform().Orchestrator().Pods().Create(createPodInput)

	return LaunchPodOutput{
		Resource: res,
		Error:    err,
	}, nil
}
