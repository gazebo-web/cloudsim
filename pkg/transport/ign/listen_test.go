package ign

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	msgs "gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTransporterListenDontPanicHTTPClosed(t *testing.T) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var conn *websocket.Conn
	var err error

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err = upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		assert.NoError(t, err)
		msg := msgs.StringMsg{
			Data: "test",
		}
		pub := NewPublicationMessage("test", "ignition.msgs.StringMsg", msg.String())
		conn.WriteMessage(websocket.TextMessage, pub.ToByteSlice())
	}))

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		tr, err := NewIgnWebsocketTransporter(u.Host, u.Path, transport.WebsocketScheme, "")
		defer tr.Disconnect()
		assert.NoError(t, err)
		tr.Subscribe("test", func(message transport.Message) {
			var msg msgs.StringMsg
			err = message.GetPayload(&msg)
			assert.NoError(t, err)
		})
		server.Close()
	})
}

func TestTransporterListenDontPanicWSClosed(t *testing.T) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var conn *websocket.Conn
	var err error

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err = upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		assert.NoError(t, err)
		msg := msgs.StringMsg{
			Data: "test",
		}
		pub := NewPublicationMessage("test", "ignition.msgs.StringMsg", msg.String())
		conn.WriteMessage(websocket.TextMessage, pub.ToByteSlice())
	}))

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		tr, err := NewIgnWebsocketTransporter(u.Host, u.Path, transport.WebsocketScheme, "")
		defer tr.Disconnect()
		assert.NoError(t, err)
		tr.Subscribe("test", func(message transport.Message) {
			var msg msgs.StringMsg
			err = message.GetPayload(&msg)
			assert.NoError(t, err)
		})
		conn.Close()
	})
}
