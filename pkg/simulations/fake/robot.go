package fake

import "github.com/gazebo-web/cloudsim/v4/pkg/simulations"

// fakeRobot is a fake simulations.Robot implementation to use with tests.
type fakeRobot struct {
	// name is the robot's name.
	name string
	// kind is the robot's kind.
	kind string
	// image is the robot's image.
	image string
}

// GetImage returns the fake robot image.
func (f fakeRobot) GetImage() string {
	return f.image
}

// IsEqual checks if the given robot has the same name as the current robot.
func (f fakeRobot) IsEqual(robot simulations.Robot) bool {
	return f.GetName() == robot.GetName()
}

// GetName returns the fake robot's name.
func (f fakeRobot) GetName() string {
	return f.name
}

// GetKind returns the fake robot's kind.
func (f fakeRobot) GetKind() string {
	return f.kind
}

// NewRobot initializes a new fake robot.
func NewRobot(name, kind string) simulations.Robot {
	return &fakeRobot{
		name:  name,
		kind:  kind,
		image: "test.org/image",
	}
}
