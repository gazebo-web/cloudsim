package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
)

var _ cloud.Machines = (*Machines)(nil)

// Machines is a fake implementation of the cloud.Machines interface.
type Machines struct {
	*mock.Mock
}

// List mocks the List method.
func (m *Machines) List(input cloud.ListMachinesInput) (*cloud.ListMachinesOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*cloud.ListMachinesOutput), args.Error(1)
}

// Create mocks the Create method.
func (m *Machines) Create(input []cloud.CreateMachinesInput) ([]cloud.CreateMachinesOutput, error) {
	args := m.Called(input)
	return args.Get(0).([]cloud.CreateMachinesOutput), args.Error(1)
}

// Terminate mocks the Terminate method.
func (m *Machines) Terminate(input cloud.TerminateMachinesInput) error {
	args := m.Called(input)
	return args.Error(0)
}

// Count mocks the Count method.
func (m *Machines) Count(input cloud.CountMachinesInput) int {
	args := m.Called(input)
	return args.Int(0)
}

// WaitOK mocks the WaitOK method.
func (m *Machines) WaitOK(input []cloud.WaitMachinesOKInput) error {
	args := m.Called(input)
	return args.Error(0)
}

// NewMachines initializes a new cloud.Machines fake implementation.
func NewMachines() *Machines {
	return &Machines{
		Mock: new(mock.Mock),
	}
}
