package requester

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Requesters interface {
	Start() Requester
	Shutdown() Requester
	Launch() Requester
	Restart() Requester
}

type Requester interface {
	Do(payload interface{}) (interface{}, *ign.ErrMsg)
	Validate(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg
}
