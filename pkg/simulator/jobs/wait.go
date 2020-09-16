package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"time"
)

// WaitInput is the input of the Wait job.
type WaitInput struct {
	GroupID       simulations.GroupID
	Request       waiter.Waiter
	PollFrequency time.Duration
	Timeout       time.Duration
}

// WaitOutput is the output of the Wait job.
type WaitOutput simulations.GroupID

// Wait is a job that is in charge of waiting for a certain process to happen.
var Wait = &actions.Job{
	Execute: wait,
}

// wait is the Wait execute function. It's used to trigger the Wait method in the Request passed inside the WaitInput
// value with the given WaitInput.Timeout and WaitInput.PollFrequency.
// It returns an error if the request fails.
func wait(ctx actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input, ok := value.(WaitInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	err := input.Request.Wait(input.Timeout, input.PollFrequency)
	if err != nil {
		return nil, err
	}

	return WaitOutput(input.GroupID), nil
}
