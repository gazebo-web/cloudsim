package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"testing"
	"time"
)

func TestWait_InvalidWaitRequest(t *testing.T) {
	assert.Panics(t, func() {
		Wait.Run(nil, nil, nil, WaitInput{
			Timeout:       time.Minute,
			PollFrequency: time.Second,
		})
	})
}

func TestWait_InvalidTimeout(t *testing.T) {
	result, _ := Wait.Run(nil, nil, nil, WaitInput{
		Request: waiter.NewWaitRequest(func() (bool, error) {
			return true, nil
		}),
		Timeout:       -1,
		PollFrequency: time.Second,
	})
	output, ok := result.(WaitOutput)
	assert.True(t, ok)
	assert.Error(t, output.Error)
}

func TestWait_InvalidPollFrequency(t *testing.T) {
	result, _ := Wait.Run(nil, nil, nil, WaitInput{
		Request: waiter.NewWaitRequest(func() (bool, error) {
			return true, nil
		}),
		Timeout:       time.Second,
		PollFrequency: -1,
	})
	output, ok := result.(WaitOutput)
	assert.True(t, ok)
	assert.Error(t, output.Error)
}
