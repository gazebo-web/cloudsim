package simulations

import "github.com/jinzhu/gorm"

type IRepository interface {
	GetDB() *gorm.DB
	SetDB(db *gorm.DB)
	Get(application, groupID string) (*Simulation, error)
	GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(application string, statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
}

type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{ Db: db }
	return r
}

func (r *Repository) GetDB() *gorm.DB {
	return r.Db
}

func (r *Repository) SetDB(db *gorm.DB) {
	r.Db = db
}

func (r *Repository) Get(application, groupID string) (*Simulation, error) {
	var sim Simulation
	if err := r.Db.Model(&Simulation{}).
		Where("group_id = ? AND application = ?", groupID, application).
		First(&sim).Error; err != nil {
		return nil, err
	}
	return &sim, nil
}

// GetAllByOwner gets a list of simulation deployment records for given application
// filtered by the given owner. The returned set will only contain simulations whose
// Status is between the given statuses range.
func (r *Repository) GetAllByOwner(application string, owner string, statusFrom, statusTo Status) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", application).
		Where("owner = ?", owner).
		Where("status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&sims).Error; err != nil {
			return nil, err
	}
	return &sims, nil
}

// GetChildren returns the child simulation of a given
// GroupID. The returned set will only contain children simulations whose
// deploymentStatus is within the given statuses range, and with NO Error status.
func (r *Repository) GetChildren(application string, groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", application).
		Where("parent_group_id = ?", groupID).
		Where("multi_sim = ?", 2). // TODO: Replace 2 with multiSimChild value.
		Where("error_status IS NULL").
		Where("status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&sims).Error; err != nil {
		return nil, err
	}
	return &sims, nil
}

// GetAllParents returns all the "parent" simulations.
// The returned set will only contain simulations whose status is between the given statuses range,
// and with within the validErrors.
func (r *Repository) GetAllParents(application string, statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", application).
		Where("multi_sim = ?", 1). // TODO: Replace 1 with multiSimParent value.
		Where("(error_status IS NULL OR error_status IN (?))", validErrors).
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&sims).Error; err != nil {
		return nil, err
	}
	return &sims, nil
}

func (r *Repository) Update(groupID string, simulation Simulation) (*Simulation, error) {
	panic("Not implemented")
}