package localstack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_NotNil(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
}
