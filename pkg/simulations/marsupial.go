package simulations

// Marsupial is a combination of robots.
type Marsupial interface {
	// Parent returns the marsupial parent robot.
	Parent() Robot
	// Child returns the marsupial child robot.
	Child() Robot
}
