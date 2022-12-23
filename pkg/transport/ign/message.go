package ign

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"strings"
)

// Message represents a message format for Ignition Robotics websocket server communication.
type Message struct {
	// Operation is the operation request that will be performed in the websocket server.
	//
	// auth: Authorize connection to the websocket server.
	//
	// sub: Subscribe to a certain topic.
	//
	// topics: List the available topics.
	//
	// pub: Publish to a certain topic - Not implemented yet.
	//
	// protos: Request message definitions.
	Operation string
	// Topic is the name of the topic to subscribe or publish to.
	Topic string
	// Type is the name of the ign-msgs protobuf type used to marshal/unmarshal payloads.
	Type string
	// Payload is the actual message. This message needs to be unmarshalled with protobuf before using.
	Payload string
}

// GetPayload returns the message payload.
func (m *Message) GetPayload(out interface{}) error {
	protoMessage, ok := out.(proto.Message)
	if !ok {
		return errors.New("Message.Get received an out value of invalid type")
	}
	return m.unmarshal(protoMessage)
}

// unmarshal assigns the current payload to the given protobuf struct pointer.
// It returns an error if the operation failed.
func (m *Message) unmarshal(p proto.Message) error {
	return proto.Unmarshal([]byte(m.Payload), p)
}

// NewMessageFromByteSlice initializes a new message from the given slice.
// It returns an error if the slice is invalid
func NewMessageFromByteSlice(slice []byte) (*Message, error) {
	var m Message
	values := strings.SplitN(string(slice), ",", 4)
	if len(values) != 4 {
		return nil, errors.New("invalid slice length")
	}
	m.Operation = values[0]
	m.Topic = values[1]
	m.Type = values[2]
	m.Payload = values[3]
	return &m, nil
}

// ToSlice converts the message into a slice.
func (m Message) ToSlice() []string {
	return []string{m.Operation, m.Topic, m.Type, m.Payload}
}

// ToString converts the message into a string.
func (m Message) ToString() string {
	return strings.Join(m.ToSlice(), ",")
}

// ToByteSlice converts the message into a slice of bytes.
func (m Message) ToByteSlice() []byte {
	return []byte(m.ToString())
}

// NewAuthorizationMessage creates a new authorization message from the given token.
func NewAuthorizationMessage(token string) Message {
	return Message{
		Operation: "auth",
		Payload:   token,
	}
}

// NewSubscriptionMessage creates a new subscription message from the given topic.
func NewSubscriptionMessage(topic string) Message {
	return Message{
		Operation: "sub",
		Topic:     topic,
	}
}

// NewPublicationMessage creates a new publication message with a certain type to send to the given topic.
// messageType is the name of a protobuf type used to marshal the message.
func NewPublicationMessage(topic string, messageType, message string) Message {
	return Message{
		Operation: "pub",
		Topic:     topic,
		Type:      messageType,
		Payload:   message,
	}
}
