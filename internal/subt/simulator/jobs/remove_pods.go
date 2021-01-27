package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemovePods extends the generic jobs.RemovePods job. It's in charge of removing simulation pods.
var RemovePods = jobs.RemovePods.Extend(actions.Job{
	Name:       "remove-pods",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemovePodsInput},
	PostHooks:  []actions.JobFunc{checkRemovePodsNoError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// checkRemovePodsNoError is a post-hook in charge of checking that no errors were thrown while removing pods.
func checkRemovePodsNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)
	out := value.(*jobs.RemovePodsOutput)
	if out.Error != nil || len(out.Resources) != len(s.PodList) {
		err := deployment.SetJobData(tx, nil, actions.DeploymentJobData, out)
		if err != nil {
			return nil, err
		}
		return nil, out.Error
	}
	return nil, nil
}

// prepareRemovePodsInput is a pre-hook in charge of setting up the selector needed for the generic jobs to delete pods.
func prepareRemovePodsInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	robots, err := s.SubTServices().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	ns := s.Platform().Store().Orchestrator().Namespace()

	// The max amount of pods is given by 3 pods per robot (fc, comms, copy) + gzserver + gzserver copy pod
	list := make([]orchestrator.Resource, 0, 3*len(robots)+2)

	// Add robot-related pods
	for i := range robots {
		robotID := subtapp.GetRobotID(i)

		// Field computer
		list = append(list, orchestrator.NewResource(subtapp.GetPodNameFieldComputer(s.GroupID, robotID), ns, nil))

		// Comms bridge
		list = append(list, orchestrator.NewResource(subtapp.GetPodNameCommsBridge(s.GroupID, robotID), ns, nil))

		// And if logs are enabled, copy pod for comms bridge.
		if s.Platform().Store().Ignition().LogsCopyEnabled() {
			list = append(list, orchestrator.NewResource(subtapp.GetPodNameCommsBridgeCopy(s.GroupID, robotID), ns, nil))
		}
	}

	// Gazebo server
	list = append(list, orchestrator.NewResource(subtapp.GetPodNameGazeboServer(s.GroupID), ns, nil))

	// And if logs are enabled, gazebo server copy pod.
	if s.Platform().Store().Ignition().LogsCopyEnabled() {
		list = append(list, orchestrator.NewResource(subtapp.GetPodNameGazeboServerCopy(s.GroupID), ns, nil))
	}

	s.PodList = list
	store.SetState(s)

	return jobs.RemovePodsInput(list), nil
}
