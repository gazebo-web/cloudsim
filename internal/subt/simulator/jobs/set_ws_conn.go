package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"time"
)

// SetWebsocketConnection is a job in charge of setting a websocket connection to the Ignition Gazebo server.
var SetWebsocketConnection = &actions.Job{
	Name:            "set-ws-conn",
	PreHooks:        []actions.JobFunc{setStartState},
	Execute:         connectWebsocket,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: revertWebsocketConnection,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
}

// connectWebsocket is the execute function of the SetWebsocketConnection job.
func connectWebsocket(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	host := s.Platform().Store().Orchestrator().IngressHost()

	path := s.Platform().Store().Ignition().GetWebsocketPath(s.GroupID)

	token, err := s.SubTServices().Simulations().GetWebsocketToken(s.GroupID)
	if err != nil {
		return nil, err
	}

	var t ignws.PubSubWebsocketTransporter
	for i := 0; i < 10; i++ {
		t, err = ignws.NewIgnWebsocketTransporter(host, path, transport.WebsocketSecureScheme, token)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i*10) * time.Second)
	}

	if err != nil {
		return nil, err
	}

	s.WebsocketConnection = t
	store.SetState(s)

	return s, nil
}

// revertWebsocketConnection is in charge of disconnecting the websocket server if an error happens while setting
// the connection up.
func revertWebsocketConnection(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if s.WebsocketConnection != nil && s.WebsocketConnection.IsConnected() {
		s.WebsocketConnection.Disconnect()
	}

	return nil, nil
}
