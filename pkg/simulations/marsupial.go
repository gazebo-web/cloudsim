package simulations

// Marsupial is a combination of robots.
type Marsupial interface {
	// Parent returns the marsupial parent robot.
	Parent() Robot
	// Child returns the marsupial child robot.
	Child() Robot
}

// marsupial is a Marsupial implementation.
type marsupial struct {
	// parent has a referencere to the parent robot.
	parent Robot
	// child has a reference to the child robot.
	child Robot
}

// Returns the parent robot.
func (m marsupial) Parent() Robot {
	return m.parent
}

// Child returns the child robot.
func (m marsupial) Child() Robot {
	return m.child
}

// NewMarsupial initializes a new Marsupial from the given pair of parent and child robots.
func NewMarsupial(parent, child Robot) Marsupial {
	return &marsupial{
		parent: parent,
		child:  child,
	}
}
