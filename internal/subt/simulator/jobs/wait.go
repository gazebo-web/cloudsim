package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"time"
)

type WaitInput struct {
	GroupID       simulations.GroupID
	Request       waiter.Waiter
	PollFrequency time.Duration
	Timeout       time.Duration
}

var Wait = &actions.Job{
	Name:       "wait",
	Execute:    wait,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

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
