package ign

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
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
	testTopicMessage string

	subscribeLock sync.Mutex
	messageLock   sync.Mutex

	message     []byte
	messageType int
	subscribed  bool
}

func (suite *subscriberTestSuite) SetupTest() {
	suite.messageLock.Lock()
	defer suite.messageLock.Unlock()

	suite.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	suite.testTopic = "test"
	suite.testTopicType = "ignition.msgs.StringMsg"
	suite.testTopicMessage = "test-server"

	suite.handler = suite.testSubscriberHandler(suite.upgrader, suite.testTopic)
	suite.router = http.NewServeMux()
	suite.router.Handle("/", suite.handler)
	suite.server = httptest.NewServer(suite.router)
}

func (suite *subscriberTestSuite) testSubscriberHandler(upgrader websocket.Upgrader, testTopic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		// TODO: Change how this works if more than one connection needs to be opened.
		for {
			suite.subscribeLock.Lock()
			if suite.subscribed {
				msg := fmt.Sprintf(
					"pub,%s,%s,%s",
					suite.testTopic,
					suite.testTopicType,
					suite.testTopicMessage,
				)
				err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			suite.subscribeLock.Unlock()

			suite.messageLock.Lock()
			suite.messageType, suite.message, err = conn.ReadMessage()
			suite.messageLock.Unlock()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			suite.messageLock.Lock()
			m, err := NewMessageFromByteSlice(suite.message)
			suite.messageLock.Unlock()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if m.Topic == testTopic {
				suite.subscribeLock.Lock()
				suite.subscribed = true
				suite.subscribeLock.Unlock()
				return
			}
			conn.WriteMessage(websocket.CloseUnsupportedData, []byte{})
		}
	}
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
	err := suite.transport.Subscribe("test", func(message transport.Message) {
		var payload string
		suite.NoError(message.GetPayload(&payload))
		suite.EqualValues(suite.testTopicMessage, payload)
	})
	suite.NoError(err)
}

func (suite *subscriberTestSuite) TestSubscribe_Rejected() {
	suite.init()
	_ = suite.transport.Subscribe("wrong-test", func(message transport.Message) {
		suite.Equal(nil, message)
	})
}

func (suite *subscriberTestSuite) AfterTest() {
	if suite.transport != nil {
		suite.transport.Disconnect()
	}
	suite.server.Close()
}
