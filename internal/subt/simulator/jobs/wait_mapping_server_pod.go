package jobs

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForMappingServerPod waits for the simulation mapping server pod to finish launching.
var WaitForMappingServerPod = jobs.Wait.Extend(actions.Job{
	Name:       "wait-mapping-server-pod",
	PreHooks:   []actions.JobFunc{createWaitRequestMappingServerPod},
	PostHooks:  []actions.JobFunc{checkMappingServerWaitError, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestMappingServerPod is the pre hook in charge of passing the needed input to the Wait job.
func createWaitRequestMappingServerPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	name := subtapp.GetPodNameMappingServer(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()
	labels := subtapp.GetPodLabelsMappingServer(s.GroupID, s.ParentGroupID)

	res := resource.NewResource(name, ns, labels)

	// Create wait for condition request
	req := s.Platform().Orchestrator().Pods().WaitForCondition(res, resource.HasIPStatusCondition)

	// Get timeout and poll frequency from store
	timeout := s.Platform().Store().Orchestrator().Timeout()
	pollFreq := s.Platform().Store().Orchestrator().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}

// checkMappingServerWaitError validates the Wait return value.
// If the track does not have a mapping server image defined, the Missing Pods errors is ignored.
// In any other case, an error will result in a rollback.
func checkMappingServerWaitError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	// Get simulation
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Parse to subt simulation
	subtSim := sim.(simulations.Simulation)

	// Get track
	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), subtSim.GetRunIndex())
	if err != nil {
		return nil, err
	}

	output := value.(jobs.WaitOutput)

	// Ignore missing pods error if mapping image is not defined.
	if track.MappingImage == nil && errors.Is(output.Error, pods.ErrMissingPods) {
		return jobs.LaunchPodsInput{}, nil
	} else if output.Error != nil {
		return nil, output.Error
	}

	return nil, nil
}
