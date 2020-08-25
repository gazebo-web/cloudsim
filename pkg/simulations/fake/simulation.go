package fake

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type fakeSimulation struct {
	groupID simulations.GroupID
	status  simulations.Status
	kind    simulations.Kind
	err     *simulations.Error
}

func (f fakeSimulation) Error() *simulations.Error {
	return f.err
}

func (f fakeSimulation) ToCreateMachinesInput() []cloud.CreateMachinesInput {
	instanceProfile := "instance-profile"
	return []cloud.CreateMachinesInput{
		{
			InstanceProfile: &instanceProfile,
			KeyName:         "secret-key",
			Type:            "machine.type",
			Image:           "busybox",
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   []string{"some-rule", "some-test-rule"},
			SubnetID:        "subnet-1fjq4378",
			Zone:            "us-east-1a",
			Tags: []cloud.Tag{
				{
					Resource: "instance",
					Map: map[string]string{
						"Name": fmt.Sprintf("fake-node-group-%s-gzserver", f.groupID),
					},
				},
			},
			InitScript: nil,
			Retries:    10,
		},
		{
			InstanceProfile: &instanceProfile,
			KeyName:         "secret-key",
			Type:            "machine.type",
			Image:           "busybox",
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   []string{"some-rule", "some-test-rule"},
			SubnetID:        "subnet-1fjq4378",
			Zone:            "us-east-1a",
			Tags: []cloud.Tag{
				{
					Resource: "instance",
					Map: map[string]string{
						"Name": fmt.Sprintf("fake-node-group-%s-fc", f.groupID),
					},
				},
			},
			InitScript: nil,
			Retries:    10,
		},
		{
			InstanceProfile: &instanceProfile,
			KeyName:         "secret-key",
			Type:            "machine.type",
			Image:           "busybox",
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   []string{"some-rule", "some-test-rule"},
			SubnetID:        "subnet-1fjq4378",
			Zone:            "us-east-1a",
			Tags: []cloud.Tag{
				{
					Resource: "instance",
					Map: map[string]string{
						"Name": fmt.Sprintf("fake-node-group-%s-comms", f.groupID),
					},
				},
			},
			InitScript: nil,
			Retries:    10,
		},
	}
}

func (f fakeSimulation) GroupID() simulations.GroupID {
	return f.groupID
}

func (f fakeSimulation) Status() simulations.Status {
	return f.status
}

func (f fakeSimulation) Kind() simulations.Kind {
	return f.kind
}

func NewSimulation(groupID simulations.GroupID, status simulations.Status, kind simulations.Kind, err *simulations.Error) simulations.Simulation {
	return &fakeSimulation{
		groupID: groupID,
		status:  status,
		kind:    kind,
		err:     err,
	}
}
