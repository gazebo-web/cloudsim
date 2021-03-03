package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// GetCommsBridgePodIP is a job in charge of getting the IP from the simulation's comms bridge pods.
var GetCommsBridgePodIP = &actions.Job{
	Name:       "get-comms-bridge-pod-ip",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    getCommsBridgeIP,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
}

// getGazeboIP gets the gazebo server pod IP and assigns it to the start simulation state.
func getCommsBridgeIP(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	ips := make([]string, len(robots))
	for i := range robots {
		robotID := application.GetRobotID(i)

		name := application.GetPodNameCommsBridge(s.GroupID, robotID)
		ns := s.Platform().Store().Orchestrator().Namespace()

		ip, err := s.Platform().Orchestrator().Pods().GetIP(name, ns)
		if err != nil {
			return nil, err
		}

		ips[i] = ip
	}

	s.CommsBridgeIPs = ips

	store.SetState(s)

	return s, nil
}
