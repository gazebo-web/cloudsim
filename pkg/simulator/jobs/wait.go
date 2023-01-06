package jobs

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/actions"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulator"
	"github.com/gazebo-web/cloudsim/v4/pkg/waiter"
	"github.com/jinzhu/gorm"
	"time"
)

// WaitInput is the input of the Wait job.
type WaitInput struct {
	Request       waiter.Waiter
	PollFrequency time.Duration
	Timeout       time.Duration
}

// WaitOutput is the output of the Wait job.
type WaitOutput struct {
	Error error
}

// Wait is a job that is in charge of waiting for a certain process to happen.
var Wait = &actions.Job{
	Execute: wait,
}

// wait is the Wait execute function. It's used to trigger the Wait method in the Request passed inside the WaitInput
// value with the given WaitInput.Timeout and WaitInput.PollFrequency.
// It returns an error if the request fails.
// wait will be used for any resource or event that implements the waiter interface.
// Examples:
//
//	Waiting for nodes to be registered in the cluster
//	Waiting for pods to have an ip assigned.
//	Waiting for pods to be on the "Ready" state.
func wait(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// If value is nil, bypass the job.
	if value == nil {
		return WaitOutput{
			Error: nil,
		}, nil
	}

	input, ok := value.(WaitInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	err := input.Request.Wait(input.Timeout, input.PollFrequency)
	return WaitOutput{
		Error: err,
	}, nil
}
