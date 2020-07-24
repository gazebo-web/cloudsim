package tracks

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
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
	repository repositories.Repository
	service    Service
}

func (s *trackServiceTestSuite) SetupTest() {
	var db *gorm.DB
	dbConfig, err := ign.NewDatabaseConfigFromEnvVars()
	s.NoError(err)
	db, err = ign.InitDbWithCfg(&dbConfig)
	s.db = db
	s.db.DropTableIfExists(&Track{})
	s.db.AutoMigrate(&Track{})
	logger := ign.NewLoggerNoRollbar("track-service-test", ign.VerbosityDebug)
	s.repository = repositories.NewGormRepository(s.db, logger, &Track{})
	s.service = NewService(s.repository, validator.New(), logger)
}

func (s *trackServiceTestSuite) init() {
	tracks := []domain.Entity{
		&Track{
			Name:          "Virtual TestA",
			Image:         "testA",
			BridgeImage:   "testA",
			StatsTopic:    "testA",
			WarmupTopic:   "testA",
			MaxSimSeconds: 30,
			Public:        false,
		},
		&Track{
			Name:          "Virtual TestB",
			Image:         "testB",
			BridgeImage:   "testB",
			StatsTopic:    "testB",
			WarmupTopic:   "testB",
			MaxSimSeconds: 30,
			Public:        false,
		},
		&Track{
			Name:          "Virtual TestC",
			Image:         "testC",
			BridgeImage:   "testC",
			StatsTopic:    "testC",
			WarmupTopic:   "testC",
			MaxSimSeconds: 30,
			Public:        false,
		},
	}
	_, err := s.repository.Create(tracks)
	s.NoError(err)
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
	s.init()

	tracks, err := s.service.GetAll(nil, nil)
	s.NoError(err)
	s.Len(tracks, 3)
	s.Equal("Virtual TestA", tracks[0].Name)
	s.Equal("Virtual TestB", tracks[1].Name)
	s.Equal("Virtual TestC", tracks[2].Name)
}

func (s *trackServiceTestSuite) TestGetAllPaginated() {
	s.init()

	page := 0
	size := 2
	tracks, err := s.service.GetAll(&page, &size)

	s.NoError(err)
	s.Len(tracks, 2)
	s.Equal("Virtual TestA", tracks[0].Name)
	s.Equal("Virtual TestB", tracks[1].Name)
}

func (s *trackServiceTestSuite) TestGetAllPaginated_InvalidPage() {
	page := 99
	size := 2
	_, err := s.service.GetAll(&page, &size)
	s.Error(err)
}
func (s *trackServiceTestSuite) TestGetOne_Exists() {
	s.init()

	result, err := s.service.Get("Virtual TestA")

	s.NoError(err)
	s.Equal("Virtual TestA", result.Name)
}

func (s *trackServiceTestSuite) TestGetOne_NonExistent() {
	_, err := s.service.Get("Test")
	s.Error(err)
}

func (s *trackServiceTestSuite) TestUpdate() {
	s.init()
	name := "Virtual TestZ"
	strValue := "testZ"
	numValue := 30
	logicValue := true
	updatedTrackInput := UpdateTrackInput{
		Name:          &name,
		Image:         &strValue,
		BridgeImage:   &strValue,
		StatsTopic:    &strValue,
		WarmupTopic:   &strValue,
		MaxSimSeconds: &numValue,
		Public:        &logicValue,
	}

	s.service.Update("Virtual TestA", updatedTrackInput)

	result, err := s.service.Get("Virtual TestZ")
	s.NoError(err)
	s.Equal(*updatedTrackInput.Name, result.Name)
}

func (s *trackServiceTestSuite) TestUpdate_InvalidInput() {
	updateTrackInput := UpdateTrackInput{}
	_, err := s.service.Update("Virtual TestA", updateTrackInput)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestUpdate_NonExistent() {
	updateTrackInput := UpdateTrackInput{}
	_, err := s.service.Update("Virtual TestA", updateTrackInput)
	s.Error(err)
}

func (s *trackServiceTestSuite) TestDelete() {
	s.init()

	before, err := s.service.Get("Virtual TestA")
	s.NoError(err)

	after, err := s.service.Delete(before.Name)

	s.Equal(before.ID, after.ID)

	result, err := s.service.Get(before.Name)
	s.Error(err)
	s.Nil(result)
}

func (s *trackServiceTestSuite) TestDelete_NonExistent() {
	_, err := s.service.Delete("Test")
	s.Error(err)
}
