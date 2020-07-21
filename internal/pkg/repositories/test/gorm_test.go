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
	s.db.DropTableIfExists(&test{})
	s.db.AutoMigrate(&test{})
	s.repository = newTestRepository(s.db, ign.NewLoggerNoRollbar("track-repository-test", ign.VerbosityDebug))
}

func (s testRepositorySuite) AfterTest() {
	s.db.Close()
}

func (s testRepositorySuite) init() {
	_, err := s.repository.create(test{
		Name:  "Test1",
		Value: 1,
	})
	s.NoError(err, "Should not throw an error when creating Test1")
	_, err = s.repository.create(test{
		Name:  "Test2",
		Value: 2,
	})
	s.NoError(err, "Should not throw an error when creating Test2")
	_, err = s.repository.create(test{
		Name:  "Test3",
		Value: 3,
	})
	s.NoError(err, "Should not throw an error when creating Test3")
}

func (s testRepositorySuite) TestCreate() {
	t := newTest("test", 1234)
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(0, count, "Before creating a test the count should be 0.")

	_, err = s.repository.create(t)
	s.NoError(err, "Creating a test with the repository should not throw an error.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(1, count, "After creating a test the count should be 1.")
}

func (s testRepositorySuite) TestGetByName() {
	s.init()
	result, err := s.repository.getByName("Test1")
	s.NoError(err, "Should not throw an error when getting by name.")
	s.Equal(uint(1), result.ID, "First database entry should be ID=1")
	s.Equal("Test1", result.Name, "Names should match")
}

func (s testRepositorySuite) TestGetByValue() {
	s.init()
	result, err := s.repository.getByValue(1)
	s.NoError(err, "Should not throw an error when getting by name.")
	s.Len(result, 1, "The result slice should be length=1.")
}

func (s testRepositorySuite) TestDelete() {
	s.init()
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "Before removing a test the count should be 3.")

	err = s.repository.delete("Test1")
	s.NoError(err, "Should not throw an error when removing an entity.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(2, count, "After removing a test the count should be 2.")
}
