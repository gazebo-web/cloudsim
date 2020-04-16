package subt

import (
	"context"
	"errors"
	"fmt"
	// TODO: Find a way of avoiding the usage of .
	. "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

// SubT is an IApplication implementation
type SubT struct {
	*application.Application
}

// New creates a new SubT application.
func New(p *platform.Platform) *SubT {
	app := application.New(p)
	subt := &SubT{
		Application: app,
	}
	repository := NewRepository(p.Server.Db)
	app.Services.Simulation = NewService(repository)
	return subt
}

// Name returns the SubT application's name.
func (app *SubT) Name() string {
	return "subt"
}

// Version returns the SubT application's version.
func (app *SubT) Version() string {
	return "2.0"
}

func (app *SubT) getGazeboConfig(sim *simulations.Simulation) simulator.GazeboConfig {
	// GetCircuitRules
	return simulator.GazeboConfig{
		WorldStatsTopic:  "/",
		WorldWarmupTopic: "/",
		MaxSeconds:       0,
	}
}

// ValidateLaunch runs a set of checks before launching a simulation. It will return an error if one of those checks fail.
func (app *SubT) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) error {
	if err := app.isSimulationHeld(ctx, simulation); err != nil {
		logger.Logger(ctx).Warning(fmt.Sprintf("[LAUNCH|VALIDATE] Cannot run a held simulation. Group ID: [%s]", *simulation.GroupID))
		return err
	}
	return nil
}

// isSimulationHeld checks if the simulations is being held.
func (app *SubT) isSimulationHeld(ctx context.Context, simulation *simulations.Simulation) error {
	if simulation.Held {
		// TODO: Replace with ign.NewErrorMessage(ign.ErrorLaunchHeld).
		return errors.New("launch held simulation")
	}
	return nil
}

// Register runs a set of instructions to initialize an application for the given platform.
func Register(p *platform.Platform) application.IApplication {
	return New(p)
}