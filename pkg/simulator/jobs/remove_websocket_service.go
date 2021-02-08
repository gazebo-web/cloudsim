package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemoveWebsocketServiceInput is the input of the RemoveWebsocketService job.
type RemoveWebsocketServiceInput struct {
	Name      string
	Namespace string
}

// RemoveWebsocketServiceOutput is the output of the RemoveWebsocketService job.
// This struct was set in place to let the post-hook handle errors.
type RemoveWebsocketServiceOutput struct {
	Error error
}

// RemoveWebsocketService is generic to job to remove simulation's websocket services.
var RemoveWebsocketService = &actions.Job{
	Execute: removeWebsocketService,
}

// removeWebsocketService is the main function executed by the RemoveWebsocketService job.
func removeWebsocketService(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input, ok := value.(RemoveWebsocketServiceInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Get the service
	res, err := s.Platform().Orchestrator().Services().Get(input.Name, input.Namespace)
	if err != nil {
		return RemoveWebsocketServiceOutput{
			Error: err,
		}, nil
	}

	// Remove the service
	err = s.Platform().Orchestrator().Services().Remove(res)

	return RemoveWebsocketServiceOutput{
		Error: err,
	}, nil
}
