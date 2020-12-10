package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationIsParent is a job in charge of checking if a simulation is parent.
var CheckSimulationIsParent = GenerateCheckSimulationKindJob(
	"check-simulation-parent",
	simulations.SimParent,
	&state.StartSimulation{},
	&state.StartSimulation{},
)
