package statistics

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// Service groups a set of methods to handle statistics for simulations.
type Service interface {
	// Save saves the given score and Statistics for the given simulation's GroupID.
	Save(groupID simulations.GroupID, score *float64, stats simulations.Statistics) error
}
