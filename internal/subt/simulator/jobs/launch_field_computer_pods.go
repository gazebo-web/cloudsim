package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchFieldComputerPods launches the list of field computer pods.
var LaunchFieldComputerPods = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-field-computer-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareFieldComputerPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareFieldComputerPodInput prepares the input for the generic LaunchPods job to launch field computer pods.
func prepareFieldComputerPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	pods := make([]orchestrator.CreatePodInput, len(subtSim.GetRobots()))

	for i, r := range subtSim.GetRobots() {
		robotID := subtapp.GetRobotID(i)
		// Create field computer input
		privileged := false
		allowPrivilegesEscalation := true
		pods[i] = orchestrator.CreatePodInput{
			Name:                          subtapp.GetPodNameFieldComputer(s.GroupID, robotID),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 orchestrator.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			Containers: []orchestrator.Container{
				{
					Name:                     subtapp.GetContainerNameFieldComputer(),
					Image:                    r.GetImage(),
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegesEscalation,
					EnvVars:                  subtapp.GetEnvVarsFieldComputer(r.GetName(), s.CommsBridgeIPs[i]),
				},
			},
		}
	}

	return jobs.LaunchPodsInput(pods), nil
}
