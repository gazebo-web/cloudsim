package simulator

import (
	"context"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

const (
	// actionNameStartSimulation is the name used to register the start simulation action.
	actionNameStartSimulation = "start-simulation"

	// actionNameStopSimulation is the name used to register the stop simulation action.
	actionNameStopSimulation = "stop-simulation"

	// applicationName is the name of the current simulator's application.
	applicationName = "subt"
)

// subTSimulator is a simulator.Simulator implementation.
type subTSimulator struct {
	applicationName string
	platform        platform.Platform
	services        subtapp.Services
	actions         actions.Servicer
	db              *gorm.DB
}

// Start triggers the action that will be in charge of launching a simulation with the given Group ID.
func (s *subTSimulator) Start(ctx context.Context, groupID simulations.GroupID) error {
	state := state.NewStartSimulation(s.platform, s.services, groupID)
	store := actions.NewStore(state)

	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameStartSimulation,
	}

	err := s.actions.Execute(store, s.db, execInput, groupID)
	if err != nil {
		return err
	}

	return nil
}

// Stop triggers the action that will be in charge of stopping a simulation with the given Group ID.
func (s *subTSimulator) Stop(ctx context.Context, groupID simulations.GroupID) error {
	state := state.NewStopSimulation(s.platform, s.services, groupID)
	store := actions.NewStore(state)

	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameStopSimulation,
	}

	err := s.actions.Execute(store, s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Config is used to initialize a new simulator for SubT.
type Config struct {
	DB                  *gorm.DB
	Platform            platform.Platform
	ApplicationServices subtapp.Services
	ActionService       actions.Servicer
}

// NewSimulator initializes a new Simulator implementation for SubT.
func NewSimulator(config Config) simulator.Simulator {
	registerActions(applicationName, config.ActionService)
	return &subTSimulator{
		platform:        config.Platform,
		applicationName: applicationName,
		services:        config.ApplicationServices,
		actions:         config.ActionService,
	}
}

// registerActions register a set of actions into the given service with the given application's name.
// It panics whenever an action could not be registered.
func registerActions(name string, service actions.Servicer) {
	actions := map[string]actions.Jobs{
		actionNameStartSimulation: JobsStartSimulation,
		actionNameStopSimulation:  JobsStopSimulation,
	}

	for actionName, jobs := range actions {
		err := registerAction(name, service, actionName, jobs)
		if err != nil {
			panic(err)
		}
	}
}

// registerAction registers the given jobs as a new action called actionName.
// The action gets registered into the given service for the given application name.
func registerAction(applicationName string, service actions.Servicer, actionName string, jobs actions.Jobs) error {
	action, err := actions.NewAction(jobs)
	if err != nil {
		return err
	}

	err = service.RegisterAction(&applicationName, actionName, action)
	if err != nil {
		return err
	}
	return nil
}
