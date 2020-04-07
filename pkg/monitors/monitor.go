package monitors

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"time"
)

type Runner func()

type Job func(ctx context.Context) error

type Monitor struct {
	UUID string
	Name string
	Ticker *time.Ticker
	Done chan bool
}

func New(uuid, name string, d time.Duration) *Monitor {
	return &Monitor{
		UUID: uuid,
		Name: name,
		Ticker: time.NewTicker(d),
		Done:   make(chan bool, 1),
	}
}

func NewRunner(baseCtx context.Context, monitor *Monitor, job Job) Runner {
	newLogger := logger.Logger(baseCtx).Clone(monitor.UUID)
	ctx := ign.NewContextWithLogger(baseCtx, newLogger)

	return func() {
		for {
			select {
			case <-monitor.Done:
				newLogger.Info(fmt.Sprintf("%s is done", monitor.Name))
				return
			case <-monitor.Ticker.C:
				_ = job(ctx)
			}
		}
	}
}