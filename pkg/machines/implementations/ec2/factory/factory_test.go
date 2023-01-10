package factory

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/machines"
	"github.com/gazebo-web/cloudsim/v4/pkg/machines/implementations/ec2"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestEC2FactorySuite(t *testing.T) {
	suite.Run(t, new(testEC2FactorySuite))
}

type testEC2FactorySuite struct {
	suite.Suite
}

func (s *testEC2FactorySuite) TestInitializeAPIDependencyIsNil() {
	config := Config{
		Region: "us-east-1",
	}
	dependencies := Dependencies{}
	s.Nil(initializeAPI(&config, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testEC2FactorySuite) TestInitializeAPIDependencyIsNotNil() {
	// Prepare dependencies
	ec2API := struct {
		ec2iface.EC2API
	}{}
	dependencies := Dependencies{
		API: ec2API,
	}

	s.Nil(initializeAPI(nil, &dependencies))
	s.Exactly(ec2API, dependencies.API)
}

func (s *testEC2FactorySuite) TestNewFuncDefaultConfig() {
	config := Config{
		Region: "test",
		Zones: []ec2.Zone{
			{
				Zone:     "test",
				SubnetID: "subnet-0123456789abcdefg",
			},
		},
	}

	// Prepare dependencies
	ec2API := struct {
		ec2iface.EC2API
	}{}
	dependencies := factory.Dependencies{
		"api":    ec2API,
		"logger": gz.NewLoggerNoRollbar("test", gz.VerbosityWarning),
	}

	var out machines.Machines
	s.Require().NoError(NewFunc(config, dependencies, &out))
}
