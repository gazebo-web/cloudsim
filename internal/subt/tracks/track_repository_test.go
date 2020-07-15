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

func (s *trackRepositoryTest) TestCreate() {
	s.repository.Create(Track{
		Name:          "Name test",
		Image:         "test.org/image",
		BridgeImage:   "test.org/bridge-image",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 10,
		Public:        true,
	})

	var count int
	err := s.db.Model(&Track{}).Count(&count).Error
	s.NoError(err)
	s.Equal(1, count)
}

func (s *trackRepositoryTest) AfterTest() {
	s.db.Close()
}
