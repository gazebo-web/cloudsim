package simulations

import (
	"github.com/jinzhu/gorm"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Repository represents a set of methods of a Data Access Object for Simulations.
type Repository interface {
	GetDB() *gorm.DB
	SetDB(db *gorm.DB)
	Create(simulation SimulationCreatePersistentInput) (SimulationCreateOutput, error)
	Get(groupID string) (*Simulation, error)
	GetAllPaginated(input GetAllPaginatedInput) (*Simulations, *ign.PaginationResult, error)
	GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error)
	Update(groupID string, simulation *Simulation) (*Simulation, error)
	Reject(simulation *Simulation) (*Simulation, error)
}

// repository is the Repository implementation
type repository struct {
	Application string
	Db *gorm.DB
}

// NewRepository
func NewRepository(db *gorm.DB, application string) Repository {
	var r Repository
	r = &repository{
		Db: db,
		Application: application,
	}
	return r
}

// GetDB
func (r *repository) GetDB() *gorm.DB {
	return r.Db
}

// SetDB
func (r *repository) SetDB(db *gorm.DB) {
	r.Db = db
}

func (r *repository) Create(simulationCreate SimulationCreatePersistentInput) (SimulationCreateOutput, error) {
	simulation := simulationCreate.Input()
	if err := r.Db.Create(simulation).Error; err != nil {
		return nil, err
	}
	return simulation, nil
}

func (r *repository) Update(groupID string, simulation *Simulation) (*Simulation, error) {
	panic("implement me")
}

func (r *repository) Reject(simulation *Simulation) (*Simulation, error) {
	if err := r.Db.Model(simulation).Update(Simulation{
		Status: StatusRejected.ToIntPtr(),
		ErrorStatus:      ErrRejected.ToStringPtr(),
	}).Delete(simulation).Error; err != nil {
		return nil, err
	}
	simulation.Status = StatusRejected.ToIntPtr()
	simulation.ErrorStatus = ErrRejected.ToStringPtr()
	return simulation, nil
}

// Get gets a simulation deployment record by its GroupID
// Fails if not found.
func (r *repository) Get(groupID string) (*Simulation, error) {
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

func (r *repository) GetAllPaginated(input GetAllPaginatedInput) (*Simulations, *ign.PaginationResult, error)  {
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
func (r *repository) GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error) {
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
func (r *repository) GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error) {
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
func (r *repository) GetAllParents(statusFrom, statusTo Status, validErrors []ErrorStatus) (*Simulations, error) {
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