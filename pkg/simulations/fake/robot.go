package fake

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type fakeRobot struct {
	name string
	kind string
}

func (f fakeRobot) Name() string {
	return f.name
}

func (f fakeRobot) Kind() string {
	return f.kind
}

func NewRobot(name, kind string) simulations.Robot {
	return &fakeRobot{
		name: name,
		kind: kind,
	}
}
