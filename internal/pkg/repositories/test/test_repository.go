package test

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type testRepository interface {
	Create(test Test) (*Test, error)
	GetByName(name string) (*Test, error)
	GetByValue(value int) ([]Test, error)
}

type testRepositoryImpl struct {
	repository repositories.GormRepository
}

func (t *testRepositoryImpl) Create(test Test) (*Test, error) {
	var tests []domain.Entity
	tests = append(tests, &test)
	_, err := t.repository.Create(tests)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryImpl) GetByName(name string) (*Test, error) {
	nameFilter := repositories.NewGormFilter("name", name)
	e, err := t.repository.FindOne(nameFilter)
	if err != nil {
		return nil, err
	}
	result := e.(*Test)
	return result, err
}

func (t *testRepositoryImpl) GetByValue(value int) ([]Test, error) {
	valueFilter := repositories.NewGormFilter("value", value)
	output, err := t.repository.Find(nil, nil, valueFilter)
	if err != nil {
		return nil, err
	}
	var result []Test
	for _, o := range output {
		test := o.(*Test)
		result = append(result, *test)
	}
	return result, nil
}

// Model returns a pointer to the entity struct for this repository.
func (t *testRepositoryImpl) Model() domain.Entity {
	return &Test{}
}

func NewTestRepository(db *gorm.DB, logger ign.Logger) testRepository {
	return &testRepositoryImpl{
		repository: repositories.GormRepository{
			DB:     db,
			Logger: logger,
			Entity: &Test{},
		},
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
