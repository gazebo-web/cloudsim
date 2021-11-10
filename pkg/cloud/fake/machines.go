package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
)

var _ machines.Machines = (*Machines)(nil)

// Machines is a fake implementation of the machines.Machines interface.
type Machines struct {
	*mock.Mock
}

// CalculateCosts mocks the CalculateCosts method.
func (m *Machines) CalculateCosts(input []machines.CreateMachinesInput) (machines.CalculateCostsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(machines.CalculateCostsOutput), args.Error(1)
}

// List mocks the List method.
func (m *Machines) List(input machines.ListMachinesInput) (*machines.ListMachinesOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*machines.ListMachinesOutput), args.Error(1)
}

// Create mocks the Create method.
func (m *Machines) Create(input []machines.CreateMachinesInput) ([]machines.CreateMachinesOutput, error) {
	args := m.Called(input)
	return args.Get(0).([]machines.CreateMachinesOutput), args.Error(1)
}

// Terminate mocks the Terminate method.
func (m *Machines) Terminate(input machines.TerminateMachinesInput) error {
	args := m.Called(input)
	return args.Error(0)
}

// Count mocks the Count method.
func (m *Machines) Count(input machines.CountMachinesInput) int {
	args := m.Called(input)
	return args.Int(0)
}

// WaitOK mocks the WaitOK method.
func (m *Machines) WaitOK(input []machines.WaitMachinesOKInput) error {
	args := m.Called(input)
	return args.Error(0)
}

// NewMachines initializes a new machines.Machines fake implementation.
func NewMachines() *Machines {
	return &Machines{
		Mock: new(mock.Mock),
	}
}
