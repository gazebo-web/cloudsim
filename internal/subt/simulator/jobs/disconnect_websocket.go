package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// DisconnectWebsocket is a job in charge of disconnecting the websocket client.
var DisconnectWebsocket = &actions.Job{
	Name:            "disconnect-websocket",
	PreHooks:        []actions.JobFunc{setStopState},
	Execute:         disconnectWebsocket,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: revertWebsocketConnection,
	InputType:       actions.GetJobDataType(&state.StopSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StopSimulation{}),
}

// disconnectWebsocket is in charge of disconnecting the websocket client for a certain running simulation.
func disconnectWebsocket(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	s.Platform().RunningSimulations().Free(s.GroupID)

	return s, nil
}
