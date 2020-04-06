package platform

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"os"
	"testing"
)

func TestInitializers_Logger(t *testing.T) {
	p := Platform{}
	p.setupLogger()
	var interfaceType ign.Logger
	assert.Implements(t, &interfaceType, p.Logger)
}

func TestInitializers_Context(t *testing.T) {
	p := Platform{}
	p.setupLogger()
	p.setupContext()
	var interfaceType context.Context
	assert.Implements(t, &interfaceType, p.Context)
}

func TestInitializers_Server(t *testing.T) {
	p := Platform{}
	os.Setenv("IGN_CLOUDSIM_HTTP_PORT", "80")
	os.Setenv("IGN_CLOUDSIM_SSL_PORT", "445")
	p.Config = NewConfig()
	p.setupLogger()
	p.setupContext()
	p.setupServer()
	assert.NotNil(t, p.Server)
}

func TestInitializers_Router(t *testing.T) {
	p := Platform{}
	os.Setenv("IGN_CLOUDSIM_HTTP_PORT", "80")
	os.Setenv("IGN_CLOUDSIM_SSL_PORT", "445")
	p.Config = NewConfig()
	p.setupLogger()
	p.setupContext()
	p.setupServer()
	p.setupRouter()
	assert.NotNil(t, p.Server.Router)
}

func TestInitializers_Email(t *testing.T) {

}
