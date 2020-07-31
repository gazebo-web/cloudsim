package ec2

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"testing"
)

func TestMachines(t *testing.T) {
	suite.Run(t, new(ec2MachinesTestSuite))
}

type ec2MachinesTestSuite struct {
	suite.Suite
	session  *session.Session
	ec2API   ec2iface.EC2API
	machines cloud.Machines
}

func (s *ec2MachinesTestSuite) SetupTest() {
	s.session = session.Must(session.NewSession())
	s.ec2API = ec2.New(s.session)
	s.machines = NewMachines(s.ec2API)

}

func (s *ec2MachinesTestSuite) TestNewMachines() {
	e, ok := s.machines.(*machines)
	s.True(ok)
	s.NotNil(e.API)
	s.IsType(&ec2.EC2{}, e.API)
}
