package simulator

import (
	"context"
	"github.com/jinzhu/gorm"
	npsapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
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

// Config is used to configure the NPS simulator when calling NewSimulatorNPS.
type Config struct {
	DB                  *gorm.DB
	Platform            platform.Platform
	ApplicationServices npsapp.Services
	ActionService       actions.Servicer
}

// NewSimulatorNPS initializes a simulator for NPS.
func NewSimulatorNPS(config Config) simulator.Simulator {
	return &nps{}
}
