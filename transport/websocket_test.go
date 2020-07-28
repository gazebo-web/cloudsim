package transport

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

func TestWebsocketSuite(t *testing.T) {
	suite.Run(t, new(websocketTestSuite))
}

type websocketTestSuite struct {
	suite.Suite
	transport WebsocketTransporter
	upgrader  websocket.Upgrader
	handler   http.HandlerFunc
	server    *httptest.Server
	router    *http.ServeMux

	message     []byte
	lock sync.Mutex
}

func (suite *websocketTestSuite) SetupTest() {
	suite.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	suite.handler = func(w http.ResponseWriter, r *http.Request) {
		conn, err := suite.upgrader.Upgrade(w, r, nil)
		suite.NoError(err)
		defer conn.Close()
		if err = conn.WriteMessage(websocket.TextMessage, []byte("test-server")); err != nil {
			return
		}
		suite.lock.Lock()
		defer suite.lock.Unlock()
		var msg []byte
		if _, msg, err = conn.ReadMessage(); err != nil {
			suite.message = msg
			return
		}
	}
	suite.router = http.NewServeMux()
	suite.router.Handle("/", suite.handler)
	suite.server = httptest.NewServer(suite.router)
}

func (suite *websocketTestSuite) init() Transporter {
	var err error
	u, err := url.Parse(suite.server.URL)
	suite.NoError(err)
	suite.transport, err = NewWebsocketTransporter(u.Host, u.Path, WebsocketScheme)
	suite.NoError(err)
	return suite.transport
}

func (suite *websocketTestSuite) AfterTest() {
	if suite.transport != nil {
		suite.transport.Disconnect()
	}
	suite.server.Close()
}

func (suite *websocketTestSuite) TestConnection_Accepted() {
	suite.init().Disconnect()
	suite.NoError(suite.transport.Connect())
}

func (suite *websocketTestSuite) TestConnection_Rejected() {
	suite.init()
	var err error
	suite.transport.Disconnect()
	suite.transport, err = NewWebsocketTransporter("wrong-host", "wrong-path", WebsocketScheme)
	suite.Error(err)
}

func (suite *websocketTestSuite) TestIsConnected() {
	suite.init()
	suite.True(suite.transport.IsConnected())
}

func (suite *websocketTestSuite) TestIsNotConnected() {
	suite.init().Disconnect()
	suite.False(suite.transport.IsConnected())
}

func (suite *websocketTestSuite) TestDisconnect() {
	suite.init()
	suite.True(suite.transport.IsConnected())
	suite.transport.Disconnect()
	suite.False(suite.transport.IsConnected())
}
