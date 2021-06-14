package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// GetMappingServerIP is a job in charge of getting the IP from the simulation's mapping server pod.
// WaitForMappingServerPod should be run before running this job.
var GetMappingServerIP = &actions.Job{
	Name:       "get-mapping-server-pod-ip",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    getMappingServerIP,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
}

// getGazeboIP gets the mapping server pod IP and assigns it to the start simulation state.
func getMappingServerIP(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	ip, err := s.Platform().Orchestrator().Pods().GetIP(application.GetPodNameMappingServer(s.GroupID), s.Platform().Store().Orchestrator().Namespace())
	if err != nil {
		return nil, err
	}

	s.MappingServerIP = ip
	store.SetState(s)

	return s, nil
}
