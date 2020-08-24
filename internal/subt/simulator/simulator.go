package simulator

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

// subTSimulator is a simulator.Simulator implementation.
type subTSimulator struct {
	applicationName string
	platform        platform.Platform
	services        application.Services
	actions         actions.Servicer
	db              *gorm.DB
}

// Start triggers the action that will be in charge of launching a simulation with the given Group ID.
func (s *subTSimulator) Start(ctx context.Context, groupID simulations.GroupID) error {
	ctx = context.WithValue(ctx, contextPlatform, s.platform)
	ctx = context.WithValue(ctx, contextServices, s.services)
	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      "startSimulation",
	}
	err := s.actions.Execute(s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Stop triggers the action that will be in charge of stopping a simulation with the given Group ID.
func (s *subTSimulator) Stop(ctx context.Context, groupID simulations.GroupID) error {
	ctx = context.WithValue(ctx, contextPlatform, s.platform)
	ctx = context.WithValue(ctx, contextServices, s.services)
	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      "stopSimulation",
	}
	err := s.actions.Execute(s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Restart triggers the action that will be in charge of restarting a simulation with the given Group ID.
func (s *subTSimulator) Restart(ctx context.Context, groupID simulations.GroupID) error {
	ctx = context.WithValue(ctx, contextPlatform, s.platform)
	ctx = context.WithValue(ctx, contextServices, s.services)
	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      "restartSimulation",
	}
	err := s.actions.Execute(s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Config is used to initialize a new simulator for SubT.
type Config struct {
	DB                  *gorm.DB
	ApplicationName     string
	Platform            platform.Platform
	ApplicationServices application.Services
	ActionService       actions.Servicer
}

// NewSimulator initializes a new Simulator implementation for SubT.
func NewSimulator(config Config) simulator.Simulator {
	// [START SIMULATION] Create action
	startSimulation, err := actions.NewAction(JobsStartSimulation)
	if err != nil {
		panic(err)
	}
	// [START SIMULATION] Register action
	err = config.ActionService.RegisterAction(&config.ApplicationName, "startSimulation", startSimulation)
	if err != nil {
		panic(err)
	}

	// [STOP SIMULATION] Create action
	stopSimulation, err := actions.NewAction(JobsStopSimulation)
	if err != nil {
		panic(err)
	}
	// [STOP SIMULATION] Register action
	err = config.ActionService.RegisterAction(&config.ApplicationName, "stopSimulation", stopSimulation)
	if err != nil {
		panic(err)
	}

	// [RESTART SIMULATION] Create action
	restartSimulation, err := actions.NewAction(JobsRestartSimulation)
	if err != nil {
		panic(err)
	}
	// [RESTART SIMULATION] Register action
	err = config.ActionService.RegisterAction(&config.ApplicationName, "restartSimulation", restartSimulation)
	if err != nil {
		panic(err)
	}

	return &subTSimulator{
		platform:        config.Platform,
		applicationName: config.ApplicationName,
		services:        config.ApplicationServices,
		actions:         config.ActionService,
	}
}
