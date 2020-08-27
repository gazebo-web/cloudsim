package fake

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type fakeSimulation struct {
	groupID simulations.GroupID
	status  simulations.Status
	kind    simulations.Kind
	err     *simulations.Error
	image   string
}

func (f fakeSimulation) Image() string {
	return f.image
}

func (f fakeSimulation) Error() *simulations.Error {
	return f.err
}

func (f fakeSimulation) GroupID() simulations.GroupID {
	return f.groupID
}

func (f fakeSimulation) Status() simulations.Status {
	return f.status
}

func (f fakeSimulation) Kind() simulations.Kind {
	return f.kind
}

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
