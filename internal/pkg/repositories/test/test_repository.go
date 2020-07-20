package test

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type testRepository interface {
	create(test Test) (*Test, error)
	getByName(name string) (*Test, error)
	getByValue(value int) ([]Test, error)
}

type testRepositoryImpl struct {
	repository repositories.Repository
}

func (t *testRepositoryImpl) create(test Test) (*Test, error) {
	var tests []domain.Entity
	tests = append(tests, &test)
	_, err := t.repository.Create(tests)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryImpl) getByName(name string) (*Test, error) {
	nameFilter := repositories.NewGormFilter("name", name)
	output := Test{}
	err := t.repository.FindOne(&output, nameFilter)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (t *testRepositoryImpl) getByValue(value int) ([]Test, error) {
	valueFilter := repositories.NewGormFilter("value", value)
	var output []Test
	err := t.repository.Find(&output, nil, nil, valueFilter)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Model returns a pointer to the entity struct for this repository.
func (t *testRepositoryImpl) Model() domain.Entity {
	return &Test{}
}

func NewTestRepository(db *gorm.DB, logger ign.Logger) testRepository {
	return &testRepositoryImpl{
		repository: repositories.NewGormRepository(db, logger, &Test{}),
	}
}

type Test struct {
	gorm.Model
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (t *Test) ParseIn(input domain.Entity) error {
	data, ok := input.(*Test)
	if !ok {
		return errors.New("invalid input data")
	}
	*t = *data
	return nil
}

func (t Test) ParseOut(input interface{}) (domain.Entity, error) {
	data, ok := input.(*Test)
	if !ok {
		return nil, errors.New("invalid input data")
	}
	return data, nil
}

func (Test) TableName() string {
	return "test"
}

func (Test) SingularName() string {
	return "Test"
}

func (Test) PluralName() string {
	return "Tests"
}

func newTest(name string, value int) Test {
	return Test{
		Name:  name,
		Value: value,
	}
}
