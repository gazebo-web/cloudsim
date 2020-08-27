package transport

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

// websocketTransport is a WebsocketTransporter implementation.
type websocketTransport struct {
	Address    url.URL
	connection *websocket.Conn
}

// WebsocketConnector is a group of methods that handle websocket connections.
type WebsocketConnector interface {
	Connection() *websocket.Conn
}

// WebsocketTransporter extends the Transporter default behavior for websockets.
type WebsocketTransporter interface {
	Transporter
	WebsocketConnector
}

// Connect establishes a connection to the websocket server.
func (w *websocketTransport) Connect() error {
	if w.connection != nil {
		return errors.New("connection already established")
	}
	conn, err := createConnection(w.Address)
	if err != nil {
		return err
	}
	w.connection = conn
	return nil
}

// createConnection creates a new websocket connection based on the given URL and topic.
// It returns an error if the connection wasn't established.
func createConnection(addr url.URL) (*websocket.Conn, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(addr.String(), nil)
	// Temporary debug code
	if err == websocket.ErrBadHandshake {
		if resp == nil {
			fmt.Println("Websocket debug: resp is nil")
		} else {
			fmt.Println("Websocket debug:", resp.Status, resp.Body, resp)
		}
	}
	return conn, err
}

// IsConnected checks if the connection has been established
func (w *websocketTransport) IsConnected() bool {
	if w.connection == nil {
		return false
	}
	err := w.connection.WriteMessage(websocket.PingMessage, []byte{})
	return err == nil
}

// Disconnect closes the connection.
func (w *websocketTransport) Disconnect() {
	w.connection.Close()
	w.connection = nil
}

// Connection returns the active websocket connection.
func (w *websocketTransport) Connection() *websocket.Conn {
	return w.connection
}

// NewWebsocketTransporter initializes a new WebsocketTransporter instance using a websocket implementation.
// It will also establish a connection to the given addr.
// It will return an error if the connection to the given address failed.
func NewWebsocketTransporter(host, path, scheme string) (WebsocketTransporter, error) {
	wst := &websocketTransport{
		Address: url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   path,
		},
	}

	if err := wst.Connect(); err != nil {
		return nil, err
	}

	return wst, nil
}
