package fake

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// fakeRobot is a fake simulations.Robot implementation to use with tests.
type fakeRobot struct {
	// name is the robot's name.
	name string
	// kind is the robot's kind.
	kind string
}

// Name returns the fake robot's name.
func (f fakeRobot) Name() string {
	return f.name
}

// Kind returns the fake robot's kind.
func (f fakeRobot) Kind() string {
	return f.kind
}

// NewRobot initializes a new fake robot.
func NewRobot(name, kind string) simulations.Robot {
	return &fakeRobot{
		name: name,
		kind: kind,
	}
}
