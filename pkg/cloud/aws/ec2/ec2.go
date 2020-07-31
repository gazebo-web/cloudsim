package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
)

type machines struct {
	API ec2iface.EC2API
}

func (e machines) Create(input cloud.CreateMachinesInput) error {
	panic("implement me")
}

func (e machines) Terminate(input cloud.TerminateMachinesInput) error {
	panic("implement me")
}

func (e machines) Count(input cloud.CountMachinesInput) int {
	panic("implement me")
}

func NewMachines(api ec2iface.EC2API) cloud.Machines {
	return &machines{
		API: api,
	}
}
