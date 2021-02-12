package platform

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(testMapSuite))
}

type testMapSuite struct {
	suite.Suite
	selector1 Selector
	selector2 Selector
	selector3 Selector
	platform1 Platform
	platform2 Platform
	platform3 Platform
	set       Manager
}

func (s *testMapSuite) SetupSuite() {
	s.selector1 = "test_1"
	s.platform1 = NewPlatform(Components{})

	s.selector2 = "test_2"
	s.platform2 = NewPlatform(Components{})

	s.selector3 = "test_3"
	s.platform3 = NewPlatform(Components{})

	s.set = Map{
		s.selector1: s.platform1,
		s.selector2: s.platform2,
		s.selector3: s.platform3,
	}
}

func (s *testMapSuite) TestSelectors() {
	selectors := s.set.Selectors()
	expected := []Selector{s.selector1, s.selector2, s.selector3}
	s.ElementsMatch(expected, selectors)
}

func (s *testMapSuite) TestPlatforms() {
	platforms := s.set.Platforms()
	expected := []Platform{s.platform1, s.platform2, s.platform3}
	s.ElementsMatch(expected, platforms)
}

func (s *testMapSuite) TestPlatformValidSelector() {
	// Get the first platform
	platform, err := s.set.Platform(s.selector1)
	s.NoError(err)
	s.Equal(s.platform1, platform)

	// Get the third platform
	platform, err = s.set.Platform(s.selector3)
	s.NoError(err)
	s.Equal(s.platform3, platform)
}

func (s *testMapSuite) TestPlatformInvalidSelector() {
	// Provide an invalid selector
	selector := Selector("invalid")
	platform, err := s.set.Platform(selector)
	s.EqualError(err, ErrPlatformNotFound.Error())
	s.Nil(platform)
}
