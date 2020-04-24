package simulations

import (
	"github.com/jinzhu/gorm"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IRepository represents a set of methods of a Data Access Object for Simulations.
type IRepository interface {
	GetDB() *gorm.DB
	SetDB(db *gorm.DB)
	Get(groupID string) (*Simulation, error)
	GetAllPaginated(input GetAllPaginatedInput) (*Simulations, *ign.PaginationResult, error)
	GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
}

// Repository is the IRepository implementation
type Repository struct {
	Application string
	Db *gorm.DB
}

// NewRepository
func NewRepository(db *gorm.DB, application string) IRepository {
	var r IRepository
	r = &Repository{
		Db: db,
		Application: application,
	}
	return r
}

// GetDB
func (r *Repository) GetDB() *gorm.DB {
	return r.Db
}

// SetDB
func (r *Repository) SetDB(db *gorm.DB) {
	r.Db = db
}

func (r *Repository) BeingTX() *gorm.DB {
	return r.Db.Begin()
}

func (r *Repository) CommitTX() *gorm.DB {
	return r.Db.Commit()
}

// Get gets a simulation deployment record by its GroupID
// Fails if not found.
func (r *Repository) Get(groupID string) (*Simulation, error) {
	var sim Simulation
	if err := r.Db.Model(&Simulation{}).
		Where("group_id = ? AND application = ?", groupID, r.Application).
		First(&sim).Error; err != nil {
		return nil, err
	}
	return &sim, nil
}

type GetAllPaginatedInput struct {
	PaginationRequest *ign.PaginationRequest
	ByStatus *Status
	InvertStatus bool
	ByErrorStatus *ErrorStatus
	InvertErrorStatus bool
	IncludeChildren bool
	CanPerformWithRole bool
	QueryForResourceVisibility func(q *gorm.DB, owner *string, user *fuel.User) *gorm.DB
	User *fuel.User
}

func (r *Repository) GetAllPaginated(input GetAllPaginatedInput) (*Simulations, *ign.PaginationResult, error)  {
	var sims Simulations
	q := r.Db.Order("created_at desc, id", true).Where("application = ?", r.Application)

	if !input.IncludeChildren {
		// TODO: Replace 2 with multisimChild value.
		q = q.Where("multi_sim != ?", 2)
	}

	if input.ByStatus != nil {
		query := "status = ?"
		if input.InvertStatus {
			query = "status != ?"
		}
		q = q.Where(query, input.ByStatus.ToInt())
	}

	if input.ByErrorStatus != nil {
		query := "error_status = ?"
		if input.InvertErrorStatus {
			query = "error_status != ?"
		}
		q = q.Where(query, input.ByErrorStatus.ToString())
	}

	if !input.CanPerformWithRole {
		q = input.QueryForResourceVisibility(q, nil, input.User)
	}

	pagination, err := ign.PaginateQuery(q, &sims, *input.PaginationRequest)
	if err != nil {
		return nil, nil, err
	}

	return &sims, pagination, nil
}

// GetAllByOwner gets a list of simulation deployment records for given application
// filtered by the given owner. The returned set will only contain simulations whose
// Status is between the given statuses range.
func (r *Repository) GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", r.Application).
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
func (r *Repository) GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", r.Application).
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
func (r *Repository) GetAllParents(statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error) {
	var sims Simulations
	if err := r.Db.Model(&Simulation{}).
		Where("application = ?", r.Application).
		Where("multi_sim = ?", 1). // TODO: Replace 1 with multiSimParent value.
		Where("(error_status IS NULL OR error_status IN (?))", validErrors).
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&sims).Error; err != nil {
		return nil, err
	}
	return &sims, nil
}

// Update
func (r *Repository) Update(groupID string, simulation Simulation) (*Simulation, error) {
	panic("Not implemented")
}