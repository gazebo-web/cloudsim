package tracks

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	utilsgorm "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestTrackService(t *testing.T) {
	suite.Run(t, new(trackServiceTestSuite))
}

type trackServiceTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository Repository
	service    Service
}

func (s *trackServiceTestSuite) SetupTest() {
	db, err := utilsgorm.GetTestDBFromEnvVars()
	if err != nil {
		s.FailNow("Failed to initialize the database.")
	}
	s.db = db
	s.db.DropTableIfExists(&Track{})
	s.db.AutoMigrate(&Track{})
	logger := ign.NewLoggerNoRollbar("track-service-test", ign.VerbosityDebug)
	s.repository = NewRepository(s.db, logger)
	s.service = NewService(s.repository, validator.New(), logger)
}

func (s *trackServiceTestSuite) TestCreate_OK() {
	input := CreateTrackInput{
		Name:          "Virtual Stix",
		Image:         "https://dkr.ecr.us-east-1.amazonws.com/stix:latest",
		BridgeImage:   "https://dkr.ecr.us-east-1.amazonws.com/stix-bridge:latest",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 3600,
		Public:        true,
		World:         "virtual_stix_headless.ign",
	}
	output, err := s.service.Create(input)
	s.NoError(err)
	s.Equal(input.Name, output.Name)
	s.Equal(input.Image, output.Image)
	s.Equal(input.BridgeImage, output.BridgeImage)
	s.Equal(input.StatsTopic, output.StatsTopic)
	s.Equal(input.MaxSimSeconds, output.MaxSimSeconds)
	s.Equal(input.Public, output.Public)
}

func (s *trackServiceTestSuite) TestCreate_EmptyFields() {
	input := CreateTrackInput{}
	_, err := s.service.Create(input)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestGetAll() {
	trackA, _ := s.repository.Create(Track{
		Name:          "Virtual TestA",
		Image:         "testA",
		BridgeImage:   "testA",
		StatsTopic:    "testA",
		WarmupTopic:   "testA",
		MaxSimSeconds: 30,
		Public:        false,
	})
	trackB, _ := s.repository.Create(Track{
		Name:          "Virtual TestB",
		Image:         "testB",
		BridgeImage:   "testB",
		StatsTopic:    "testB",
		WarmupTopic:   "testB",
		MaxSimSeconds: 30,
		Public:        false,
	})
	trackC, _ := s.repository.Create(Track{
		Name:          "Virtual TestC",
		Image:         "testC",
		BridgeImage:   "testC",
		StatsTopic:    "testC",
		WarmupTopic:   "testC",
		MaxSimSeconds: 30,
		Public:        false,
	})

	tracks, err := s.service.GetAll()
	s.NoError(err)
	s.Len(tracks, 3)
	s.Equal(trackA.Name, tracks[0].Name)
	s.Equal(trackB.Name, tracks[1].Name)
	s.Equal(trackC.Name, tracks[2].Name)
}

func (s *trackServiceTestSuite) TestGetOne_Exists() {
	createdTrack, _ := s.repository.Create(Track{
		Name:          "Virtual TestA",
		Image:         "testA",
		BridgeImage:   "testA",
		StatsTopic:    "testA",
		WarmupTopic:   "testA",
		MaxSimSeconds: 30,
		Public:        false,
	})

	result, err := s.service.Get(createdTrack.Name, 0)

	s.NoError(err)
	s.Equal(createdTrack.Name, result.Name)
}

func (s *trackServiceTestSuite) TestGetOne_NonExistent() {
	_, err := s.service.Get("Test", 0)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestUpdate() {
	_, err := s.repository.Create(Track{
		Name:          "Virtual TestA",
		Image:         "testA",
		BridgeImage:   "testA",
		StatsTopic:    "testA",
		WarmupTopic:   "testA",
		MaxSimSeconds: 30,
		Public:        false,
		World:         "virtual_stix_headless.ign",
	})
	s.Require().NoError(err)

	before, err := s.service.Get("Virtual TestA", 0)
	s.Require().NoError(err)

	updatedTrackInput := UpdateTrackInput{
		Name:          "Virtual TestB",
		Image:         "testB",
		BridgeImage:   "testB",
		StatsTopic:    "testB",
		WarmupTopic:   "testB",
		MaxSimSeconds: 30,
		Public:        true,
		World:         "virtual_testb_headless.ign",
	}

	_, err = s.service.Update("Virtual TestA", updatedTrackInput)
	s.Require().NoError(err)

	result, err := s.service.Get("Virtual TestB", 0)
	s.Require().NoError(err)

	s.Assert().Equal(before.ID, result.ID)
	s.Assert().Equal(updatedTrackInput.Name, result.Name)
	s.Assert().Equal(updatedTrackInput.World, result.World)
}

func (s *trackServiceTestSuite) TestUpdate_InvalidInput() {
	_, err := s.repository.Create(Track{
		Name:          "Virtual TestA",
		Image:         "testA",
		BridgeImage:   "testA",
		StatsTopic:    "testA",
		WarmupTopic:   "testA",
		MaxSimSeconds: 30,
		Public:        false,
	})
	s.NoError(err)

	updateTrackInput := UpdateTrackInput{}
	_, err = s.service.Update("Virtual TestA", updateTrackInput)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestUpdate_NonExistent() {
	updateTrackInput := UpdateTrackInput{}
	_, err := s.service.Update("Virtual TestA", updateTrackInput)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestDelete() {
	_, err := s.repository.Create(Track{
		Name:          "Virtual TestA",
		Image:         "testA",
		BridgeImage:   "testA",
		StatsTopic:    "testA",
		WarmupTopic:   "testA",
		MaxSimSeconds: 30,
		Public:        false,
	})
	s.NoError(err)

	before, err := s.service.Get("Virtual TestA", 0)
	s.NoError(err)

	after, err := s.service.Delete(before.Name)

	s.Equal(before.ID, after.ID)

	result, err := s.service.Get(before.Name, 0)
	s.Error(err)
	s.Nil(result)
}

func (s *trackServiceTestSuite) TestDelete_NonExistent() {
	_, err := s.service.Delete("Test")
	s.Error(err)
}
