package jobs

import (
	"context"
	"github.com/gazebo-web/cloudsim/v4/pkg/actions"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/services"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// LaunchWebsocketServiceInput is the input of the LaunchWebsocketService job.
type LaunchWebsocketServiceInput services.CreateServiceInput

// LaunchWebsocketServiceOutput is the output of the LaunchWebsocketService job.
// This struct was set in place to let the post-hook handle errors.
type LaunchWebsocketServiceOutput struct {
	Resource resource.Resource
	Error    error
}

// LaunchWebsocketService is generic to job to launch a simulation's websocket service.
var LaunchWebsocketService = &actions.Job{
	Execute: launchWebsocketService,
}

// launchWebsocketService is the main function executed by the LaunchWebsocketService job.
func launchWebsocketService(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input := value.(LaunchWebsocketServiceInput)

	// Create service
	res, err := s.Platform().Orchestrator().Services().Create(context.TODO(), services.CreateServiceInput(input))

	return LaunchWebsocketServiceOutput{
		Resource: res,
		Error:    err,
	}, nil
}
