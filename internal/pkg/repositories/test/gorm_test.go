package test

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestGormRepository(t *testing.T) {
	suite.Run(t, new(testRepositorySuite))
}

type testRepositorySuite struct {
	suite.Suite
	db         *gorm.DB
	repository testRepository
}

func (s *testRepositorySuite) SetupTest() {
	var db *gorm.DB
	dbConfig, err := ign.NewDatabaseConfigFromEnvVars()
	s.NoError(err)
	db, err = ign.InitDbWithCfg(&dbConfig)
	s.NoError(err)
	s.db = db
	s.db.DropTableIfExists(&Test{})
	s.db.AutoMigrate(&Test{})
	s.repository = NewTestRepository(s.db, ign.NewLoggerNoRollbar("track-repository-Test", ign.VerbosityDebug))
}

func (s testRepositorySuite) AfterTest() {
	s.db.Close()
}

func (s testRepositorySuite) TestCreate() {
	t := newTest("test", 1234)
	var count int
	err := s.db.Model(&Test{}).Count(&count).Error
	s.NoError(err, "Counting should not throw an error.")
	s.Equal(0, count, "Before creating a test the count should be 0.")

	_, err = s.repository.Create(t)
	s.NoError(err, "Creating a test with the repository should not throw an error.")

	err = s.db.Model(&Test{}).Count(&count).Error
	s.NoError(err, "Counting should not throw an error.")
	s.Equal(1, count, "After creating a test the count should be 1.")
}
