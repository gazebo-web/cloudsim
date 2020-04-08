package simulations

func Get(groupID string) (*[]Simulation, error) {
	panic("Not implemented")
}

func GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")

}

func GetChildren(groupID string, application string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")
}

func GetAllParents(application string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")
}