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
	}, nil)
	s.Require().NoError(err)
}

func (s *testRoundRobinSuite) TestSelectors() {
	selectors := s.platformRoundRobin.Selectors()
	expected := []string{s.selector1, s.selector2, s.selector3}
	s.Require().ElementsMatch(expected, selectors)
}

func (s *testRoundRobinSuite) TestPlatformsNoSelector() {
	platforms := s.platformRoundRobin.Platforms(nil)
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
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

func (s *testRoundRobinSuite) TestPlatformsValidSelector() {
	platforms := s.platformRoundRobin.Platforms(&s.selector2)
	s.Require().Equal(s.platform2, platforms[0])
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
}

func (s *testRoundRobinSuite) TestPlatformsInvalidSelector() {
	selector := "invalid"
	platforms := s.platformRoundRobin.Platforms(&selector)
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
}

func (s *testRoundRobinSuite) TestPlatformValidSelector() {
	// Get the first platform
	platform, err := s.platformRoundRobin.Platform(s.selector1)
	s.Require().NoError(err)
	s.Require().Equal(s.platform1, platform)

	// Get the third platform
	platform, err = s.platformRoundRobin.Platform(s.selector3)
	s.Require().NoError(err)
	s.Require().Equal(s.platform3, platform)
}

func (s *testRoundRobinSuite) TestPlatformInvalidSelector() {
	// Provide an invalid selector
	selector := "invalid"
	platform, err := s.platformRoundRobin.Platform(selector)
	s.Assert().EqualError(err, ErrPlatformNotFound.Error())
	s.Assert().Nil(platform)
}
