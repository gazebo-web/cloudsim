package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// GetGazeboIP is a job in charge of getting the IP from the simulation's gazebo server pod.
// WaitForGazeboServerPod should be run before running this job.
var GetGazeboIP = &actions.Job{
	Name:       "get-gzserver-pod-ip",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    getGazeboIP,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
}

// getGazeboIP gets the gazebo server pod IP and assigns it to the start simulation state.
func getGazeboIP(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	ip, err := s.Platform().Orchestrator().Pods().GetIP(application.GetPodNameGazeboServer(s.GroupID), s.Platform().Store().Orchestrator().Namespace())
	if err != nil {
		return nil, err
	}

	s.GazeboServerIP = ip
	store.SetState(s)

	return s, nil
}
