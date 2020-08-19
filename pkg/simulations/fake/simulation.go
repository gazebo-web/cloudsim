package fake

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type fakeSimulation struct {
	groupID simulations.GroupID
	status  simulations.Status
	kind    simulations.Kind
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

func NewSimulation(groupID simulations.GroupID, status simulations.Status, kind simulations.Kind) simulations.Simulation {
	return &fakeSimulation{
		groupID: groupID,
		status:  status,
		kind:    kind,
	}
}
