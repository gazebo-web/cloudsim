package tracks

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
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
	var db *gorm.DB
	dbConfig, err := ign.NewDatabaseConfigFromEnvVars()
	s.NoError(err)
	db, err = ign.InitDbWithCfg(&dbConfig)
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
		Image:         "dkr.ecr.us-east-1.amazonws.com/stix:latest",
		BridgeImage:   "dkr.ecr.us-east-1.amazonws.com/stix-bridge:latest",
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
