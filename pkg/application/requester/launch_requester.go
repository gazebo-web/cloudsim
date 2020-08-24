package requester

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type launchRequester struct{}

func (req *launchRequester) Do(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (req *launchRequester) Validate(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func NewLaunchRequester() Requester {
	return &launchRequester{}
}
