package transport

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

// WebsocketTransportMock represents a Transporter mock implementation.
type WebsocketTransportMock struct {
	WebsocketTransporter
	*mock.Mock
}

// Connect is a mock for the Connect method.
func (m *WebsocketTransportMock) Connect() error {
	args := m.Called()
	return args.Error(0)
}

// IsConnected is a mock for the IsConnected method.
func (m *WebsocketTransportMock) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

// Disconnect is a mock for the Disconnect method.
func (m *WebsocketTransportMock) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

// Connection is a mock for the Connection method.
func (m *WebsocketTransportMock) Connection() *websocket.Conn {
	args := m.Called()
	var c *websocket.Conn
	c = args.Get(0).(*websocket.Conn)
	return c
}

// NewWebsocketTransporterMock initializes a new WebsocketTransportMock object.
func NewWebsocketTransporterMock() *WebsocketTransportMock {
	return &WebsocketTransportMock{
		Mock: new(mock.Mock),
	}
}
