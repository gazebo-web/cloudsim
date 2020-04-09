package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"time"
)

func (app *SubT) RegisterTasks() []tasks.Task {
	return []tasks.Task{
		{
			Job:  func() {

			},
			Date: time.Time{},
		},
	}
}