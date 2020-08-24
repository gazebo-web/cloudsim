package requester

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type shutdownRequester struct{}

func (req *shutdownRequester) Do(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (req *shutdownRequester) Validate(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func NewShutdownRequester() Requester {
	return &shutdownRequester{}
}
