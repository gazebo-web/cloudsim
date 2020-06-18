package application

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/monitors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strings"
	"time"
)

// Application describes a set of methods for an application.
type Application interface {
	Name() string
	Version() string
	Platform() platform.Platform
	RegisterRoutes() ign.Routes
	RegisterTasks() []monitors.Task
	RegisterMonitors(ctx context.Context)
	RegisterValidators(ctx context.Context)
	RebuildState(ctx context.Context) error
	Stop(ctx context.Context) error

	Launch(payload interface{}) (interface{}, *ign.ErrMsg)
	ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg

	Shutdown(payload interface{}) (interface{}, *ign.ErrMsg)
	ValidateShutdown(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg

	LaunchHeld(payload interface{}) (interface{}, *ign.ErrMsg)
	ValidateLaunchHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg

	Restart(payload interface{}) (interface{}, *ign.ErrMsg)
	ValidateRestart(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg
}

// application is a generic implementation of an application to be extended by a specific application.
type application struct {
	platform platform.Platform
	Services Services
	Cleaner  *monitors.Monitor
	Updater  *monitors.Monitor
}

// Services group a list of services to be used by the application.
type Services struct {
	Simulation simulations.Service
	User       users.Service
}

// New creates a new application for the given platform.
func New(p platform.Platform, simulationService simulations.Service, userService users.Service) Application {
	app := &application{
		platform: p,
		Cleaner:  monitors.New("expired-simulations-cleaner", "Expired Simulations Cleaner", 20*time.Second),
		Updater:  monitors.New("multisim-status-updater", "MultiSim Parent Status Updater", time.Minute),
		Services: Services{
			Simulation: simulationService,
			User:       userService,
		},
	}
	return app
}

// Name returns the application's name.
// Needs to be implemented by the specific application.
func (app *application) Name() string {
	panic("Name should be implemented by the application")
}

// Version returns the application's version.
// If the specific application doesn't implement this method, it will return 1.0.
func (app *application) Version() string {
	return "1.0"
}

// platform returns the reference the application's platform.
func (app *application) Platform() platform.Platform {
	return app.platform
}

// RegisterRoutes returns the slice of the application's routes.
// Needs to be implemented by the specific application.
func (app *application) RegisterRoutes() ign.Routes {
	panic("RegisterRoutes should be implemented by the application")
}

// RegisterTasks returns an array of the tasks that need to be executed by the scheduler.
// If the specific application doesn't implement this method, it will return an empty slice.
func (app *application) RegisterTasks() []monitors.Task {
	return []monitors.Task{}
}

// RegisterMonitors runs the Cleaner Job and the Updater job.
func (app *application) RegisterMonitors(ctx context.Context) {
	cleanerRunner := monitors.NewRunner(
		ctx,
		app.Cleaner,
		func(ctx context.Context) error { return app.checkForExpiredSimulations() },
	)
	go cleanerRunner()

	updaterRunner := monitors.NewRunner(
		ctx,
		app.Updater,
		func(ctx context.Context) error { return app.updateMultiSimStatuses() },
	)
	go updaterRunner()
}

func (app *application) RegisterValidators(ctx context.Context) {
	return
}

// Stop executes a set of instructions to turn off the application.
func (app *application) Stop(ctx context.Context) error {
	app.Updater.Ticker.Stop()
	app.Cleaner.Ticker.Stop()
	return nil
}

// RebuildState runs a set of instructions to restore the application to the previous state before a restart.
func (app *application) RebuildState(ctx context.Context) error {
	err := app.Platform().Simulator.Recover(ctx, app.getLabel, app.getGazeboConfig)
	if err != nil {
		return err
	}

	app.Platform().Simulator.RLock()
	defer app.Platform().Simulator.RUnlock()

	var sims simulations.Simulations
	// if err := db.Model(&SimulationDeployment{}).Where("error_status IS NULL").Where("multi_sim != ?", multiSimParent).
	//		Where("deployment_status BETWEEN ? AND ?", int(simPending), int(simTerminatingInstances)).Find(&deps).Error; err != nil {
	//		return err
	//	}

	for _, sim := range sims {
		switch sim.GetStatus() {
		case simulations.StatusPending:
			logger.Logger(ctx).Info(fmt.Sprintf("[APP|REBUILDING] Resuming launch process. GroupID: [%s]", *sim.GroupID))
			app.Platform().RequestLaunch(ctx, *sim.GroupID)
			continue
		case simulations.StatusRunning:
			running := app.Platform().Simulator.GetRunningSimulation(*sim.GroupID)
			if running != nil {
				logger.Logger(ctx).Info(fmt.Sprintf("[APP|RECOVER] The expected running simulation doesn't have any node running. GroupID: [%s]. Marking with error.", *sim.GroupID))
				updateSim := simulations.SimulationUpdate{
					ErrorStatus: simulations.ErrServerRestart.ToStringPtr(),
				}
				if _, err := app.Services.Simulation.Update(ctx, *sim.GroupID, updateSim, nil); err != nil {
					logger.Logger(ctx).Error(fmt.Sprintf("[APP|REBUILDING] Error while updating simulation. GroupID: [%s]", *sim.GroupID))
				}
			}
			continue
		default:
			logger.Logger(ctx).Info(fmt.Sprintf("[APP|REBUILDING] Simulation found with intermediate Status: [%s]. GroupID: [%s]. Marking with error.", sim.GetStatus().ToString(), *sim.GroupID))
			updateSim := simulations.SimulationUpdate{
				ErrorStatus: simulations.ErrServerRestart.ToStringPtr(),
			}
			if _, err := app.Services.Simulation.Update(ctx, *sim.GroupID, updateSim, nil); err != nil {
				logger.Logger(ctx).Error(fmt.Sprintf("[APP|REBUILDING] Error while updating simulation. GroupID: [%s]", *sim.GroupID))
			}
		}
	}
	return nil
}

// getLabel returns the label that's being used to identify the application's running simulations.
func (app *application) getLabel() *string {
	return nil
}

// getGazeboConfig returns a GazeboConfig for the application.
func (app *application) getGazeboConfig(sim *simulations.Simulation) simulator.GazeboConfig {
	panic("getGazeboConfig should be implemented by the application.")
}

// LaunchSimulation receives a Simulation and requests a Launch to the platform.
func (app *application) Launch(payload interface{}) (interface{}, *ign.ErrMsg) {
	simulation := payload.(*simulations.Simulation)
	ctx := context.Background()

	sims, err := app.Services.Simulation.Prepare(ctx, simulation)

	if err != nil {
		return nil, err
	}

	for _, sim := range sims {
		groupID := *sim.GroupID
		logger.Logger(ctx).Info(
			fmt.Sprintf(
				"[%s|SIMULATIONS] About to submit launch task for GroupID: [%s]", strings.ToUpper(app.Name()), groupID,
			),
		)
		if err := app.ValidateLaunch(ctx, &sim); err != nil {
			return nil, err
		}
		app.Platform().RequestLaunch(ctx, *sim.GroupID)
	}

	return simulation, nil
}

// ValidateLaunch receives a simulation and performs a set of checks to
// ensure that the simulation is good to be launched.
func (app *application) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	return nil
}

func (app *application) Shutdown(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (app *application) ValidateShutdown(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func (app *application) LaunchHeld(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (app *application) ValidateLaunchHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}

func (app *application) Restart(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (app *application) ValidateRestart(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}
