package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSynchronic_Serve(t *testing.T) {
	flag := false

	synchronicPool, err := NewSynchronicPool(1, func(i interface{}) {
		value, ok := i.(bool)
		assert.True(t, ok)
		flag = value
	})
	assert.NoError(t, err)
	synchronicPool.Serve(!flag)
	assert.Equal(t, true, flag)
	synchronicPool.Serve(!flag)
	assert.Equal(t, false, flag)
}