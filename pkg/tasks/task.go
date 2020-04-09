package tasks

import "time"

type Task struct {
	Job func()
	Date time.Time
}
