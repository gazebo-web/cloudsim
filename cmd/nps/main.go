package main

import (
	"context"
	npsapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/simulator"
	gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"os"
)

func main() {
	logger := ign.NewLoggerNoRollbar("NPS", ign.VerbosityDebug)

	if err := run(logger); err != nil {
		logger.Error("main: error:", err)
		os.Exit(1)
	}
}

func run(logger ign.Logger) error {
	// Database ---
	logger.Debug("main: Initializing database connection")
	db, err := gorm.GetDBFromEnvVars()
	if err != nil {
		return err
	}

	// Queue ---
	startQueue := ign.NewQueue()
	stopQueue := ign.NewQueue()

	// Simulations ---
	simulationRepository := gormrepo.NewRepository(db, logger, &simulations.Simulation{})
	simulationService := simulations.NewService(simulationRepository, startQueue, stopQueue, logger)
	simulationController := simulations.NewController(simulationService)

	// Users & Permissions ---
	perm := &permissions.Permissions{}
	err = perm.Init(db, "sysadmin")
	if err != nil {
		return err
	}

	userService, err := users.NewService(context.TODO(), perm, db, "sysadmin")
	if err != nil {
		return err
	}

	// Router ---
	logger.Debug("main: Initializing router")
	router := ign.NewRouter()
	routerConfig := ign.NewRouterConfigurer(router, nil)

	routerConfig.ConfigureRouter("/1.0/simulations", simulationController.GetRoutes())

	// Platform ---
	logger.Debug("main: Initializing NPS cloudsim platform")
	p := platform.NewPlatform(platform.Components{
		Machines: nil,
		Storage:  nil,
		Cluster:  nil,
		Store:    nil,
		Secrets:  nil,
	})

	// Application services
	base := application.NewServices(simulationService, userService)
	services := npsapp.NewServices(base)

	// Simulator ---
	logger.Debug("main: Initializing NPS simulator")
	sim := simulator.NewSimulatorNPS(simulator.Config{
		DB:                  db,
		Platform:            p,
		ApplicationServices: services,
		ActionService:       actions.NewService(),
	})

	// API Server ---
	logger.Debug("main: Initializing API server")
	s := server.NewServer(server.Config{
		Router:     router,
		DB:         db,
		Logger:     logger,
		Simulator:  sim,
		StartQueue: startQueue,
		StopQueue:  stopQueue,
	})

	// HTTP listener ---
	logger.Debug("main: Listening on port :3030")
	err = s.ListenAndServe(":3030")
	if err != nil {
		return err
	}

	return nil
}
