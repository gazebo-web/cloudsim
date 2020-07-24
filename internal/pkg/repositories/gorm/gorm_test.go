package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"reflect"
	"testing"
)

func TestGormRepository(t *testing.T) {
	suite.Run(t, new(testRepositorySuite))
}

type testRepositorySuite struct {
	suite.Suite
	db             *gorm.DB
	baseEntity     domain.Entity
	baseRepository domain.Repository
	repository     testRepository
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
	s.baseEntity = &test{}
	s.baseRepository = NewRepository(s.db, ign.NewLoggerNoRollbar("track-repository-test", ign.VerbosityDebug), s.baseEntity)
	s.repository = newTestRepository(s.baseRepository)
}

func (s testRepositorySuite) AfterTest() {
	s.db.Close()
}

func (s testRepositorySuite) init() {
	test1 := &test{
		Name:  "Test1",
		Value: 1,
	}
	test2 := &test{
		Name:  "Test2",
		Value: 2,
	}

	test3 := &test{
		Name:  "Test3",
		Value: 3,
	}
	_, err := s.repository.create([]*test{test1, test2, test3})
	s.NoError(err, "Should not throw an error when creating test entries")
}

func (s testRepositorySuite) TestCreate() {
	var tests []*test
	tests = append(tests, newTest("test", 1234), newTest("test2", 12345))
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(0, count, "Before creating a test the count should be 0.")

	_, err = s.repository.create(tests)
	s.NoError(err, "Creating a test with the repository should not throw an error.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(2, count, "After creating a test the count should be 2.")
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
	s.NoError(err, "Should not throw an error when getting by value.")
	s.Equal("Test1", result[0].Name, "First database entry should have name Test1.")
	s.Len(result, 1, "The result slice should be length=1.")
}


func (s testRepositorySuite) TestGetAll() {
	s.init()

	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "The total amount of entries should be 3.")

	result, err := s.repository.getAll()
	s.NoError(err, "Should not throw an error when getting all entities.")
	s.Len(result, count, "The result slice should have the same length at the previous count.")
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

func (s testRepositorySuite) TestDeleteAll() {
	s.init()
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "Before removing a test the count should be 3.")

	err = s.repository.deleteAll()
	s.NoError(err, "Should not throw an error when removing all entities.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(0, count, "After removing all tests the count should be 0.")
}

func (s testRepositorySuite) TestDeleteSome() {
	s.init()
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "Before removing a test the count should be 3.")

	err = s.repository.deleteSome([]string{"Test1", "Test2"})
	s.NoError(err, "Should not throw an error when removing some entities.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(1, count, "After removing all tests the count should be 1.")
}

func (s testRepositorySuite) TestDeleteInvalid() {
	s.init()
	var count int
	err := s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "Before removing a test the count should be 3.")

	err = s.repository.delete("1234")
	s.Error(err, "Should throw an error when removing an invalid entity.")

	err = s.db.Model(&test{}).Count(&count).Error
	s.NoError(err, "Should not throw an error when counting.")
	s.Equal(3, count, "After deleting a non existent test failed, the count should remain the same 3.")
}

func (s testRepositorySuite) TestUpdate() {
	s.init()

	err := s.repository.update("Test1", map[string]interface{}{ "name": "Test111", "value": 12345 })
	s.NoError(err, "Should not throw an error when updating an entity.")

	_, err = s.repository.getByName("Test1")
	s.Error(err, "Should throw an error when getting a test with the former name of the updated entity.")

	result, err := s.repository.getByName("Test111")
	s.NoError(err, "Should not throw an error when getting the updated entity.")

	s.Equal("Test111", result.Name)
	s.Equal(12345, result.Value)
}

func (s testRepositorySuite) TestUpdateAll() {
	s.init()

	err := s.repository.updateAll(map[string]interface{}{ "name": "Test123" })
	s.NoError(err, "Should not throw an error when updating all entities.")

	result, err := s.repository.getAll()
	s.NoError(err, "Should not throw an error when getting all entries")

	s.Equal("Test123", result[0].Name)
	s.Equal("Test123", result[1].Name)
	s.Equal("Test123", result[2].Name)
}

func (s testRepositorySuite) TestUpdateZeroValue() {
	s.init()
	err := s.repository.update("Test1", map[string]interface{}{ "value": 0 })
	s.NoError(err, "Should not throw an error when updating an entity.")

	result, err := s.repository.getByName("Test1")
	s.NoError(err, "Should not throw an error when getting the updated entity.")

	s.Equal(0, result.Value)
}

func (s testRepositorySuite) TestUpdateInvalid() {
	s.init()

	err := s.repository.update("12345", nil)
	s.Error(err, "Should throw an error when trying to update an invalid entity.")
}

func (s testRepositorySuite) TestUpdateSomeValues() {
	s.init()

	err := s.repository.updateSome([]string{"Test1", "Test2"}, map[string]interface{}{ "value": 99 })

	result, err := s.repository.getAll()
	s.NoError(err, "Should not throw an error when getting the updated entities.")
	s.Equal(99, result[0].Value)
	s.Equal(99, result[1].Value)
	s.Equal(3, result[2].Value)

}

func (s testRepositorySuite) TestModel() {
	baseEntity := reflect.ValueOf(s.baseEntity)
	baseRepositoryModel := reflect.ValueOf(s.baseRepository.Model())
	s.NotEqual(baseEntity.Pointer(), baseRepositoryModel.Pointer())
}
