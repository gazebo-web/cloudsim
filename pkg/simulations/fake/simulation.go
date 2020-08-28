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

// Image returns the fake simulation's image.
func (f fakeSimulation) Image() string {
	return f.image
}

// Error returns the fake simulation's error.
// It returns nil if no error has been set.
func (f fakeSimulation) Error() *simulations.Error {
	return f.err
}

// GroupID returns the fake simulation's group id.
func (f fakeSimulation) GroupID() simulations.GroupID {
	return f.groupID
}

// Status returns the fake simulation's status.
func (f fakeSimulation) Status() simulations.Status {
	return f.status
}

// Kind returns the simulation's kind.
func (f fakeSimulation) Kind() simulations.Kind {
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
