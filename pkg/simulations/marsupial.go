package simulations

// Marsupial is a combination of robots.
type Marsupial interface {
	// Parent returns the marsupial parent robot.
	Parent() Robot
	// Child returns the marsupial child robot.
	Child() Robot
}

type marsupial struct {
	parent Robot
	child  Robot
}

func (m marsupial) Parent() Robot {
	return m.parent
}

func (m marsupial) Child() Robot {
	return m.child
}

func NewMarsupial(parent, child Robot) Marsupial {
	return &marsupial{
		parent: parent,
		child:  child,
	}
}
