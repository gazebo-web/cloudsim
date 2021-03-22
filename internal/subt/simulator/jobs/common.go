package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	cstate "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
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
	if err := deployment.SetJobData(tx, nil, actions.DeploymentJobData, output.Resources); err != nil {
		return nil, err
	}
	if output.Error == nil {
		return value, nil
	}
	return nil, output.Error
}

// rollbackPodCreation is an actions.JobErrorHandler implementation meant to be used as rollback handler to delete pods
// that were initialized in the jobs.LaunchPods job.
func rollbackPodCreation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, thrownError error) (interface{}, error) {
	out, err := deployment.GetJobData(tx, nil, actions.DeploymentJobData)
	if err != nil {
		return nil, err
	}

	s := store.State().(cstate.PlatformGetter)

	list := out.([]orchestrator.Resource)
	for _, pod := range list {
		_, _ = s.Platform().Orchestrator().Pods().Delete(pod)
	}

	return nil, nil
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
func readFileContentFromPod(p orchestrator.Pods, podName, namespace, path string) ([]byte, error) {
	res, err := p.Get(podName, namespace)
	if err != nil {
		return nil, err
	}

	reader, err := p.Reader(res).File(path)
	if err != nil {
		return nil, err
	}

	var c []byte
	_, err = reader.Read(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
