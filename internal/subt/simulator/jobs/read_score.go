package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"strconv"
)

// ReadScore is a job in charge of reading the score from a gzserver copy pod for the simulation that is being processed.
var ReadScore = actions.Job{
	Name:       "read-simulation-score",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    readScore,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// readScore is the main execute function for the ReadScore job.
func readScore(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	// Get the gzserver copy pod
	res, err := s.Platform().Orchestrator().Pods().Get(subtapp.GetPodNameGazeboServerCopy(s.GroupID), s.Platform().Store().Orchestrator().Namespace())
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/score.yml", s.Platform().Store().Ignition().GazeboServerLogsPath())

	// Get a reader to read the score from the gzserver pod
	reader, err := s.Platform().Orchestrator().Pods().Reader(res).File(path)
	if err != nil {
		return nil, err
	}

	// Read the file's body
	var body []byte
	_, err = reader.Read(body)
	if err != nil {
		return nil, err
	}

	// Parse the score
	score, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return nil, err
	}

	// Set the score
	s.Score = score
	store.SetState(s)

	// Return state
	return s, nil
}
