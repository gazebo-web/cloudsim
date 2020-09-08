package actions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewStore()
	assert.NotNil(t, s)
}
