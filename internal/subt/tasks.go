package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/monitors"
	"time"
)

// RegisterTasks returns an array of the tasks that need to be executed by the scheduler.
func (app *SubT) RegisterTasks() []monitors.Task {
	return []monitors.Task{
		{
			Job:  func() {

			},
			Date: time.Time{},
		},
	}
}