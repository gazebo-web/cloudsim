package simulations

// Marsupial is a combination of robots.
type Marsupial interface {
	// Parent returns the marsupial parent robot.
	GetParent() Robot
	// Child returns the marsupial child robot.
	GetChild() Robot
}

// marsupial is a Marsupial implementation.
type marsupial struct {
	// parent has a reference to the parent robot.
	parent Robot
	// child has a reference to the child robot.
	child Robot
}

// Returns the parent robot.
func (m marsupial) GetParent() Robot {
	return m.parent
}

// Child returns the child robot.
func (m marsupial) GetChild() Robot {
	return m.child
}

// NewMarsupial initializes a new Marsupial from the given pair of parent and child robots.
func NewMarsupial(parent, child Robot) Marsupial {
	return &marsupial{
		parent: parent,
		child:  child,
	}
}
