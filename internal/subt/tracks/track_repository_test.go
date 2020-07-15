package tracks

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestTrackRepository(t *testing.T) {
	suite.Run(t, new(trackRepositoryTest))
}

type trackRepositoryTest struct {
	suite.Suite
	db         *gorm.DB
	repository Repository
}

func (s *trackRepositoryTest) SetupTest() {
	var db *gorm.DB
	dbConfig, err := ign.NewDatabaseConfigFromEnvVars()
	s.NoError(err)
	db, err = ign.InitDbWithCfg(&dbConfig)
	s.db = db
	s.db.DropTableIfExists(&Track{})
	s.db.AutoMigrate(&Track{})
	s.repository = NewRepository(s.db, ign.NewLoggerNoRollbar("track-repository-test", ign.VerbosityDebug))
}

func (s *trackRepositoryTest) addMockData(key string) *Track {
	value := &Track{
		Name:          "Test" + key,
		Image:         "test.org/image-" + key,
		BridgeImage:   "test.org/image-" + key,
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 10,
		Public:        false,
	}
	err := s.db.Model(&Track{}).Save(value).Error
	s.NoError(err)
	return value
}

func (s *trackRepositoryTest) TestCreate() {
	var count int
	err := s.db.Model(&Track{}).Count(&count).Error
	s.NoError(err)
	s.Equal(0, count)

	track := Track{
		Name:          "Name test",
		Image:         "test.org/image",
		BridgeImage:   "test.org/bridge-image",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 10,
		Public:        true,
	}
	_, err = s.repository.Create(track)
	s.NoError(err)

	err = s.db.Model(&Track{}).Count(&count).Error
	s.NoError(err)
	s.Equal(1, count)
}

func (s *trackRepositoryTest) TestGetOne() {
	value := s.addMockData("Practice1")
	s.addMockData("Practice2")
	s.addMockData("Practice3")

	result, err := s.repository.Get("TestPractice1")
	s.NoError(err)
	s.Equal(value.ID, result.ID)
	s.Equal(value.Name, result.Name)
}

func (s *trackRepositoryTest) TestGetAll() {
	valueA := s.addMockData("Practice1")
	valueB := s.addMockData("Practice2")
	valueC := s.addMockData("Practice3")

	result, err := s.repository.GetAll()
	s.NoError(err)
	s.Equal(valueA.ID, result[0].ID)
	s.Equal(valueB.ID, result[1].ID)
	s.Equal(valueC.ID, result[2].ID)
}

func (s *trackRepositoryTest) TestUpdate() {
	value := s.addMockData("Practice1")
	value.BridgeImage = "test.org/bridge-image-changed"
	result, err := s.repository.Update("TestPractice1", *value)
	s.NoError(err)
	s.Equal(value.ID, result.ID)
	s.Equal("test.org/bridge-image-changed", result.BridgeImage)
	s.False(value.UpdatedAt.Equal(result.UpdatedAt))
}

func (s *trackRepositoryTest) TestUpdateNonExistent() {
	_, err := s.repository.Update("TestPractice1", Track{
		Name:          "test",
		Image:         "test",
		BridgeImage:   "test",
		StatsTopic:    "test",
		WarmupTopic:   "test",
		MaxSimSeconds: 20,
		Public:        false,
	})
	s.Error(err)
}

func (s *trackRepositoryTest) TestUpdateEmptyField() {
	value := s.addMockData("Practice1")
	value.BridgeImage = ""
	result, err := s.repository.Update("TestPractice1", *value)
	s.NoError(err)
	s.Empty(result.BridgeImage)
}

func (s *trackRepositoryTest) TestDelete() {
	value := s.addMockData("Practice1")

	var count int
	err := s.db.Model(&Track{}).Count(&count).Error
	s.NoError(err)
	s.Equal(1, count)

	result, err := s.repository.Delete("TestPractice1")
	s.NoError(err)
	s.Equal(value.ID, result.ID)

	err = s.db.Model(&Track{}).Count(&count).Error
	s.NoError(err)
	s.Equal(0, count)
}

func (s *trackRepositoryTest) TestDeleteNonExistent() {
	_, err := s.repository.Delete("TestPractice")
	s.Error(err)
}

func (s *trackRepositoryTest) AfterTest() {
	s.db.Close()
}
