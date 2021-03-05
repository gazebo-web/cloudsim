package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveWebsocketService is a job extending the generic jobs.RemoveWebsocketService to remove the websocket service pointing to the
// Ignition Gazebo Server inside the
var RemoveWebsocketService = jobs.RemoveWebsocketService.Extend(actions.Job{
	Name:       "remove-websocket-service",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemoveWebsocketServiceInput},
	PostHooks:  []actions.JobFunc{checkRemoveWebsocketServiceError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// prepareRemoveWebsocketServiceInput is a pre-hook for the RemoveWebsocketService in charge of generating the input for
// the generic jobs.RemoveWebsocketService job.
func prepareRemoveWebsocketServiceInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	return jobs.RemoveWebsocketServiceInput{
		Name:      subtapp.GetServiceNameWebsocket(s.GroupID),
		Namespace: s.Platform().Store().Orchestrator().Namespace(),
	}, nil
}

// checkRemoveWebsocketServiceError is a post-hook for the RemoveWebsocketService in charge of checking that the generic
// jobs.RemoveWebsocketService job returns no error.
func checkRemoveWebsocketServiceError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.RemoveWebsocketServiceOutput)

	if out.Error != nil {
		return nil, out.Error
	}

	return nil, nil
}
