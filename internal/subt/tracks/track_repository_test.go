package tracks

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/suite"
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
	db, err := gorm.Open("sqlite3", "/tmp/test.db")
	s.NoError(err)
	s.db = db
	s.db.DropTableIfExists(&Track{})
	s.db.AutoMigrate(&Track{})
	s.repository = NewRepository(db)
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
	track := Track{
		Name:          "Name test",
		Image:         "test.org/image",
		BridgeImage:   "test.org/bridge-image",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 10,
		Public:        true,
	}
	result, err := s.repository.Create(track)
	s.NoError(err)
	s.EqualValues(track, result)

	var count int
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
	s.EqualValues(value, result)
}

func (s *trackRepositoryTest) TestGetAll() {
	valueA := s.addMockData("Practice1")
	valueB := s.addMockData("Practice2")
	valueC := s.addMockData("Practice3")

	result, err := s.repository.GetAll()
	s.NoError(err)
	s.EqualValues(valueA, result[0])
	s.EqualValues(valueB, result[1])
	s.EqualValues(valueC, result[2])
}

func (s *trackRepositoryTest) AfterTest() {
	s.db.Close()
}
