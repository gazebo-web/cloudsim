package ign

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/transport"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestTransporterListenDontPanicConnClosed(t *testing.T) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Used to guarantee that the waiter goroutine is always executed before the websocket listener
	var waiterLock sync.Mutex
	// Used to guarantee that the websocket server closes the connection before the websocket listener is terminated
	var connLock sync.Mutex

	// Start the websocket server
	// This server opens and immediately closes websocket connections
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)
		wsConn.Close()
	}))

	waiterLock.Lock()

	// Waiter goroutine
	// This block buys time to allow the server connection to close and cause the listener websocket to fail
	go func() {
		connLock.Lock()
		waiterLock.Unlock()
		// Give the server some time to close the websocket connection
		time.Sleep(100 * time.Millisecond)
		connLock.Unlock()
	}()

	// Ensure that the previous block runs before the next one
	waiterLock.Lock()

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		// Allow the test to terminate
		defer waiterLock.Unlock()

		tr, err := NewIgnWebsocketTransporter(context.TODO(), u.Host, u.Path, transport.WebsocketScheme, "")
		defer tr.Disconnect()
		assert.NoError(t, err)

		// Wait until the server is given time to terminate the connection
		connLock.Lock()
	})

	// Wait until the websocket listener has finished processing
	waiterLock.Lock()

	// Close the server
	server.Close()
}
