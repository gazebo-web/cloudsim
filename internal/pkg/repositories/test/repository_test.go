package test

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
)

type testRepository interface {
	create(tests []*test) ([]*test, error)
	getByName(name string) (*test, error)
	getByValue(value int) ([]test, error)
	getAll() ([]test, error)
	delete(name string) error
	deleteAll() error
	update(name string, data map[string]interface{}) error
	updateAll(data map[string]interface{}) error
}

type testRepositoryImpl struct {
	repository repositories.Repository
}

func (t *testRepositoryImpl) deleteAll() error {
	return t.repository.Delete()
}

func (t *testRepositoryImpl) updateAll(data map[string]interface{}) error {
	return t.repository.Update(data)
}

func (t *testRepositoryImpl) update(name string, data map[string]interface{}) error {
	f := repositories.NewGormFilter("name = ?", name)
	return t.repository.Update(data, f)
}

func (t *testRepositoryImpl) getAll() ([]test, error) {
	var tests []test
	err := t.repository.Find(&tests, nil, nil)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *testRepositoryImpl) create(tests []*test) ([]*test, error) {
	var input []domain.Entity
	for _, test := range tests {
		input = append(input, test)
	}
	_, err := t.repository.Create(input)
	if err != nil {
		return nil, err
	}
	return tests, nil
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
func newTestRepository(base repositories.Repository) testRepository {
	return &testRepositoryImpl{
		repository: base,
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
	return "Test"
}

func (test) PluralName() string {
	return "Tests"
}

func newTest(name string, value int) *test {
	return &test{
		Name:  name,
		Value: value,
	}
}
