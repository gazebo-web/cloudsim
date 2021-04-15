package manager

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/loader"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(testMapSuite))
}

type testMapSuite struct {
	suite.Suite
	selector1   Selector
	selector2   Selector
	selector3   Selector
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
	expected := []Selector{s.selector1, s.selector2, s.selector3}
	s.ElementsMatch(expected, selectors)
}

func (s *testMapSuite) TestPlatforms() {
	platforms := s.platformMap.Platforms()
	expected := []platform.Platform{s.platform1, s.platform2, s.platform3}
	s.ElementsMatch(expected, platforms)
}

func (s *testMapSuite) TestPlatformValidSelector() {
	// Get the first platform
	platform, err := s.platformMap.Platform(s.selector1)
	s.NoError(err)
	s.Equal(s.platform1, platform)

	// Get the third platform
	platform, err = s.platformMap.Platform(s.selector3)
	s.NoError(err)
	s.Equal(s.platform3, platform)
}

func (s *testMapSuite) TestPlatformInvalidSelector() {
	// Provide an invalid selector
	selector := Selector("invalid")
	platform, err := s.platformMap.Platform(selector)
	s.EqualError(err, ErrPlatformNotFound.Error())
	s.Nil(platform)
}

func (s *testMapSuite) TestSet() {
	// Provide an invalid selector
	selector := Selector("test")
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
		ConfigPath: "manager_config_test.yaml",
		Logger:     logger,
		Loader:     yamlLoader,
	}
	manager, err := NewMapFromConfig(input)
	s.Require().NoError(err)
	s.Require().NotNil(manager)
	s.Require().GreaterOrEqual(len(manager.Selectors()), 2)
}
