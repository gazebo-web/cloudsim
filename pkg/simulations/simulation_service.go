package simulations

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strings"
	"time"
)

// IService
type IService interface {
	GetRepository() IRepository
	SetRepository(repository IRepository)
	Get(groupID string) (*Simulation, error)
	GetAll(ctx context.Context, input GetAllInput) (*Simulations, *ign.PaginationResult, *ign.ErrMsg)
	GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(statusFrom, statusTo Status) (*Simulations, error)
	GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error)
	GetParent(groupID string) (*Simulation, error)
	Create(ctx context.Context, simulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Launch(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Restart(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Shutdown(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Update(ctx context.Context, groupID string, simulationUpdate SimulationUpdate) (*Simulation, *ign.ErrMsg)
	UpdateParentFromChildren(parent *Simulation) (*Simulation, *ign.ErrMsg)
	Reject(ctx context.Context, simulation *Simulation) (*Simulation, *ign.ErrMsg)
	addPermissionsToOwner(resourceID string, permissions []per.Action, owner string) (bool, *ign.ErrMsg)
	addPermissionsToOwners(resourceID string, permissions []per.Action, owners ...string) *ign.ErrMsg
	Prepare(ctx context.Context, sim *Simulation) (Simulations, *ign.ErrMsg)
}

// Service
type Service struct {
	repository  IRepository
	userService users.IService
	config      ServiceConfig
}

type NewServiceInput struct {
	Repository IRepository
	Config     ServiceConfig
}

type ServiceConfig struct {
	Platform    string
	Application string
	MaxDuration time.Duration
}

// NewService
func NewService(repository IRepository) IService {
	var s IService
	s = &Service{
		repository: repository,
	}
	return s
}

// GetRepository
func (s *Service) GetRepository() IRepository {
	return s.repository
}

// SetRepository
func (s *Service) SetRepository(repository IRepository) {
	s.repository = repository
}

// Get
func (s *Service) Get(groupID string) (*Simulation, error) {
	panic("Not implemented")
}

type GetAllInput struct {
	p               *ign.PaginationRequest
	byStatus        *Status
	invertStatus    bool
	byErrStatus     *ErrorStatus
	invertErrStatus bool
	user            *fuel.User
	includeChildren bool
}

// GetAll
func (s *Service) GetAll(ctx context.Context, input GetAllInput) (*Simulations, *ign.PaginationResult, *ign.ErrMsg) {
	canPerformWithRole, _ := s.userService.CanPerformWithRole(&s.config.Application, *input.user.Username, per.Member)

	includeChildren := false
	if s.userService.IsSystemAdmin(*input.user.Username) {
		includeChildren = true
	}

	sims, pagination, err := s.repository.GetAllPaginated(GetAllPaginatedInput{
		PaginationRequest:          input.p,
		ByStatus:                   input.byStatus,
		InvertStatus:               input.invertStatus,
		ByErrorStatus:              input.byErrStatus,
		InvertErrorStatus:          input.invertErrStatus,
		IncludeChildren:            includeChildren && input.includeChildren,
		CanPerformWithRole:         canPerformWithRole,
		QueryForResourceVisibility: s.userService.QueryForResourceVisibility,
		User:                       input.user,
	})

	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}

	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorInvalidPaginationRequest)
	}

	return sims, pagination, nil
}

// GetAllByOwner
func (s *Service) GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")

}

// GetChildren
func (s *Service) GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

// GetAllParents
func (s *Service) GetAllParents(statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

// GetAllParentsWithErrors
func (s *Service) GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error) {
	panic("Not implemented")
}

// GetParent
func (s *Service) GetParent(groupID string) (*Simulation, error) {
	panic("implement me")
}

