package fake

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// fakeSimulation is a fake simulations.Simulation implementation.
type fakeSimulation struct {
	groupID simulations.GroupID
	status  simulations.Status
	kind    simulations.Kind
	err     *simulations.Error
	image   string
}

func (f *fakeSimulation) HasStatus(status simulations.Status) bool {
	panic("implement me")
}

func (f *fakeSimulation) IsKind(kind simulations.Kind) bool {
	panic("implement me")
}

// Image returns the fake simulation's image.
func (f *fakeSimulation) GetImage() string {
	return f.image
}

// Error returns the fake simulation's error.
// It returns nil if no error has been set.
func (f *fakeSimulation) GetError() *simulations.Error {
	return f.err
}

// GroupID returns the fake simulation's group id.
func (f *fakeSimulation) GetGroupID() simulations.GroupID {
	return f.groupID
}

// Status returns the fake simulation's status.
func (f *fakeSimulation) GetStatus() simulations.Status {
	return f.status
}

// SetStatus sets the fake simulation's status to the given status.
func (f *fakeSimulation) SetStatus(status simulations.Status) {
	f.status = status
}

// Kind returns the simulation's kind.
func (f *fakeSimulation) GetKind() simulations.Kind {
	return f.kind
}

// NewSimulation initializes a new fake simulation.
func NewSimulation(groupID simulations.GroupID, status simulations.Status, kind simulations.Kind,
	err *simulations.Error, image string) simulations.Simulation {
	return &fakeSimulation{
		groupID: groupID,
		status:  status,
		kind:    kind,
		err:     err,
		image:   image,
	}
}
