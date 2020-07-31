package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
)

// machines is a cloud.Machines implementation.
type machines struct {
	API ec2iface.EC2API
}

// Create creates EC2 machines.
func (e machines) Create(input cloud.CreateMachinesInput) error {
	panic("implement me")
}

// Terminate terminates EC2 machines.
func (e machines) Terminate(input cloud.TerminateMachinesInput) error {
	panic("implement me")
}

// Count counts EC2 machines.
func (e machines) Count(input cloud.CountMachinesInput) int {
	panic("implement me")
}

// NewMachines initializes a new cloud.Machines implementation using EC2.
func NewMachines(api ec2iface.EC2API) cloud.Machines {
	return &machines{
		API: api,
	}
}
