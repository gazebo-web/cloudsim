package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchWebsocketService is a job extending the generic jobs.LaunchWebsocketService to launch a websocket service running inside
// the gazebo server pod.
var LaunchWebsocketService = jobs.LaunchWebsocketService.Extend(actions.Job{
	Name:            "launch-websocket-service",
	PreHooks:        []actions.JobFunc{setStartState, prepareLaunchWebsocketServiceInput},
	PostHooks:       []actions.JobFunc{checkLaunchServiceError, returnState},
	RollbackHandler: rollbackLaunchWebsocketService,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchWebsocketService(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetServiceNameWebsocket(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_ = s.Platform().Orchestrator().Services().Remove(resource.NewResource(name, ns, nil))

	return nil, nil
}

// prepareLaunchWebsocketServiceInput is a pre-hook of LaunchWebsocketService in charge of preparing the input for the generic
// jobs.LaunchWebsocketService job.
func prepareLaunchWebsocketServiceInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	return jobs.LaunchWebsocketServiceInput{
		Name:          subtapp.GetServiceNameWebsocket(s.GroupID),
		Type:          "ClusterIP",
		Namespace:     s.Platform().Store().Orchestrator().Namespace(),
		ServiceLabels: subtapp.GetWebsocketServiceLabels(s.GroupID).Map(),
		TargetLabels:  subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID).Map(),
		Ports: map[string]int32{
			"websockets": 9002,
		},
	}, nil
}