func (s *Service) Create(ctx context.Context, createSimulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	if createSimulation.Platform == "" {
		createSimulation.Platform = s.config.Platform
	}
	if createSimulation.Application == "" {
		createSimulation.Application = s.config.Application
	}

	// Set the owner, if missing
	if createSimulation.Owner == "" {
		createSimulation.Owner = *user.Username
	}

	owner := createSimulation.Owner
	if owner == "" {
		owner = *user.Username
	} else {
		// VerifyOwner checks to see if the 'owner' arg is an organization or a user. If the
		// 'owner' is an organization, it verifies that the given 'user' arg has the expected
		// permission in the organization. If the 'owner' is a user, it verifies that the
		// 'user' arg is the same as the owner.
		if ok, em := s.userService.VerifyOwner(owner, *user.Username, per.Read); !ok {
			return nil, em
		}
	}

	private := true
	if createSimulation.Private != nil {
		private = *createSimulation.Private
	}

	stopOnEnd := false
	// Only system admins can request instances to stop on end
	if createSimulation.StopOnEnd != nil && s.userService.IsSystemAdmin(*user.Username) {
		stopOnEnd = *createSimulation.StopOnEnd
	}

	// Create and assign a new GroupID
	groupID := uuid.NewV4().String()

	// Create the SimulationDeployment record in DB. Set initial status.
	creator := *user.Username
	imageStr := strings.Join(createSimulation.Image, ",")
	sim := Simulation{
		Owner:         &owner,
		Name:          &createSimulation.Name,
		Creator:       &creator,
		Private:       &private,
		StopOnEnd:     &stopOnEnd,
		Platform:      &createSimulation.Platform,
		Application:   &createSimulation.Application,
		Image:         &imageStr,
		GroupID:       &groupID,
		Status:        StatusPending.ToIntPtr(),
		Extra:         createSimulation.Extra,
		ExtraSelector: createSimulation.ExtraSelector,
		Robots:        createSimulation.Robots,
		Held:          false,
	}

	// Set the maximum simulation expiration time.
	validFor := s.config.MaxDuration.String()
	sim.ValidFor = &validFor

	createdSim, err := s.repository.Create(&sim)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// Set read and write permissions to owner (eg, the team) and to the Application
	// organizing team (eg. subt).
	if em := s.addPermissionsToOwners(groupID, []per.Action{per.Read, per.Write}, owner, *createdSim.Application); em != nil {
		return nil, em
	}

	return createdSim, nil
}

func (s *Service) Restart(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	panic("implement me")
}

func (s *Service) Launch(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	panic("implement me")
}

func (s *Service) Shutdown(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	panic("Not implemented")
}

// Update
func (s *Service) Update(ctx context.Context, groupID string, simulationUpdate SimulationUpdate) (*Simulation, *ign.ErrMsg) {
	var simulation *Simulation
	var err error

	simulation, err = s.Get(groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	if simulationUpdate.Held != nil {
		simulation.Held = *simulationUpdate.Held
	}

	if simulationUpdate.ErrorStatus != nil {
		simulation.ErrorStatus = simulationUpdate.ErrorStatus
	}

	updatedSim, err := s.repository.Update(groupID, simulation)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	return updatedSim, nil
}

// UpdateParentFromChildren
func (s *Service) UpdateParentFromChildren(parent *Simulation) (*Simulation, *ign.ErrMsg) {
	panic("implement me")
}

func (s *Service) Reject(ctx context.Context, simulation *Simulation) (*Simulation, *ign.ErrMsg) {
	var err error
	if simulation, err = s.repository.Reject(simulation); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	return simulation, nil
}

func (s *Service) addPermissionsToOwner(resourceID string, permissions []per.Action, owner string) (bool, *ign.ErrMsg) {
	var ok bool
	var em *ign.ErrMsg
	for _, p := range permissions {
		ok, em = s.userService.AddResourcePermission(owner, resourceID, p)
		if !ok {
			return ok, em
		}
	}
	return ok, em
}

func (s *Service) addPermissionsToOwners(resourceID string, permissions []per.Action, owners ...string) *ign.ErrMsg {
	for _, owner := range owners {
		ok, err := s.addPermissionsToOwner(resourceID, permissions, owner)
		if !ok {
			return err
		}
	}
	return nil
}

func (s *Service) Prepare(ctx context.Context, sim *Simulation) (Simulations, *ign.ErrMsg) {
	panic("implement me")
}
