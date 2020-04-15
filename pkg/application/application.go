package application

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/monitors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"time"
)

// IApplication describes a set of methods for an Application.
type IApplication interface {
	Name() string
	Version() string
	RegisterRoutes() ign.Routes
	RegisterTasks() []tasks.Task
	RegisterMonitors()
	RebuildState(ctx context.Context) error
	Shutdown(ctx context.Context)
}

// Application is a generic implementation of an application to be extended by a specific application.
type Application struct {
	Platform *platform.Platform
	services services
	Cleaner	 *monitors.Monitor
	Updater	 *monitors.Monitor
}

type services struct {
	simulation simulations.IService
}

// New creates a new application for the given platform.
func New(p *platform.Platform) *Application {
	app := &Application{
		Platform: p,
		Cleaner: monitors.New("expired-simulations-cleaner", "Expired Simulations Cleaner", 20 * time.Second),
		Updater: monitors.New("multisim-status-updater", "MultiSim Parent Status Updater", time.Minute),
	}
	return app
}

// Name returns the application's name.
// Needs to be implemented by the specific application.
func (app *Application) Name() string {
	panic("Name should be implemented by the application")
}

// Version returns the application's version.
// If the specific application doesn't implement this method, it will return 1.0.
func (app *Application) Version() string {
	return "1.0"
}

// RegisterRoutes returns the slice of the application's routes.
// Needs to be implemented by the specific application.
func (app *Application) RegisterRoutes() ign.Routes {
	panic("RegisterRoutes should be implemented by the application")
}

// RegisterTasks returns an array of the tasks that need to be executed by the scheduler.
// If the specific application doesn't implement this method, it will return an empty slice.
func (app *Application) RegisterTasks() []tasks.Task {
	return []tasks.Task{}
}

func (app *Application) RegisterMonitors(ctx context.Context) {
	cleanerRunner := monitors.NewRunner(
		ctx,
		app.Cleaner,
		// TODO: Add checkForExpiredSimulations
		func(ctx context.Context) error { return nil },
	)
	go cleanerRunner()

	updaterRunner := monitors.NewRunner(
		ctx,
		app.Updater,
		// TODO: Add updateMultiSimStatuses
		func(ctx context.Context) error { return nil },
	)
	go updaterRunner()
}

func (app *Application) Shutdown(ctx context.Context) {
	app.Updater.Ticker.Stop()
	app.Cleaner.Ticker.Stop()
}

func (app *Application) RebuildState(ctx context.Context) error {
	err := app.Platform.Simulator.Recover(ctx, app.getLabel, app.getGazeboConfig)
	if err != nil {
		return err
	}

	app.Platform.Simulator.RLock()
	defer app.Platform.Simulator.RUnlock()

	var sims simulations.Simulations
	// if err := db.Model(&SimulationDeployment{}).Where("error_status IS NULL").Where("multi_sim != ?", multiSimParent).
	//		Where("deployment_status BETWEEN ? AND ?", int(simPending), int(simTerminatingInstances)).Find(&deps).Error; err != nil {
	//		return err
	//	}

	for _, sim := range sims {
		switch sim.GetStatus() {
		case simulations.StatusPending:
			logger.Logger(ctx).Info(fmt.Sprintf("[APP|REBUILDING] Resuming launch process. GroupID: [%s]", *sim.GroupID))
			if err := app.Platform.LaunchQueue.Enqueue(); err != nil {
				logger.Logger(ctx).Error(fmt.Sprintf("[APP|REBUILDING] Error while launching simulation. GroupID: [%s]", *sim.GroupID))
			}
			continue
		case simulations.StatusRunning:
			running := app.Platform.Simulator.GetRunningSimulation(*sim.GroupID)
			if running != nil {
				logger.Logger(ctx).Info(fmt.Sprintf("[APP|RECOVER] The expected running simulation doesn't have any node running. GroupID: [%s]. Marking with error.", *sim.GroupID))
				sim.ErrorStatus = simulations.ErrServerRestart.ToStringPtr()
				if _, err := app.services.simulation.Update(*sim.GroupID, sim); err != nil {
					logger.Logger(ctx).Error(fmt.Sprintf("[APP|REBUILDING] Error while updating simulation. GroupID: [%s]", *sim.GroupID))
				}
			}
			continue
		default:
			logger.Logger(ctx).Info(fmt.Sprintf("[APP|REBUILDING] Simulation found with intermediate Status: [%s]. GroupID: [%s]. Marking with error.", sim.GetStatus().ToString(), *sim.GroupID))
			sim.ErrorStatus = simulations.ErrServerRestart.ToStringPtr()
			if _, err := app.services.simulation.Update(*sim.GroupID, sim); err != nil {
				logger.Logger(ctx).Error(fmt.Sprintf("[APP|REBUILDING] Error while updating simulation. GroupID: [%s]", *sim.GroupID))
			}
		}
	}
	return nil
}

func (app *Application) getLabel() *string {
	return nil
}

func (app *Application) getGazeboConfig() simulator.GazeboConfig {
	return simulator.GazeboConfig{
		WorldStatsTopic:  "",
		WorldWarmupTopic: "",
		MaxSeconds:       0,
	}
}