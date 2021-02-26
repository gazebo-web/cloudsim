package jobs

import (
  "fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
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
// 		Waiting for nodes to be registered in the cluster
// 		Waiting for pods to have an ip assigned.
// 		Waiting for pods to be on the "Ready" state.
func wait(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\n\nwait\n\n")
	input, ok := value.(WaitInput)
	if !ok {
    fmt.Println(simulator.ErrInvalidInput)
		return nil, simulator.ErrInvalidInput
	}

	err := input.Request.Wait(input.Timeout, input.PollFrequency)
  fmt.Printf("\n\ndone waiting\n\n")
	return WaitOutput{
		Error: err,
	}, nil
}
