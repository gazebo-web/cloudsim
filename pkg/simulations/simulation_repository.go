package simulations

import "github.com/jinzhu/gorm"

type IRepository interface {
	GetDB() *gorm.DB
	SetDB(db *gorm.DB)
	Get(groupID string) (*[]Simulation, error)
	GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *Repository) Get(groupID string) (*Simulations, error) {
	panic("Not implemented")
}

func (r *Repository) GetAllByOwner(application string, owner string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

func (r *Repository) GetChildren(application string, groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

func (r *Repository) GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

func (r *Repository) Update(groupID string, simulation Simulation) (*Simulation, error) {
	panic("Not implemented")
}