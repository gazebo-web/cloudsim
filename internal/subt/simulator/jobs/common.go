package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	cstate "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
	"time"
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

// returnGroupIDFromStartState parses the input value as the StartSimulation state and returns the group id.
func returnGroupIDFromStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
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
	err := deployment.SetJobData(tx, nil, actions.DeploymentJobData, value)
	if err != nil {
		return nil, err
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

// configCreatePodInput is a set of configurations that need to be passed in order to configure a orchestrator.CreatePodInput using
// the prepareCreatePodInput function.
type configCreatePodInput struct {
	name                      string
	namespace                 string
	labels                    map[string]string
	restartPolicy             orchestrator.RestartPolicy
	terminationGracePeriod    time.Duration
	nodeSelector              orchestrator.Selector
	containerName             string
	image                     string
	command                   []string
	args                      []string
	privileged                bool
	allowPrivilegesEscalation bool
	ports                     []int32
	volumes                   []orchestrator.Volume
	envVars                   map[string]string
	nameservers               []string
}

// prepareCreatePodInput is in charge of preparing the input for the create pod job.
func prepareCreatePodInput(c configCreatePodInput) orchestrator.CreatePodInput {
	return orchestrator.CreatePodInput{
		Name:                          c.name,
		Namespace:                     c.namespace,
		Labels:                        c.labels,
		RestartPolicy:                 c.restartPolicy,
		TerminationGracePeriodSeconds: c.terminationGracePeriod,
		NodeSelector:                  c.nodeSelector,
		Containers: []orchestrator.Container{
			{
				Name:                     c.containerName,
				Image:                    c.image,
				Command:                  c.command,
				Args:                     c.args,
				Privileged:               &c.privileged,
				AllowPrivilegeEscalation: &c.allowPrivilegesEscalation,
				Ports:                    c.ports,
				Volumes:                  c.volumes,
				EnvVars:                  c.envVars,
			},
		},
		Volumes:     c.volumes,
		Nameservers: c.nameservers,
	}
}
