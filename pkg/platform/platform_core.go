package platform

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/monitors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/workers"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IPlatformCore represents a set of methods to start, stop, restart and reload the application.
type IPlatformCore interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	Reload(ctx context.Context) error
}

// Start starts the platform.
func (p *Platform) Start(ctx context.Context) error {
	go func() {
		for {
			var element interface{}
			var err *ign.ErrMsg
			if element, err = p.LaunchQueue.DequeueOrWait(); err != nil {
				continue
			}

			dto, ok := element.(workers.LaunchDTO)
			if !ok {
				continue
			}

			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] About to process launch action. Group ID: [%s]", dto.GroupID))
			if err := p.LaunchPool.Serve(dto); err != nil {
				p.Logger.Error(fmt.Sprintf("[QUEUE|LAUNCH] Error while serving launch action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] The launch action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
		}
	}()

	go func() {
		for dto := range p.TerminationQueue {
			p.Logger.Info(fmt.Sprintf("[QUEUE|TERMINATE] About to process terminate action. Group ID: [%s]", dto.GroupID))
			if err := p.TerminationPool.Serve(dto); err != nil {
				p.Logger.Error(fmt.Sprintf("[QUEUE|TERMINATE] Error while serving terminate action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger.Info(fmt.Sprintf("[QUEUE|TERMINATE] The terminate action was successfully served to the worker pool. Group ID: [%s]", dto.GroupID))
		}
	}()

	// TODO: Rebuild state


	cleanerRunner := monitors.NewRunner(
		ctx,
		p.Cleaner,
		// TODO: Add checkForExpiredSimulations
		func(ctx context.Context) error { return nil },
	)
	go cleanerRunner()

	updaterRunner := monitors.NewRunner(
		ctx,
		p.Updater,
		// TODO: Add updateMultiSimStatuses
		func(ctx context.Context) error { return nil },
	)
	go updaterRunner()

	// TODO: Register tasks for the scheduler
	return nil
}

// Stop stops the platform.
func (p *Platform) Stop(ctx context.Context) error {
	p.Updater.Ticker.Stop()
	p.Cleaner.Ticker.Stop()
	close(p.TerminationQueue)
	return nil
}
