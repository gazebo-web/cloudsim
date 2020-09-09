package simulator

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	simctx "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

const (
	// actionNameStartSimulation is the name used to register the start simulation action.
	actionNameStartSimulation = "start-simulation"

	// actionNameStopSimulation is the name used to register the stop simulation action.
	actionNameStopSimulation = "stop-simulation"

	// actionNameRestartSimulation is the name used to register the restart simulation action.
	actionNameRestartSimulation = "restart-simulation"

	// applicationName is the name of the current simulator's application.
	applicationName = "subt"
)

// subTSimulator is a simulator.Simulator implementation.
type subTSimulator struct {
	applicationName string
	platform        platform.Platform
	services        application.Services
	actions         actions.Servicer
	db              *gorm.DB
	logger          ign.Logger
}

// Start triggers the action that will be in charge of launching a simulation with the given Group ID.
func (s *subTSimulator) Start(ctx context.Context, groupID simulations.GroupID) error {
	ctx = s.setupContext(ctx)

	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameStartSimulation,
	}
	err := s.actions.Execute(ctx, s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Stop triggers the action that will be in charge of stopping a simulation with the given Group ID.
func (s *subTSimulator) Stop(ctx context.Context, groupID simulations.GroupID) error {
	ctx = s.setupContext(ctx)
	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameStopSimulation,
	}
	err := s.actions.Execute(ctx, s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Restart triggers the action that will be in charge of restarting a simulation with the given Group ID.
func (s *subTSimulator) Restart(ctx context.Context, groupID simulations.GroupID) error {
	ctx = s.setupContext(ctx)
	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameRestartSimulation,
	}
	err := s.actions.Execute(ctx, s.db, execInput, groupID)
	if err != nil {
		return err
	}
	return nil
}

// setupContext is in charge of setting up the context for jobs.
func (s *subTSimulator) setupContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, simctx.CtxPlatform, s.platform)
	ctx = context.WithValue(ctx, simctx.CtxServices, s.services)
	return ctx
}

// Config is used to initialize a new simulator for SubT.
type Config struct {
	DB                  *gorm.DB
	Platform            platform.Platform
	ApplicationServices application.Services
	ActionService       actions.Servicer
	Logger              ign.Logger
}

// NewSimulator initializes a new Simulator implementation for SubT.
func NewSimulator(config Config) simulator.Simulator {
	registerActions(applicationName, config.ActionService)
	return &subTSimulator{
		platform:        config.Platform,
		applicationName: applicationName,
		services:        config.ApplicationServices,
		actions:         config.ActionService,
		logger:          config.Logger,
	}
}

// registerActionInput is used as the input for the registerAction function.
type registerActionInput struct {
	Jobs  actions.Jobs
	Store actions.Store
}

// registerActions register a set of actions into the given service with the given application's name.
// It panics whenever an action could not be registered.
func registerActions(name string, service actions.Servicer) {
	actions := map[string]registerActionInput{
		actionNameStartSimulation: {
			Jobs:  JobsStartSimulation,
			Store: fake.NewFakeStore(new(startSimulationData)),
		},
		actionNameStopSimulation: {
			Jobs:  JobsStopSimulation,
			Store: nil,
		},
		actionNameRestartSimulation: {
			Jobs:  JobsRestartSimulation,
			Store: nil,
		},
	}

	for actionName, input := range actions {
		err := registerAction(name, service, actionName, input)
		if err != nil {
			panic(err)
		}
	}
}

// registerAction registers the given jobs as a new action called actionName.
// The action gets registered into the given service for the given application name.
func registerAction(applicationName string, service actions.Servicer, actionName string, input registerActionInput) error {
	action, err := actions.NewAction(input.Jobs)
	if err != nil {
		return err
	}

	err = service.RegisterAction(&applicationName, actionName, action, input.Store)
	if err != nil {
		return err
	}
	return nil
}

// startSimulationData has all the information needed to start a simulation.
// It's used as the data type for the action's store.
type startSimulationData struct {
	GroupID             simulations.GroupID
	GazeboServerPodName string
	MachineList         []string
	GazeboServerPodIP   string
	BaseLabels          map[string]string
	GazeboLabels        map[string]string
	BridgeLabels        map[string]string
	FieldComputerLabels map[string]string
}
