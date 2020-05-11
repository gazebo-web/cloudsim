package simulations

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type RepositoryMock struct {
	*mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	m := &RepositoryMock{
		Mock: new(mock.Mock),
	}
	return m
}

func (r *RepositoryMock) GetDB() *gorm.DB {
	return nil

}

func (r *RepositoryMock) SetDB(db *gorm.DB) {
	return
}

func (r *RepositoryMock) Create(simulation *Simulation) (*Simulation, error) {
	args := r.Called(simulation)
	result := args.Get(0).(*Simulation)
	return result, args.Error(1)
}

func (r *RepositoryMock) Get(groupID string) (*Simulation, error) {
	args := r.Called(groupID)
	result := args.Get(0).(*Simulation)
	return result, args.Error(1)
}

func (r *RepositoryMock) GetAllPaginated(input GetAllPaginatedInput) (*Simulations, *ign.PaginationResult, error) {
	args := r.Called(input)
	result := args.Get(0).(*Simulations)
	pagination := args.Get(1).(*ign.PaginationResult)
	return result, pagination, args.Error(2)
}

func (r *RepositoryMock) GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error) {
	args := r.Called(owner, statusFrom, statusTo)
	result := args.Get(0).(*Simulations)
	return result, args.Error(1)
}

func (r *RepositoryMock) GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	args := r.Called(groupID, statusFrom, statusTo)
	result := args.Get(0).(*Simulations)
	return result, args.Error(1)
}

func (r *RepositoryMock) GetAllParents(statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error) {
	args := r.Called(statusFrom, statusTo, validErrors)
	result := args.Get(0).(*Simulations)
	return result, args.Error(1)
}

func (r *RepositoryMock) Update(groupID string, simulation *Simulation) (*Simulation, error) {
	args := r.Called(groupID, simulation)
	result := args.Get(0).(*Simulation)
	return result, args.Error(1)
}

func (r *RepositoryMock) Reject(simulation *Simulation) (*Simulation, error) {
	args := r.Called(simulation)
	result := args.Get(0).(*Simulation)
	return result, args.Error(1)
}
