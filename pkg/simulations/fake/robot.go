package fake

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// fakeRobot is a fake simulations.Robot implementation to use with tests.
type fakeRobot struct {
	// name is the robot's name.
	name string
	// kind is the robot's kind.
	kind string
}

// IsEqual checks if the given robot has the same name as the current robot.
func (f fakeRobot) IsEqual(robot simulations.Robot) bool {
	return f.Name() == robot.Name()
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
