package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckStartSimulationIsNotParent is a job in charge of checking if a simulation is not a parent simulation.
var CheckStartSimulationIsNotParent = GenerateCheckSimulationNotOfKindJob(
	"check-start-simulation-no-parent",
	simulations.SimParent,
	&state.StartSimulation{},
	&state.StartSimulation{},
)

// CheckStopSimulationIsNotParent is a job in charge of checking if a simulation is not a parent simulation.
var CheckStopSimulationIsNotParent = GenerateCheckSimulationNotOfKindJob(
	"check-stop-simulation-no-parent",
	simulations.SimParent,
	&state.StartSimulation{},
	&state.StartSimulation{},
)
