package platform

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestInitializers_Logger(t *testing.T) {
	p := Platform{}
	p.initializeLogger()
	var interfaceType ign.Logger
	assert.Implements(t, &interfaceType, p.Logger)
}

func TestInitializers_Context(t *testing.T) {
	p := Platform{}
	p.initializeLogger()
	p.initializeContext()
	var interfaceType context.Context
	assert.Implements(t, &interfaceType, p.Context)
}

func TestInitializers_Server(t *testing.T) {
}

func TestInitializers_Router(t *testing.T) {

}

func TestInitializers_Email(t *testing.T) {

}