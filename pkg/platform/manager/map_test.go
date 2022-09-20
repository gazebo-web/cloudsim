package manager

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/loader"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	"testing"
)

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(testMapSuite))
}

type testMapSuite struct {
	suite.Suite
	selector1   string
	selector2   string
	selector3   string
	platform1   platform.Platform
	platform2   platform.Platform
	platform3   platform.Platform
	platformMap Map
}

func (s *testMapSuite) SetupSuite() {
	s.selector1 = "test_1"
	s.platform1, _ = platform.NewPlatform("p1", platform.Components{})

	s.selector2 = "test_2"
	s.platform2, _ = platform.NewPlatform("p2", platform.Components{})

	s.selector3 = "test_3"
	s.platform3, _ = platform.NewPlatform("p3", platform.Components{})

	s.platformMap = Map{
		s.selector1: s.platform1,
		s.selector2: s.platform2,
		s.selector3: s.platform3,
	}
}

func (s *testMapSuite) TestSelectors() {
	selectors := s.platformMap.Selectors()
	expected := []string{s.selector1, s.selector2, s.selector3}
	s.Require().ElementsMatch(expected, selectors)
}

func (s *testMapSuite) TestPlatformsNoSelector() {
	platforms := s.platformMap.Platforms(nil)
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
}

func (s *testMapSuite) TestPlatformsValidSelector() {
	platforms := s.platformMap.Platforms(&s.selector2)
	s.Require().Equal(s.platform2, platforms[0])
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
}

func (s *testMapSuite) TestPlatformsInvalidSelector() {
	selector := "invalid"
	platforms := s.platformMap.Platforms(&selector)
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.Require().ElementsMatch(expected, platforms)
}

func (s *testMapSuite) TestPlatformValidSelector() {
	// Get the first platform
	platform, err := s.platformMap.Platform(s.selector1)
	s.Require().NoError(err)
	s.Require().Equal(s.platform1, platform)

	// Get the third platform
	platform, err = s.platformMap.Platform(s.selector3)
	s.Require().NoError(err)
	s.Require().Equal(s.platform3, platform)
}

func (s *testMapSuite) TestPlatformInvalidSelector() {
	// Provide an invalid selector
	selector := "invalid"
	platform, err := s.platformMap.Platform(selector)
	s.Assert().EqualError(err, ErrPlatformNotFound.Error())
	s.Assert().Nil(platform)
}

func (s *testMapSuite) TestSetPlatformExists() {
	// Provide an invalid selector
	selector := "test"
	err := s.platformMap.set(selector, nil)
	s.Require().NoError(err)
	err = s.platformMap.set(selector, nil)
	s.Require().True(errors.Is(err, ErrPlatformExists))
}

func (s *testMapSuite) TestNewMap() {
	// Prepare input
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	yamlLoader := loader.NewYAMLLoader(logger)

	input := &NewInput{
		ConfigPath: "./testdata",
		Logger:     logger,
		Loader:     yamlLoader,
	}
	manager, err := NewMapFromConfig(input)
	s.Require().NoError(err)
	s.Require().NotNil(manager)
	s.Assert().GreaterOrEqual(len(manager.Selectors()), 2)
	s.Assert().Contains(manager.Selectors(), "us-east-1")
	s.Assert().Contains(manager.Selectors(), "us-east-2")
}

func (s *testMapSuite) TestNewMapWithFile() {
	// Prepare input
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	yamlLoader := loader.NewYAMLLoader(logger)

	input := &NewInput{
		ConfigPath: "./testdata/us-east-1.yaml",
		Logger:     logger,
		Loader:     yamlLoader,
	}
	manager, err := NewMapFromConfig(input)
	s.Require().NoError(err)
	s.Require().NotNil(manager)
	s.Assert().GreaterOrEqual(len(manager.Selectors()), 1)
	s.Assert().Equal([]string{"us-east-1"}, manager.Selectors())
}
