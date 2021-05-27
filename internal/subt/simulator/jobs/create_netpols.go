package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// CreateNetworkPolicyGazeboServer extends the generic jobs.CreateNetworkPolicies to create a network policy for the
// gazebo server
var CreateNetworkPolicyGazeboServer = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:            "create-netpol-gzserver",
	PreHooks:        []actions.JobFunc{setStartState, prepareNetworkPolicyGazeboServerInput},
	PostHooks:       []actions.JobFunc{checkCreateNetworkPoliciesError, returnState},
	RollbackHandler: removeCreatedNetworkPolicyGazeboServer,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// removeCreatedNetworkPolicyGazeboServer removes the created network policies for a gzserver in case an error is thrown.
func removeCreatedNetworkPolicyGazeboServer(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameGazeboServer(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_ = s.Platform().Orchestrator().NetworkPolicies().Remove(name, ns)

	return nil, nil
}

// checkCreateNetworkPoliciesError checks that the output generated by CreateNetworkPolicies has returned no errors.
func checkCreateNetworkPoliciesError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.CreateNetworkPoliciesOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return out, nil
}

// prepareNetworkPolicyGazeboServerInput is a pre-hook of the CreateNetworkPolicyGazeboServer job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
func prepareNetworkPolicyGazeboServerInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	// All robots
	selectors := make([]resource.Selector, len(robots))

	// Each robot's comms bridge will be granted with permissions to communicate to the gazebo server
	for i, r := range robots {
		selectors[i] = subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r)
	}

	// Mapping server
	selectors = append(selectors, subtapp.GetPodLabelsMappingServer(s.GroupID, s.ParentGroupID))

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
			Ingresses: network.IngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic to websocket server
					9002,
				},
			},
			Egresses: network.EgressRule{
				IPBlocks: []string{
					// Allow traffic to cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				// Allow traffic to the internet
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
	PostHooks:       []actions.JobFunc{checkCreateNetworkPoliciesError, returnState},
	RollbackHandler: removeCreatedNetworkPoliciesFieldComputer,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// removeCreatedNetworkPoliciesFieldComputer removes the created network policies for all field computers in case an error is thrown.
func removeCreatedNetworkPoliciesFieldComputer(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for i := range robots {
		robotID := subtapp.GetRobotID(i)
		name := subtapp.GetPodNameFieldComputer(s.GroupID, robotID)
		ns := s.Platform().Store().Orchestrator().Namespace()

		_ = s.Platform().Orchestrator().NetworkPolicies().Remove(name, ns)
	}

	return nil, nil
}

// prepareNetworkPolicyFieldComputersInput is a pre-hook of the CreateNetworkPolicyFieldComputers job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
func prepareNetworkPolicyFieldComputersInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	input := make(jobs.CreateNetworkPoliciesInput, len(robots))

	// Each robot's field computer should communicate to its respective comms bridge.
	for i, r := range robots {
		robotID := subtapp.GetRobotID(i)
		input[i] = network.CreateNetworkPolicyInput{
			Name:        subtapp.GetPodNameFieldComputer(s.GroupID, robotID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			PeersFrom: []resource.Selector{
				// Allow traffic from comms bridges
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			PeersTo: []resource.Selector{
				// Allow traffic to comms bridges
				subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			},
			Ingresses: network.IngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic from websocket server
					9002,
				},
			},
			Egresses: network.EgressRule{
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
	Name:            "create-netpol-comms-bridges",
	PreHooks:        []actions.JobFunc{setStartState, prepareNetworkPolicyCommsBridgesInput},
	PostHooks:       []actions.JobFunc{checkCreateNetworkPoliciesError, returnState},
	RollbackHandler: removeCreatedNetworkPoliciesCommsBridge,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// removeCreatedNetworkPoliciesCommsBridge removes the created network policies for all comms bridge in case an error is thrown.
func removeCreatedNetworkPoliciesCommsBridge(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for i := range robots {
		robotID := subtapp.GetRobotID(i)
		name := subtapp.GetPodNameCommsBridge(s.GroupID, robotID)
		ns := s.Platform().Store().Orchestrator().Namespace()

		_ = s.Platform().Orchestrator().NetworkPolicies().Remove(name, ns)
	}

	return nil, nil
}

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
		input[i] = network.CreateNetworkPolicyInput{
			Name:        subtapp.GetPodNameCommsBridge(s.GroupID, robotID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r),
			PeersFrom: []resource.Selector{
				// Allow traffic from gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				// Allow traffic from field computer
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			PeersTo: []resource.Selector{
				// Allow traffic to gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
				// Allow traffic to field computer
				subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID),
			},
			Ingresses: network.IngressRule{
				IPBlocks: []string{
					// Allow traffic from cloudsim
					fmt.Sprintf("%s/32", s.Platform().Store().Ignition().IP()),
				},
				Ports: []int32{
					// Allow traffic from websocket server
					9002,
				},
			},
			Egresses: network.EgressRule{
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

// CreateNetworkPolicyMappingServer extends the generic jobs.CreateNetworkPolicies to create a network policy for the
// mapping server.
var CreateNetworkPolicyMappingServer = jobs.CreateNetworkPolicies.Extend(actions.Job{
	Name:            "create-netpol-mapping-server",
	PreHooks:        []actions.JobFunc{setStartState, prepareNetworkPolicyMappingServerInput},
	PostHooks:       []actions.JobFunc{checkCreateNetworkPoliciesError, returnState},
	RollbackHandler: removeCreatedNetworkPoliciesMappingServer,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// removeCreatedNetworkPoliciesMappingServer removes the created network policies for the mapping server in case an error is thrown.
func removeCreatedNetworkPoliciesMappingServer(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for i := range robots {
		robotID := subtapp.GetRobotID(i)
		name := subtapp.GetPodNameCommsBridge(s.GroupID, robotID)
		ns := s.Platform().Store().Orchestrator().Namespace()

		_ = s.Platform().Orchestrator().NetworkPolicies().Remove(name, ns)
	}

	return nil, nil
}

// prepareNetworkPolicyMappingServerInput is a pre-hook of the CreateNetworkPolicyMappingServer job that prepares the input
// for the generic jobs.CreateNetworkPolicies job.
func prepareNetworkPolicyMappingServerInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	return jobs.CreateNetworkPoliciesInput{
		network.CreateNetworkPolicyInput{
			Name:        subtapp.GetPodNameMappingServer(s.GroupID),
			Namespace:   s.Platform().Store().Orchestrator().Namespace(),
			Labels:      subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID).Map(),
			PodSelector: subtapp.GetPodLabelsMappingServer(s.GroupID, s.ParentGroupID),
			PeersFrom: []resource.Selector{
				// Allow traffic from gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
			},
			PeersTo: []resource.Selector{
				// Allow traffic to gazebo server
				subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID),
			},
			Ingresses: network.IngressRule{
				Ports: []int32{
					// Allow traffic from websocket server
					// TBD: [ign port],
				},
			},
			Egresses: network.EgressRule{
				// Allow traffic to the internet
				AllowOutbound: true,
			},
		},
	}, nil
}
