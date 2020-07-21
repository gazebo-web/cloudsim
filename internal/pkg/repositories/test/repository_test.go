package test

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type testRepository interface {
	create(test test) (*test, error)
	getByName(name string) (*test, error)
	getByValue(value int) ([]test, error)
	getAll() ([]test, error)
	delete(name string) error
	update(name string, data test) error
}

type testRepositoryImpl struct {
	repository repositories.Repository
}

func (t *testRepositoryImpl) update(name string, data test) error {
	f := repositories.NewGormFilter("name = ?", name)
	return t.repository.Update(&data, f)
}

func (t *testRepositoryImpl) getAll() ([]test, error) {
	var tests []test
	err := t.repository.Find(&tests, nil, nil)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *testRepositoryImpl) create(test test) (*test, error) {
	var tests []domain.Entity
	tests = append(tests, &test)
	_, err := t.repository.Create(tests)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryImpl) getByName(name string) (*test, error) {
	f := repositories.NewGormFilter("name = ?", name)
	output := test{}
	err := t.repository.FindOne(&output, f)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (t *testRepositoryImpl) getByValue(value int) ([]test, error) {
	f := repositories.NewGormFilter("value = ?", value)
	var output []test
	err := t.repository.Find(&output, nil, nil, f)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (t *testRepositoryImpl) delete(name string) error {
	f := repositories.NewGormFilter("name = ?", name)
	return t.repository.Delete(f)
}

// Model returns a pointer to the entity struct for this repository.
func (t *testRepositoryImpl) Model() domain.Entity {
	return &test{}
}

// newTestRepository initializes a new testRepository.
func newTestRepository(db *gorm.DB, logger ign.Logger) testRepository {
	return &testRepositoryImpl{
		repository: repositories.NewGormRepository(db, logger, &test{}),
	}
}

type test struct {
	gorm.Model
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (test) TableName() string {
	return "test"
}

func (test) SingularName() string {
	return "test"
}

func (test) PluralName() string {
	return "Tests"
}

func newTest(name string, value int) test {
	return test{
		Name:  name,
		Value: value,
	}
}
