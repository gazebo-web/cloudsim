package requester

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type restartRequester struct{}

func (req *restartRequester) Do(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (req *restartRequester) Validate(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func NewRestartRequester() Requester {
	return &restartRequester{}
}
