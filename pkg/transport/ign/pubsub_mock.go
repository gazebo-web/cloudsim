package ign

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

// PubSubTransporterMock implements PubSubWebsocketTransporter to be used as mock for testing purposes.
type PubSubTransporterMock struct {
	*mock.Mock
}

// NewPubSubTransporterMock initializes a new PubSubTransporterMock
func NewPubSubTransporterMock() *PubSubTransporterMock {
	return &PubSubTransporterMock{
		Mock: new(mock.Mock),
	}
}

// Subscribe is a mock for the Subscribe method.
func (m *PubSubTransporterMock) Subscribe(topic string, cb Callback) error {
	args := m.Called(topic, cb)
	return args.Error(0)
}

// Unsubscribe is a mock for the Unsubscribe method.
func (m *PubSubTransporterMock) Unsubscribe(topic string) error {
	args := m.Called(topic)
	return args.Error(0)
}

// Publish is a mock for the Publish method.
func (m *PubSubTransporterMock) Publish(message Message) error {
	args := m.Called(message)
	return args.Error(0)
}

// Connect is a mock for the Connect method.
func (m *PubSubTransporterMock) Connect() error {
	args := m.Called()
	return args.Error(0)
}

// IsConnected is a mock for the IsConnected method.
func (m *PubSubTransporterMock) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

// Disconnect is a mock for the Disconnect method.
func (m *PubSubTransporterMock) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

// Connection is a mock for the Connection method.
func (m *PubSubTransporterMock) Connection() *websocket.Conn {
	return m.Called().Get(0).(*websocket.Conn)
}
