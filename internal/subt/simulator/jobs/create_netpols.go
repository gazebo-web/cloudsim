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

// CreateNetworkPolicyGazeboServer extends the generic jobs.CreateNetworkPolicies to create a network policy for the
// gazebo server
var CreateNetworkPolicyGazeboServer = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:       "create-netpol-gzserver",
	PreHooks:   []actions.JobFunc{setStartState, prepareNetworkPolicyGazeboServerInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareNetworkPolicyGazeboServerInput is a pre-hook of the CreateNetworkPolicyGazeboServer job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
func prepareNetworkPolicyGazeboServerInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	selectors := make([]orchestrator.Selector, len(robots))

	for i, r := range robots {
		selectors[i] = subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r)
	}

	return jobs.CreateNetworkPoliciesInput{
		{
			Name:        subtapp.GetPodNameGazeboServer(s.GroupID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
			// Allow traffic from comms bridges
			PeersFrom: selectors,
			// Allow traffic to comms bridges
			PeersTo: selectors,

			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic to websocket server
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					// Allow traffic to cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				// Allow traffic from anywhere.
				AllowOutbound: true,
			},
		},
	}, nil
}

// CreateNetworkPolicyFieldComputers extends the generic jobs.CreateNetworkPolicies to create a network policy for the
// different field computers.
var CreateNetworkPolicyFieldComputers = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:            "create-netpol-field-computers",
	PreHooks:        []actions.JobFunc{setStartState, prepareNetworkPolicyFieldComputersInput},
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareNetworkPolicyFieldComputersInput is a pre-hook of the CreateNetworkPolicyFieldComputers job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
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
				// Allow traffic from comms bridges
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			PeersTo: []orchestrator.Selector{
				// Allow traffic to comms bridges
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic from websocket server
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					// Allow traffic to cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				// Disable traffic to the internet
				AllowOutbound: false,
			},
		}
	}

	return input, nil
}

// CreateNetworkPolicyCommsBridges extends the generic jobs.CreateNetworkPolicies to create a network policy for the
// different comms bridges.
var CreateNetworkPolicyCommsBridges = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:       "create-netpol-comms-bridges",
	PreHooks:   []actions.JobFunc{setStartState, prepareNetworkPolicyCommsBridgesInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareNetworkPolicyCommsBridgesInput is a pre-hook of the CreateNetworkPolicyCommsBridges job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
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
				// Allow traffic from gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				// Allow traffic from field computer
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			PeersTo: []orchestrator.Selector{
				// Allow traffic to gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				// Allow traffic to field computer
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			Ingresses: orchestrator.NetworkIngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic from websocket server
					9002,
				},
			},
			Egresses: orchestrator.NetworkEgressRule{
				IPBlocks: []string{
					// Allow traffic to cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				// Allow traffic to the internet
				AllowOutbound: true,
			},
		}
	}

	return input, nil
}
