package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type Repository struct {
	*simulations.Repository
}

type IRepository interface {
	simulations.IRepository
	CountByOwnerAndCircuit(owner, circuit string) (int, error)
	GetSubT(groupID string) (*SubTSimulation, error)
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	parent := simulations.NewRepository(db, "subt")
	repository := parent.(*simulations.Repository)
	r = &Repository{
		Repository: repository,
	}
	return r
}


func (r *Repository) CountByOwnerAndCircuit(owner, circuit string) (int, error) {

}

func (r *Repository) GetSubT(groupID string) (*SubTSimulation, error) {

}