package ign

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
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
		msg := "test"
		pub := NewPublicationMessage("test", "ignition.msgs.StringMsg", msg)
		conn.WriteMessage(websocket.TextMessage, pub.ToByteSlice())
	}))

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		tr, err := NewIgnWebsocketTransporter(u.Host, u.Path, transport.WebsocketScheme, "")
		defer tr.Disconnect()
		assert.NoError(t, err)
		tr.Subscribe("test", func(message transport.Message) {
			var msg string
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
		msg := "test"
		pub := NewPublicationMessage("test", "ignition.msgs.StringMsg", msg)
		conn.WriteMessage(websocket.TextMessage, pub.ToByteSlice())
	}))

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		tr, err := NewIgnWebsocketTransporter(u.Host, u.Path, transport.WebsocketScheme, "")
		defer tr.Disconnect()
		assert.NoError(t, err)
		tr.Subscribe("test", func(message transport.Message) {
			var msg string
			err = message.GetPayload(&msg)
			assert.NoError(t, err)
		})
		conn.Close()
	})
}
