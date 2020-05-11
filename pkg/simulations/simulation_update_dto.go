package simulations

type SimulationUpdateInput interface {
	Input() SimulationUpdate
}

type SimulationUpdate struct {
	Held        *bool
	ErrorStatus *string
}

func (su SimulationUpdate) Input() SimulationUpdate {
	return su
}

type SimulationUpdateOutput interface {
	Output() Simulation
}
