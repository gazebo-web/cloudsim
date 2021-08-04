package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subtsims "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// returnState is an actions.JobFunc implementation that returns the state. It's usually used as a posthook.
func returnState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	return store.State(), nil
}

// setStartState parses the input value as the StartSimulation state and sets the store with that state.
func setStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)
	store.SetState(s)
	return s, nil
}

// setStopState parses the input value as the StopSimulation state and sets the store with that state.
func setStopState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StopSimulation)
	store.SetState(s)
	return s, nil
}

// returnGroupIDFromStartState parses the input value as the StartSimulation state and returns the group id.
func returnGroupIDFromStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	return s.GroupID, nil
}

// returnGroupIDFromStopState parses the input vale as the StopSimulation state and returns the group id.
func returnGroupIDFromStopState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)
	return s.GroupID, nil
}

// checkWaitError is an actions.JobFunc implementation that checks if the value returned by a job extended from the
// Wait job returns an error.
// If the previous jobs.Wait job returns an error, it will trigger the rollback handler.
// If the previous jobs.Wait does not return an error, it will pass to the next actions.JobFunc.
func checkWaitError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	output := value.(jobs.WaitOutput)
	if output.Error != nil {
		return nil, output.Error
	}
	return nil, nil
}

// checkLaunchPodsError is an actions.JobFunc implementation meant to be used as a post hook that checks if the value
// returned by the job that launches pods returns an error.
func checkLaunchPodsError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	output := value.(jobs.LaunchPodsOutput)
	if output.Error == nil {
		return value, nil
	}
	return nil, output.Error
}

// checkLaunchServiceError checks if the output from the jobs.LaunchWebsocketService has an error.
func checkLaunchServiceError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	output := value.(jobs.LaunchWebsocketServiceOutput)
	if output.Error == nil {
		return value, nil
	}
	return nil, output.Error
}

// readFileContentFromCopyPod reads the file content located in the given path
// of a certain pod in the given namespace.
func readFileContentFromPod(p pods.Pods, podName, namespace, path string) ([]byte, error) {
	res, err := p.Get(podName, namespace)
	if err != nil {
		return nil, err
	}

	// TODO: Change container once we add more containers to the different simulation pods.
	buff, err := p.Reader(res).File("", path)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func isMappingServerEnabled(svc subtapp.Services, groupID simulations.GroupID) bool {
	// Get simulation
	sim, err := svc.Simulations().Get(groupID)
	if err != nil {
		return false
	}

	// Parse to subt simulation
	subtSim := sim.(subtsims.Simulation)

	// Get track
	track, err := svc.Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), subtSim.GetRunIndex())
	if err != nil {
		return false
	}

	return track.MappingImage != nil
}
