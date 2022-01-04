package ign

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessageFromByteSlice_ConvertedMessage(t *testing.T) {
	m := &Message{
		Operation: "sub",
		Topic:     "test",
		Type:      "1234",
		Payload:   "hello-world",
	}
	text := m.ToString()
	b := []byte(text)
	result, err := NewMessageFromByteSlice(b)
	assert.NoError(t, err)
	assert.Equal(t, m, result)
}

func TestNewMessageFromByteSlice_RandomInput_WorldStatistics(t *testing.T) {
	input := "pub,/world/tunnel_circuit_practice_01/stats,ignition.msgs.WorldStatistics,\u0012   �\u0004\u0010��þ\u0002\"�\u0006\u0010����\u00010�I,Ԛ�\u001D��?"
	result, err := NewMessageFromByteSlice([]byte(input))
	assert.NoError(t, err)
	assert.Equal(t, result.Operation, "pub")
	assert.Equal(t, result.Topic, "/world/tunnel_circuit_practice_01/stats")
	assert.Equal(t, result.Type, "ignition.msgs.WorldStatistics")
	assert.Equal(t, result.Payload, "\u0012   �\u0004\u0010��þ\u0002\"�\u0006\u0010����\u00010�I,Ԛ�\u001D��?")
}

func TestNewMessageFromByteSlice_AlmostEmpty(t *testing.T) {
	input := "sub,,,test"
	_, err := NewMessageFromByteSlice([]byte(input))
	assert.NoError(t, err)
}

func TestNewMessageFromByteSlice_SubTStart(t *testing.T) {
	input := "pub,/subt/start,ignition.msgs.StringMsg,`\u0010���&\u0012\u0004init"
	result, err := NewMessageFromByteSlice([]byte(input))
	assert.NoError(t, err)
	assert.Equal(t, result.Operation, "pub")
	assert.Equal(t, result.Topic, "/subt/start")
	assert.Equal(t, result.Type, "ignition.msgs.StringMsg")
	assert.Equal(t, result.Payload, "`\u0010���&\u0012\u0004init")
}

func TestNewAuthorizationMessage(t *testing.T) {
	result := NewAuthorizationMessage("test")
	assert.Equal(t, result.Operation, "auth")
	assert.Equal(t, result.Payload, "test")
}

func TestNewSubscriptionMessage(t *testing.T) {
	result := NewSubscriptionMessage("test")
	assert.Equal(t, result.Operation, "sub")
	assert.Equal(t, result.Topic, "test")
}

func TestMessage_ToString(t *testing.T) {
	m := &Message{
		Operation: "sub",
		Topic:     "test",
		Type:      "1234",
		Payload:   "payload-test",
	}
	expected := "sub,test,1234,payload-test"
	assert.Equal(t, expected, m.ToString())
}

func TestMessage_ToByteSlice(t *testing.T) {
	m := &Message{
		Operation: "sub",
		Topic:     "test",
		Type:      "1234",
		Payload:   "payload-test",
	}
	b := []byte(m.ToString())
	result := m.ToByteSlice()
	assert.Equal(t, b, result)
}

func TestMessage_ToSlice(t *testing.T) {
	m := &Message{
		Operation: "sub",
		Topic:     "test",
		Type:      "1234",
		Payload:   "payload-test",
	}
	s := []string{"sub", "test", "1234", "payload-test"}

	assert.Equal(t, s, m.ToSlice())
}
