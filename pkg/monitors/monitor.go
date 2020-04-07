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
	Ticker *time.Ticker
	Done chan bool
}

func New(d time.Duration) *Monitor {
	return &Monitor{
		Ticker: time.NewTicker(d),
		Done:   make(chan bool, 1),
	}
}

func GetRunner(baseCtx context.Context, id string, name string, monitor *Monitor, job Job) Runner {
	newLogger := logger.Logger(baseCtx).Clone(id)
	ctx := ign.NewContextWithLogger(baseCtx, newLogger)

	return func() {
		for {
			select {
			case <-monitor.Done:
				newLogger.Info(fmt.Sprintf("%s is done", name))
				return
			case <-monitor.Ticker.C:
				_ = job(ctx)
			}
		}
	}
}