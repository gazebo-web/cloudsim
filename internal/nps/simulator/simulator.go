package simulator

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

// A similar implementation for this could be found in:
// internal/subt/simulator/simulator.go
type nps struct {
}

func (n *nps) Start(ctx context.Context, groupID simulations.GroupID) error {
	panic("todo: implement me")
}

func (n *nps) Stop(ctx context.Context, groupID simulations.GroupID) error {
	panic("todo: implement me")
}

// NewSimulatorNPS initializes a simulator for NPS.
func NewSimulatorNPS() simulator.Simulator {
	return &nps{}
}
