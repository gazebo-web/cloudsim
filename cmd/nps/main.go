package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/server"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/simulator"
	gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
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

	// Router ---
	logger.Debug("main: Initializing router")
	router := ign.NewRouter()
	routerConfig := ign.NewRouterConfigurer(router, nil)

	routerConfig.ConfigureRouter("/1.0/simulations", simulationController.GetRoutes())

	// Simulator ---
	logger.Debug("main: Initializing NPS simulator")
	sim := simulator.NewSimulatorNPS()

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
