package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"time"
)

// LaunchFieldComputers launches the list of field computer pods.
var LaunchFieldComputers = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-field-computer-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareFieldComputerPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareFieldComputerPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {

	}

	subtSim := sim.(subt.Simulation)

	var pods []orchestrator.CreatePodInput

	for i, r := range subtSim.GetRobots() {
		robotID := subtapp.GetRobotID(i)
		// Create field computer input
		pods = append(pods, prepareFieldComputerCreatePodInput(configFieldComputerPod{
			groupID:                s.GroupID,
			robotID:                robotID,
			namespace:              s.Platform().Store().Orchestrator().Namespace(),
			labels:                 subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID).Map(),
			terminationGracePeriod: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			nodeSelector:           subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			containerImage:         subtSim.GetImage(),
			robotName:              r.Name(),
			nameservers:            s.Platform().Store().Orchestrator().Nameservers(),
		}))
	}

	return pods, nil
}

type configFieldComputerPod struct {
	groupID                simulations.GroupID
	robotID                string
	namespace              string
	labels                 map[string]string
	terminationGracePeriod time.Duration
	nodeSelector           orchestrator.Selector
	containerImage         string
	robotName              string
	nameservers            []string
}

func prepareFieldComputerCreatePodInput(c configFieldComputerPod) orchestrator.CreatePodInput {
	in := configPod{
		name:                      subtapp.GetPodNameFieldComputer(c.groupID, c.robotID),
		namespace:                 c.namespace,
		labels:                    c.labels,
		restartPolicy:             orchestrator.RestartPolicyNever,
		terminationGracePeriod:    c.terminationGracePeriod,
		nodeSelector:              c.nodeSelector,
		containerName:             "field-computer",
		image:                     c.containerImage,
		args:                      nil,
		privileged:                false,
		allowPrivilegesEscalation: true,
		ports:                     nil,
		volumes:                   nil,
		envVars: map[string]string{
			"ROBOT_NAME": c.robotName,
		},
		nameservers: nil,
	}

	return preparePod(in)
}
