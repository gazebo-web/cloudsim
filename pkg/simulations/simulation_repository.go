package simulations

import "github.com/jinzhu/gorm"

type IRepository interface {
	Get(groupID string) (*[]Simulation, error)
	GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*[]Simulation, error)
	GetChildren(groupID string, application string, statusFrom, statusTo Status) (*[]Simulation, error)
	GetAllParents(application string, statusFrom, statusTo Status) (*[]Simulation, error)
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Get(groupID string) (*[]Simulation, error) {
	panic("Not implemented")
}

func (r *Repository) GetAllByOwner(application string, owner string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")
}

func (r *Repository) GetChildren(application string, groupID string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")
}

func (r *Repository) GetAllParents(application string, statusFrom, statusTo Status) (*[]Simulation, error) {
	panic("Not implemented")
}