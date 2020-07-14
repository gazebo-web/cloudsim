package ign

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	msgs "gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/transport"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSubscriberSuite(t *testing.T) {
	suite.Run(t, new(subscriberTestSuite))
}

type subscriberTestSuite struct {
	suite.Suite
	transport PubSubWebsocketTransporter
	upgrader  websocket.Upgrader
	handler   http.HandlerFunc
	server    *httptest.Server
	router    *http.ServeMux

	testTopic        string
	testTopicType    string
	testTopicMessage msgs.StringMsg

	message     []byte
	messageType int
	subscribed  bool
}

func (suite *subscriberTestSuite) SetupTest() {
	suite.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	suite.testTopic = "test"
	suite.testTopicType = "ignition.msgs.StringMsg"
	suite.testTopicMessage = msgs.StringMsg{
		Data: "test-server",
	}

	suite.handler = func(w http.ResponseWriter, r *http.Request) {
		conn, err := suite.upgrader.Upgrade(w, r, nil)
		suite.NoError(err)
		defer conn.Close()
		for {
			// TODO: Change how this works if more than one connection needs to be opened.
			if suite.subscribed {
				msg := fmt.Sprintf(
					"pub,%s,%s,%s",
					suite.testTopic,
					suite.testTopicType,
					suite.testTopicMessage.String(),
				)
				err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
				suite.NoError(err)
			}

			suite.messageType, suite.message, err = conn.ReadMessage()
			suite.NoError(err)
			m, err := NewMessageFromByteSlice(suite.message)
			suite.NoError(err)
			if m.Topic == suite.testTopic {
				suite.subscribed = true
				return
			}
			conn.WriteMessage(websocket.CloseUnsupportedData, []byte{})
		}
	}
	suite.router = http.NewServeMux()
	suite.router.Handle("/", suite.handler)
	suite.server = httptest.NewServer(suite.router)
}

func (suite *subscriberTestSuite) init() PubSubWebsocketTransporter {
	var err error
	u, err := url.Parse(suite.server.URL)
	suite.NoError(err)
	suite.transport, err = NewIgnWebsocketTransporter(u.Host, u.Path, transport.WebsocketScheme, "1234")
	suite.NoError(err)
	return suite.transport
}

func (suite *subscriberTestSuite) TestSubscribe_Accepted() {
	suite.init()
	err := suite.transport.Subscribe("test", func(message transport.Messager) {
		var payload string
		suite.NoError(message.GetPayload(&payload))
		suite.EqualValues(suite.testTopicMessage, payload)
	})
	suite.NoError(err)
}

func (suite *subscriberTestSuite) TestSubscribe_Rejected() {
	suite.init()
	_ = suite.transport.Subscribe("wrong-test", func(message transport.Messager) {
		suite.Equal(nil, message)
	})
}

func (suite *subscriberTestSuite) AfterTest() {
	if suite.transport != nil {
		suite.transport.Disconnect()
	}
	suite.server.Close()
}
