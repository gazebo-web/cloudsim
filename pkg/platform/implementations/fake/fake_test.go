package fake

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	fakeStore "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"testing"
)

func TestFakePlatformSuite(t *testing.T) {
	suite.Run(t, &FakePlatformTestSuite{})
}

type FakePlatformTestSuite struct {
	suite.Suite
}

func (s *FakePlatformTestSuite) ValidateFakePlatform(p platform.Platform) {
	s.Require().NotEmpty(p.GetName())
	s.Require().NotNil(p.Machines())
	s.Require().NotNil(p.Storage())
	s.Require().NotNil(p.Orchestrator())
	s.Require().NotNil(p.Store())
	s.Require().NotNil(p.Secrets())
	s.Require().NotNil(p.EmailSender())
}

func (s *FakePlatformTestSuite) TestNewFakePlatformNilConfig() {
	p, err := NewFakePlatform(nil)
	s.Require().NoError(err)
	s.ValidateFakePlatform(p)
}

func (s *FakePlatformTestSuite) TestNewFakePlatformNilComponents() {
	name := "test"
	p, err := NewFakePlatform(&NewInput{
		Name: name,
	})
	s.Require().NoError(err)
	s.Require().Equal(name, p.GetName())
	s.ValidateFakePlatform(p)
}

func (s *FakePlatformTestSuite) TestNewFakePlatformWithComponents() {
	store := fakeStore.NewDefaultFakeStore()
	p, err := NewFakePlatform(&NewInput{
		Components: platform.Components{
			Store: store,
		},
	})
	s.Require().NoError(err)
	s.Require().Exactly(store, p.Store())
	s.ValidateFakePlatform(p)
}
