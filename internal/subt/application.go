package subt

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/circuits"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/robots"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/stats"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"path"
)

type IApplication interface {
	application.IApplication
	isSimulationHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg
	UploadSummary(sim *sim.Simulation, score stats.Score) *ign.ErrMsg
}

// SubT is an IApplication implementation
type SubT struct {
	application.IApplication
	Services    services
	Controllers controllers
	Validator   *validator.Validate
}

type controllers struct {
	Simulation sim.Controller
}

type services struct {
	application.Services
	Simulation sim.IService
	Circuit    circuits.Service
	Robot      robots.IService
}

// New creates a new SubT application.
func New(p *platform.Platform) IApplication {
	simulationRepository := sim.NewRepository(p.Server.Db, p.Name())
	simulationService := sim.NewService(simulationRepository)

	baseApp := application.New(p, simulationService, p.UserService)

	validate := validator.New()
	subt := &SubT{
		IApplication: baseApp,
		Services: services{
			Services: application.Services{
				Simulation: simulationService,
				User:       p.UserService,
			},
			Simulation: simulationService,
		},
		Controllers: controllers{
			Simulation: sim.NewController(sim.NewControllerInput{
				Service:     simulationService,
				Decoder:     p.FormDecoder,
				Validator:   validate,
				Permissions: p.Permissions,
				UserService: p.UserService,
			}),
		},
		Validator: validate,
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
	subt := New(p)
	return subt
}

func (app *SubT) RegisterValidators(ctx context.Context) {
	app.Validator.RegisterValidation("iscircuit", app.Services.Circuit.IsValidCircuit)
	app.Validator.RegisterValidation("isrobottype", app.Services.Robot.IsValidRobotType)
}

func (app *SubT) UploadSummary(sim *sim.Simulation, score stats.Score) *ign.ErrMsg {
	b, err := json.Marshal(score)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorMarshalJSON, err)
	}

	fileName := tools.GenerateSummaryFilename(*sim.GroupID)
	key := path.Join(app.Platform().CloudProvider.S3().GetLogKey(*sim.GroupID, *sim.Base.Owner), fileName)

	// TODO: Add AWS_GZ_LOGS_BUCKET env var.
	_, err = app.Platform().CloudProvider.S3().Upload("AWS_GZ_LOGS_BUCKET", key, b)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	return nil
}

func (app *SubT) UploadLogs() *ign.ErrMsg {
	return nil
}
