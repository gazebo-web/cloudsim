package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationIsNotParent is a job in charge of checking if a simulation is parent.
var CheckSimulationIsNotParent = GenerateCheckSimulationNoKindJob(
	"check-simulation-parent",
	simulations.SimParent,
	&state.StartSimulation{},
	&state.StartSimulation{},
)
