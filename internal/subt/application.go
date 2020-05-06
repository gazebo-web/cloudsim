package subt

import (
	"context"
	"fmt"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IApplication interface {
	application.IApplication
	isSimulationHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg
}

// SubT is an IApplication implementation
type SubT struct {
	application.IApplication
	Services    services
	Controllers controllers
}

type controllers struct {
	Simulation sim.IController
}

type services struct {
	application.Services
	Simulation sim.IService
}

// New creates a new SubT application.
func New(p *platform.Platform) IApplication {
	simulationRepository := sim.NewRepository(p.Server.Db)
	simulationService := sim.NewService(simulationRepository)
	baseApp := application.New(p, simulationService.Parent(), p.UserService)

	subt := &SubT{
		IApplication: baseApp,
		Controllers: controllers{
			Simulation: sim.NewController(sim.NewControllerInput{
				Service:     simulationService,
				Decoder:     p.FormDecoder,
				Validator:   p.Validator,
				Permissions: p.Permissions,
				UserService: p.UserService,
			}),
		},
		Services: services{
			Services: application.Services{
				User: p.UserService,
			},
			Simulation: simulationService,
		},
	}

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
func (app *SubT) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	if err := app.IApplication.ValidateLaunch(ctx, simulation); err != nil {
		return err
	}

	if err := app.isSimulationHeld(ctx, simulation); err != nil {
		logger.Logger(ctx).Warning(fmt.Sprintf("[LAUNCH|VALIDATE] Cannot run a held simulation. Group ID: [%s]", *simulation.GroupID))
		return err
	}

	return nil
}

// isSimulationHeld checks if the simulations is being held.
func (app *SubT) isSimulationHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	if simulation.Held {
		// TODO: Replace with ign.NewErrorMessage(ign.ErrorLaunchHeld).
		return ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus)
	}
	return nil
}

// Register runs a set of instructions to initialize an application for the given platform.
func Register(p *platform.Platform) application.IApplication {
	return New(p)
}
