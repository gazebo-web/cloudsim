package monitors

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"time"
)

// Runner represents a function that will trigger a monitor to run a job.
type Runner func()

// Job represents the set of instructions to be executed by the monitor.
type Job func(ctx context.Context) error

// Monitor is in charge of executing a job.
type Monitor struct {
	UUID string
	Name string
	Ticker *time.Ticker
	Done chan bool
}

// New creates a new Monitor.
func New(uuid, name string, d time.Duration) *Monitor {
	return &Monitor{
		UUID: uuid,
		Name: name,
		Ticker: time.NewTicker(d),
		Done:   make(chan bool, 1),
	}
}

// NewRunner creates a new runner by the given monitor and job.
func NewRunner(baseCtx context.Context, monitor *Monitor, job Job) Runner {
	newLogger := logger.Logger(baseCtx).Clone(monitor.UUID)
	ctx := ign.NewContextWithLogger(baseCtx, newLogger)

	return func() {
		for {
			select {
			case <-monitor.Done:
				newLogger.Info(fmt.Sprintf("[RUNNER] %s is done", monitor.Name))
				return
			case <-monitor.Ticker.C:
				_ = job(ctx)
			}
		}
	}
}