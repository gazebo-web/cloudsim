package statistics

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type Service interface {
	// Save saves the given score and Statistics for the given GroupID.
	Save(groupID simulations.GroupID, score *float64, stats simulations.Statistics) error
}
