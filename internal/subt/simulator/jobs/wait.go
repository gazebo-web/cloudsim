package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"time"
)

// WaitInput is the input used by the Wait execute function.
type WaitInput struct {
	GroupID       simulations.GroupID
	Request       waiter.Waiter
	PollFrequency time.Duration
	Timeout       time.Duration
}

// Wait is a job that is in charge of waiting for a certain process to happen.
var Wait = &actions.Job{
	Name:       "wait",
	Execute:    wait,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// wait is the Wait execute function. It's used to trigger the Wait method in the Request passed inside the WaitInput
// value with the given WaitInput.Timeout and WaitInput.PollFrequency.
// It returns an error if the request fails.
func wait(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	input, ok := value.(WaitInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	err := input.Request.Wait(input.Timeout, input.PollFrequency)
	if err != nil {
		return nil, err
	}

	return input.GroupID, nil
}
