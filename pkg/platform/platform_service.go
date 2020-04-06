package platform

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/workers"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IPlatformService represents a set of methods to perform on the platform.
type IPlatformService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Start starts the platform.
func (p *Platform) Start(ctx context.Context) error {
	go func() {
		var element interface{}
		var err *ign.ErrMsg
		for {
			if element, err = p.LaunchQueue.DequeueOrWait(); err != nil {
				continue
			}

			dto, ok := element.(workers.LaunchDTO)
			if !ok {
				continue
			}

			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] About to process launch action. Group ID: [%s]", dto.GroupID))
			if err := p.LaunchPool.Serve(dto); err != nil {
				p.Logger.Error(fmt.Sprintf("[QUEUE|ERROR] Error while launching action. Group ID: [%s]. Error: [%v]", dto.GroupID, err))
				continue
			}
			p.Logger.Info(fmt.Sprintf("[QUEUE|LAUNCH] The launch action was successfully launched to the workers pool. Group ID: [%s]", dto.GroupID))
		}
	}()



	return nil
}

// Stop stops the platform.
func (p *Platform) Stop(ctx context.Context) error {
	return nil
}
