package manager

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"testing"
)

func TestRoundRobinSuite(t *testing.T) {
	suite.Run(t, new(testRoundRobinSuite))
}

type testRoundRobinSuite struct {
	suite.Suite
	selector1          string
	selector2          string
	selector3          string
	platform1          platform.Platform
	platform2          platform.Platform
	platform3          platform.Platform
	platformRoundRobin Manager
}

func (s *testRoundRobinSuite) SetupSuite() {
	s.selector1 = "test_1"
	s.platform1, _ = platform.NewPlatform("p1", platform.Components{})

	s.selector2 = "test_2"
	s.platform2, _ = platform.NewPlatform("p2", platform.Components{})

	s.selector3 = "test_3"
	s.platform3, _ = platform.NewPlatform("p3", platform.Components{})

	var err error
	s.platformRoundRobin, err = WithRoundRobin(Map{
		s.selector1: s.platform1,
		s.selector2: s.platform2,
		s.selector3: s.platform3,
	})
	s.Require().NoError(err)
}

func (s *testRoundRobinSuite) TestPlatformsNoSelectorMultipleCalls() {
	platforms := s.platformRoundRobin.Platforms(nil)
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Assert().ElementsMatch(expected, platforms)

	platforms = s.platformRoundRobin.Platforms(nil)
	expected = []platform.Platform{s.platform2, s.platform3, s.platform1}
	s.Assert().ElementsMatch(expected, platforms)

	platforms = s.platformRoundRobin.Platforms(nil)
	expected = []platform.Platform{s.platform3, s.platform1, s.platform2}
	s.Assert().ElementsMatch(expected, platforms)
}
