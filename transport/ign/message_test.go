package ign

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	msgs "gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
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

func TestNewMessageFromByteSlice(t *testing.T) {
	test := func(slice []byte, expected proto.Message, out proto.Message) {
		// Convert the message to an ign.Message struct
		message, err := NewMessageFromByteSlice(slice)
		require.NoError(t, err)
		// Unmarshal the payload
		err = message.unmarshal(out)
		require.NoError(t, err)
		// Check that the payload contains the expected value
		require.EqualValues(t, expected.String(), out.String())
	}

	// Topic: /subt/start
	// Message type: ignition.msgs.StringMsg
	slice := []byte{
		112, 117, 98, 44, 47, 115, 117, 98, 116, 47, 115, 116, 97, 114, 116, 44, 105, 103, 110, 105, 116, 105, 111,
		110, 46, 109, 115, 103, 115, 46, 83, 116, 114, 105, 110, 103, 77, 115, 103, 44, 10, 11, 10, 9, 8, 200, 15,
		16, 128, 242, 195, 215, 1, 18, 7, 115, 116, 97, 114, 116, 101, 100,
	}
	subtStartStarted := &msgs.StringMsg{
		Header: &msgs.Header{
			Stamp: &msgs.Time{
				Sec:  1992,
				Nsec: 452000000,
			},
		},
		Data: "started",
	}
	subtStartOut := &msgs.StringMsg{}
	test(slice, subtStartStarted, subtStartOut)

	// Topic: /subt/start
	// Message type: ignition.msgs.StringMsg
	slice = []byte{
		112, 117, 98, 44, 47, 115, 117, 98, 116, 47, 115, 116, 97, 114, 116, 44, 105, 103, 110, 105, 116, 105, 111,
		110, 46, 109, 115, 103, 115, 46, 83, 116, 114, 105, 110, 103, 77, 115, 103, 44, 10, 11, 10, 9, 8, 201, 29, 16,
		128, 246, 242, 182, 2, 18, 8, 102, 105, 110, 105, 115, 104, 101, 100,
	}
	subtStartFinished := &msgs.StringMsg{
		Header: &msgs.Header{
			Stamp: &msgs.Time{
				Sec:  3785,
				Nsec: 652000000,
			},
		},
		Data: "finished",
	}
	test(slice, subtStartFinished, subtStartOut)

	// Topic: /stats
	// Message type: ignition.msgs.WorldStatistics
	slice = []byte{
		112, 117, 98, 44, 47, 115, 116, 97, 116, 115, 44, 105, 103, 110, 105, 116, 105, 111, 110, 46, 109, 115, 103,
		115, 46, 87, 111, 114, 108, 100, 83, 116, 97, 116, 105, 115, 116, 105, 99, 115, 44, 18, 9, 8, 201, 29, 16,
		128, 246, 242, 182, 2, 34, 9, 8, 131, 30, 16, 132, 171, 197, 183, 1, 40, 1, 48, 237, 225, 57, 73, 97, 50, 85,
		48, 42, 169, 239, 63,
	}
	stats := &msgs.WorldStatistics{
		SimTime: &msgs.Time{
			Sec:  3785,
			Nsec: 652000000,
		},
		RealTime: &msgs.Time{
			Sec:  3843,
			Nsec: 384914820,
		},
		Paused:         true,
		Iterations:     uint64(946413),
		RealTimeFactor: float64(0.9894),
	}
	statsOut := &msgs.WorldStatistics{}
	test(slice, stats, statsOut)
}
