package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// LaunchWebsocketServiceInput is the input of the LaunchWebsocketService job.
type LaunchWebsocketServiceInput orchestrator.CreateServiceInput

// LaunchWebsocketServiceOutput is the output of the LaunchWebsocketService job.
// This struct was set in place to let the post-hook handle errors.
type LaunchWebsocketServiceOutput struct {
	Resource orchestrator.Resource
	Error    error
}

// LaunchWebsocketService is generic to job to launch a simulation's websocket service.
var LaunchWebsocketService = &actions.Job{
	Execute: launchWebsocketService,
}

// launchWebsocketService is the main function executed by the LaunchWebsocketService job.
func launchWebsocketService(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.Platform)

	// Parse input
	input, ok := value.(orchestrator.CreateServiceInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Create service
	res, err := s.Platform().Orchestrator().Services().Create(input)

	return LaunchWebsocketServiceOutput{
		Resource: res,
		Error:    err,
	}, nil
}
