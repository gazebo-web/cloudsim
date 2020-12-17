package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

var CreateNetworkPolicyGazeboServer = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:       "create-netpol-gzserver",
	PreHooks:   []actions.JobFunc{setStartState, prepareNetworkPolicyGazeboServerInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareNetworkPolicyGazeboServerInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	selectors := make([]orchestrator.Selector, len(robots))

	for _, r := range robots {
		selectors = append(selectors, subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r))
	}

	return jobs.CreateNetworkPoliciesInput{
		{
			Name:        subtapp.GetPodNameGazeboServer(s.GroupID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
			PeersFrom:   selectors,
			PeersTo:     selectors,
			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				AllowOutbound: true,
			},
		},
	}, nil
}

var CreateNetworkPolicyFieldComputers = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:            "create-netpol-field-computers",
	PreHooks:        []actions.JobFunc{setStartState, prepareNetworkPolicyFieldComputersInput},
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareNetworkPolicyFieldComputersInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	input := make(jobs.CreateNetworkPoliciesInput, len(robots))

	for i, r := range robots {
		robotID := subtapp.GetRobotID(i)
		input[i] = orchestrator.CreateNetworkPolicyInput{
			Name:        subtapp.GetPodNameFieldComputer(s.GroupID, robotID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			PeersFrom: []orchestrator.Selector{
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			PeersTo: []orchestrator.Selector{
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				AllowOutbound: false,
			},
		}
	}

	return input, nil
}

var CreateNetworkPolicyCommsBridges = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:       "create-netpol-comms-bridges",
	PreHooks:   []actions.JobFunc{setStartState, prepareNetworkPolicyCommsBridgesInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareNetworkPolicyCommsBridgesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	input := make(jobs.CreateNetworkPoliciesInput, len(robots))

	for i, r := range robots {
		robotID := subtapp.GetRobotID(i)
		input[i] = orchestrator.CreateNetworkPolicyInput{
			Name:        subtapp.GetPodNameFieldComputer(s.GroupID, robotID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			PeersFrom: []orchestrator.Selector{
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			PeersTo: []orchestrator.Selector{
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				AllowOutbound: true,
			},
		}
	}

	return input, nil
}
