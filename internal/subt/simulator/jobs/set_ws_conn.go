package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
)

var SetWebsocketConnection = &actions.Job{
	Name:            "set-ws-conn",
	PreHooks:        []actions.JobFunc{setStartState},
	Execute:         connectWebsocket,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: revertWebsocketConnection,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
}

func connectWebsocket(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	host := s.Platform().Store().Ignition().GetWebsocketHost()
	path := s.Platform().Store().Ignition().GetWebsocketPath(s.GroupID)

	token, err := s.SubTServices().Simulations().GetWebsocketToken(s.GroupID)
	if err != nil {
		return nil, err
	}

	t, err := ignws.NewIgnWebsocketTransporter(host, path, transport.WebsocketSecureScheme, token)
	if err != nil {
		return nil, err
	}

	err = t.Connect()
	if err != nil {
		return nil, err
	}

	s.WebsocketConnection = t
	store.SetState(s)

	return store, nil
}

func revertWebsocketConnection(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if s.WebsocketConnection != nil && s.WebsocketConnection.IsConnected() {
		s.WebsocketConnection.Disconnect()
	}

	return nil, nil
}
