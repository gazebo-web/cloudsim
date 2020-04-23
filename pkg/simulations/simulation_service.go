package simulations

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strings"
)

// IService
type IService interface {
	GetRepository() IRepository
	SetRepository(repository IRepository)
	Get(groupID string) (*Simulation, error)
	GetAll() []Simulation
	GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(statusFrom, statusTo Status) (*Simulations, error)
	GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error)
	GetParent(groupID string) (*Simulation, error)
	Create(ctx context.Context, simulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Launch(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg)
	Update(groupID string, simulation Simulation) (*Simulation, error)
	UpdateParentFromChildren(parent *Simulation) (*Simulation, error)
	addPermission(resourceID string, permission []per.Action, owner string) *ign.ErrMsg
	addPermissions(resourceID string, permission []per.Action, owner string) *ign.ErrMsg
}

// Service
type Service struct {
	repository IRepository
	userService users.IService
}

// NewService
func NewService(repository IRepository) IService {
	var s IService
	s = &Service{repository: repository}
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

// GetAll
func (s *Service) GetAll() []Simulation {
	panic("Not implemented")

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
func (s * Service) GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error) {
	panic("Not implemented")
}

// GetParent
func (s *Service) GetParent(groupID string) (*Simulation, error) {
	panic("implement me")
}

func (s *Service) Create(ctx context.Context, createSimulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	// TODO: Move this to SubT implementation.
	if createSimulation.Platform == "" {
		createSimulation.Platform = getDefaultPlatformName()
	}
	if createSimulation.Application == "" {
		createSimulation.Application = getDefaultApplicationName()
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
	sim := &Simulation{
		Owner:            &owner,
		Name:             &createSimulation.Name,
		Creator:          &creator,
		Private:          &private,
		StopOnEnd:        &stopOnEnd,
		Platform:         &createSimulation.Platform,
		Application:      &createSimulation.Application,
		Image:            &imageStr,
		GroupID:          &groupID,
		Status: 		  StatusPending.ToIntPtr(),
		// TODO: Move Extra and ExtraSelector to SubT implementation
		// Extra:            createSimulation.Extra,
		// ExtraSelector:    createSimulation.ExtraSelector,
		Robots:           createSimulation.Robots,
		Held:             false,
	}

	// Set the maximum simulation expiration time.
	validFor := s.getMaxDurationForSimulation(ctx, tx, sim)
	validForStr := validFor.String()
	sim.ValidFor = &validForStr

	// TODO: Move to Repository
	if err := tx.Create(sim).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// TODO: Move to SubT implementation
	// Set held state if the user is not a sysadmin and the simulations needs to be held
	if !s.userService.IsSystemAdmin(*user.Username) && s.applications[*sim.Application].simulationIsHeld(ctx, tx, sim) {
		err := sim.UpdateHeldStatus(tx, true)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
		}
	}

	// Set read and write permissions to owner (eg, the team) and to the Application
	// organizing team (eg. subt).
	if em := s.bulkAddPermissions(groupID, []per.Action{per.Read, per.Write}, owner, *sim.Application); em != nil {
		return nil, em
	}

	// Sanity check: check for maximum number of allowed simultaneous simulations per Owner.
	// Also allow Applications to provide custom validations.
	// Dev note: in this case we check 'after' creating the record in the DB to make
	// sure that in case of a race condition then both records are added with pending state
	// and one of those (or both) can be rejected immediately.
	if em := s.checkValidNumberOfSimulations(ctx, tx, sim); em != nil {
		// In case of error we delete the simulation request from DB and exit.
		tx.Model(sim).Update(SimulationDeployment{
			DeploymentStatus: simRejected.ToPtr(),
			ErrorStatus:      simErrorRejected.ToStringPtr(),
		}).Delete(sim)
		return nil, em
	}

	// By default, we launch a single simulation from a createSimulation request.
	// But we also allow specific ApplicationTypes (eg. SubT) to spawn multiple simulations
	// from a single request. When that happens, we call those "child simulations"
	// and they will be grouped by the same parent simulation's groupID.
	simsToLaunch, err := s.prepareSimulations(ctx, tx, sim)
	if err != nil {
		return nil, err
	}

	// Add a 'launch simulation' request to the Launcher Jobs-Pool
	for _, sim := range simsToLaunch {
		groupID := *sim.GroupID
		logger(ctx).Info("StartSimulationAsync about to submit launch task for groupID: " + groupID)
		if err := LaunchSimulation(s, ctx, tx, sim); err != nil {
			logger(ctx).Error(fmt.Sprintf("StartSimulationAsync -- Cannot launch simulation: %s", err.Msg))
		}
	}

	return sim, nil
}

func (s *Service) Launch(ctx context.Context, groupID string, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	panic("implement me")
}

// UpdateParentFromChildren
func (s *Service) UpdateParentFromChildren(parent *Simulation) (*Simulation, error) {
	panic("implement me")
}

// Update
func (s *Service) Update(groupID string, simulation Simulation) (*Simulation, error) {
	sim, err := s.repository.Update(groupID, simulation)
	if err != nil {
		return nil, err
	}
	return sim, nil
}