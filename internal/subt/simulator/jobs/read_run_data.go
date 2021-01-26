package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// ReadRunData is a job in charge of reading the run data from a gzserver copy pod for the simulation that is being processed.
var ReadRunData = actions.Job{
	Name:       "read-simulation-run-data",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    readRunData,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// readRunData is the main execute function for the ReadRunData job.
func readRunData(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	if !sim.IsKind(simulations.SimSingle) {
		return s, nil
	}

	path := fmt.Sprintf("%s/run.yml", s.Platform().Store().Ignition().GazeboServerLogsPath())

	body, err := readFileContentFromPod(
		s.Platform().Orchestrator().Pods(),
		subtapp.GetPodNameGazeboServerCopy(s.GroupID),
		s.Platform().Store().Orchestrator().Namespace(),
		path,
	)
	if err != nil {
		return nil, err
	}

	// Set run data
	s.RunData = string(body)
	store.SetState(s)

	// Return state
	return s, nil
}
