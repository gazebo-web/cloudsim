package transport

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
)

// websocketTransport is a WebsocketTransporter implementation.
type websocketTransport struct {
	Address    url.URL
	connection *websocket.Conn
	connLock   sync.Mutex
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
	w.connLock.Lock()
	defer w.connLock.Unlock()

	if w.connection != nil && w.IsConnected() {
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
	conn, _, err := websocket.DefaultDialer.Dial(addr.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// IsConnected checks if the connection has been established
func (w *websocketTransport) IsConnected() bool {
	w.connLock.Lock()
	defer w.connLock.Unlock()

	if w == nil || w.connection == nil {
		return false
	}
	err := w.connection.WriteMessage(websocket.PingMessage, []byte{})
	return err == nil
}

// Disconnect closes the connection.
func (w *websocketTransport) Disconnect() error {
	w.connLock.Lock()
	defer w.connLock.Unlock()

	if w == nil || w.connection == nil {
		return nil
	}
	err := w.connection.Close()
	if err != nil {
		return err
	}
	w.connection = nil
	return nil
}

// Connection returns the active websocket connection.
func (w *websocketTransport) Connection() *websocket.Conn {
	w.connLock.Lock()
	defer w.connLock.Unlock()
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
