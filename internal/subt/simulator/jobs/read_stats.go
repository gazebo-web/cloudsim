package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gopkg.in/yaml.v2"
)

// ReadStats is a job in charge of reading the statistics from a gzserver copy pod for the simulation that is being processed.
var ReadStats = &actions.Job{
	Name:       "read-simulation-stats",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    readStats,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// readStats is the main execute function for the ReadStats job.
func readStats(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	// Get file path
	path := fmt.Sprintf("%s/summary.yml", s.Platform().Store().Ignition().SidecarContainerLogsPath())

	body, err := readFileContentFromPod(
		s.Platform().Orchestrator().Pods(),
		subtapp.GetPodNameGazeboServerCopy(s.GroupID),
		s.Platform().Store().Orchestrator().Namespace(),
		path,
	)
	if err != nil {
		return nil, err
	}

	// Parse statistics using yaml
	var stats simulations.Statistics
	err = yaml.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	// Set the stats to the store
	s.Stats = stats
	store.SetState(s)

	// Return the state
	return s, nil
}
