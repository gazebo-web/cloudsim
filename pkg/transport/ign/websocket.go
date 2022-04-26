package ign

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	"sync"
)

// Callback is a function that will be executed after reading a message from a certain topic.
type Callback func(ctx context.Context, message transport.Message)

// Publisher represents a set of methods that will let some process send messages to another process.
type Publisher interface {
	Publish(message Message) error
}

// Subscriber represents a set of methods that will let some process receive messages from another process.
type Subscriber interface {
	Subscribe(topic string, cb Callback) error
	Unsubscribe(topic string) error
}

// PubSubWebsocketTransporter represents a set of methods to communicate two processes using the Publisher and
// Subscriber interfaces.
type PubSubWebsocketTransporter interface {
	transport.WebsocketTransporter
	Subscriber
	Publisher
}

// websocketPubSubTransport is a websocket transport implementation with Ignition Robotics Publish/Subscribe protocol
// support.
type websocketPubSubTransport struct {
	transport.WebsocketTransporter
	topics     map[string]Callback
	topicsLock sync.RWMutex
	listening  bool
	baseCtx    context.Context
}

func newWebsocketPubSubTransport(baseCtx context.Context, transport transport.WebsocketTransporter) (*websocketPubSubTransport, error) {
	pubsub := &websocketPubSubTransport{
		baseCtx:              baseCtx,
		WebsocketTransporter: transport,
		topics:               make(map[string]Callback, 0),
	}

	if err := pubsub.listen(); err != nil {
		return nil, err
	}

	return pubsub, nil
}

func (w *websocketPubSubTransport) listen() error {
	// Check that this transport is not listening already
	if w.listening {
		return errors.New("already listening to websocket connection")
	}
	w.listening = true

	// Start the listener
	go func() {
		// Recover from panics to prevent a websocket connection from terminating the server
		defer func() {
			if p := recover(); p != nil {
				logger := ign.NewLogger("ws_cb_proxy", true, ign.VerbosityDebug)
				logger.Critical("Panic while running websocket transport listen() function: ", p)
			}
		}()

		for {
			if w.Connection() == nil {
				return
			}
			messageType, message, err := w.Connection().ReadMessage()
			if err == nil && message != nil {
				w.processMessage(messageType, message)
			}
		}
	}()

	return nil
}

func (w *websocketPubSubTransport) processMessage(messageType int, message []byte) {
	w.topicsLock.RLock()
	defer w.topicsLock.RUnlock()

	// Try to parse the incoming message as a Message struct
	if message, err := NewMessageFromByteSlice(message); err == nil {
		if cb, ok := w.topics[message.Topic]; ok {
			cb(w.baseCtx, message)
		}
	}
}

// Subscribe establishes a subscription to the given topic.
// It will also create a new go routine that will read messages until the connection is lost or closed.
// The incoming messages will be sent as an argument for the given callback.
func (w *websocketPubSubTransport) Subscribe(topic string, cb Callback) error {
	w.topicsLock.Lock()
	defer w.topicsLock.Unlock()

	sub := NewSubscriptionMessage(topic)

	// Send a subscription message to the websocket server
	err := w.Connection().WriteMessage(websocket.TextMessage, sub.ToByteSlice())
	if err != nil {
		return err
	}

	// Register the callback
	if _, ok := w.topics[topic]; ok {
		return errors.New("already subscribed to topic")
	}
	w.topics[topic] = cb

	return nil
}

// Unsubscribe closes a connection to the given topic.
func (w *websocketPubSubTransport) Unsubscribe(topic string) error {
	w.topicsLock.Lock()
	defer w.topicsLock.Unlock()

	delete(w.topics, topic)

	return nil
}

// Publish sends a message to the given topic.
func (w *websocketPubSubTransport) Publish(message Message) error {
	return nil
}

// NewIgnWebsocketTransporter initializes a new PubSubWebsocketTransporter instance using a websocketPubSubTransport
// implementation. It also establishes a connection to the given addr and sends an authorization message with the
// given token. The token should be the same as the simulation authorization token from the simulation that the
// transporter is attempting to connect to.
func NewIgnWebsocketTransporter(baseCtx context.Context, host, path, scheme, token string) (PubSubWebsocketTransporter, error) {
	wst, err := transport.NewWebsocketTransporter(host, path, scheme)
	if err != nil {
		return nil, err
	}

	pubsub, err := newWebsocketPubSubTransport(baseCtx, wst)
	if err != nil {
		return nil, err
	}

	// Send an authorization message using the given token.
	auth := NewAuthorizationMessage(token)
	if err := pubsub.Connection().WriteMessage(websocket.TextMessage, auth.ToByteSlice()); err != nil {
		return nil, err
	}

	return pubsub, nil
}
