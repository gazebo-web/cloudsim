package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockApp struct {
	Application
}

type mockApp2 struct {
	Application
}

func (mockApp2) Name() string {
	return "test"
}

func TestName_PanicWhenNotImplemented(t *testing.T) {
	m := mockApp{}
	assert.Panics(t, func() {
		m.Name()
	})
}

func TestName_NotPanicWhenImplemented(t *testing.T) {
	m := mockApp2{}
	assert.NotPanics(t, func() {
		m.Name()
	})
}

func TestName_NameMatches(t *testing.T) {
	m := mockApp2{}
	assert.Equal(t, "test", m.Name())
}

func TestRegistry_RegisterApplication(t *testing.T) {
	apps := map[string]*IApplication{}
	RegisterApplications(apps, func() *IApplication {
		var app IApplication
		app = mockApp2{}
		return &app
	})
	assert.Len(t, apps, 1)
}
