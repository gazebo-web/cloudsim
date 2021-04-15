package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForCommsBridgePodsReady waits for the simulation comms bridge pods to be ready.
var WaitForCommsBridgePodsReady = jobs.Wait.Extend(actions.Job{
	Name:       "wait-comms-bridge-pods-ready",
	PreHooks:   []actions.JobFunc{createWaitRequestForCommsBridgePodToBeReady},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestForCommsBridgePodToBeReady is the pre hook in charge of passing the needed input to the Wait job.
func createWaitRequestForCommsBridgePodToBeReady(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	res := resource.NewResource("", "", subtapp.GetPodLabelsBase(s.GroupID, nil))

	// Create wait for condition request
	// Since only the gazebo server pod has been created and already has an IP, we only need to wait until
	// comms bridge pods have an ip.
	req := s.Platform().Orchestrator().Pods().WaitForCondition(res, resource.ReadyCondition)

	// Get timeout and poll frequency from store
	timeout := s.Platform().Store().Machines().Timeout()
	pollFreq := s.Platform().Store().Machines().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
