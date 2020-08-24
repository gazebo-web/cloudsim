package requester

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type startRequester struct{}

func (req *startRequester) Do(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (req *startRequester) Validate(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func NewStartRequester() Requester {
	return &startRequester{}
}
