package simulations

import (
	"bitbucket.org/ignitionrobotics/ign-fuelserver/bundles/users"
	"bitbucket.org/ignitionrobotics/ign-go"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"net/http"
)

// TODO: should be moved to package 'testing'

// MockService implements the functions defined for the SimService interface,
// but prevents cloud instances from starting.
type MockService struct {
}

// Start starts this simulation service
func (s *MockService) Start(ctx context.Context) error {
	// nothing to do
	logger(ctx).Debug("Mock Sim Service - Start invoked")
	return nil
}

// Stop stops this simulation service
func (s *MockService) Stop(ctx context.Context) error {
	// nothing to do
	logger(ctx).Debug("Mock Sim Service - Stop invoked")
	return nil
}

// RegisterApplication registers a new application type.
func (s *MockService) RegisterApplication(ctx context.Context, app ApplicationType) {
	// nothing to do
	logger(ctx).Debug(fmt.Sprintf("Mock Sim Service - Registered new Application [%s]", app.getApplicationName()))
}

// StartSimulationAsync is the main func to launch a new simulation
func (s *MockService) StartSimulationAsync(ctx context.Context,
	tx *gorm.DB, createSim *CreateSimulation, user *users.User) (interface{}, *ign.ErrMsg) {
	// Create and assign a new GroupId
	groupId := uuid.NewV4().String()

	private := true
	if createSim.Private != nil {
		private = *createSim.Private
	}

	// Create the SimulationDeployment record in DB. Set initial status.
	simDep := SimulationDeployment{
		Owner:            &createSim.Owner,
		Creator:          user.Username,
		Private:          &private,
		Platform:         &createSim.Platform,
		Application:      &createSim.Application,
		GroupId:          &groupId,
		DeploymentStatus: simPending.ToPtr()}

	if err := tx.Create(&simDep).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	if em := simDep.updateSimDepStatus(tx, simLaunchingNodes); em != nil {
		return nil, em
	}

	if em := simDep.updateSimDepStatus(tx, simLaunchingPods); em != nil {
		return nil, em
	}

	if em := simDep.updateSimDepStatus(tx, simRunning); em != nil {
		return nil, em
	}

	iId := "i-" + uuid.NewV4().String()
	status := "Running"
	// Create a DB record for a  single machine instance
	machine := MachineInstance{
		InstanceId:      &iId,
		LastKnownStatus: &status,
		GroupId:         &groupId,
	}
	if err := tx.Create(&machine).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	logger(ctx).Debug("Mock StartSimulationAsync created fake simulation for groupId: " + groupId)

	return simDep, nil
}

// ShutdownSimulationAsync finishes all resources associated to a cloudsim simulation.
func (s *MockService) ShutdownSimulationAsync(ctx context.Context, tx *gorm.DB,
	groupId string, user *users.User) (interface{}, *ign.ErrMsg) {

	logger(ctx).Debug("Mock ShutdownSimulationAsync requested for groupId: " + groupId)

	dep, err := GetSimulationDeployment(tx, groupId)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}
	if em := dep.updateSimDepStatus(tx, simTerminateRequested); em != nil {
		return nil, em
	}

	if em := dep.updateSimDepStatus(tx, simDeletingPods); em != nil {
		return nil, em
	}

	// Now request to begin the deletion of Nodes and instances
	// are actually removed of if there was an error.
	if _, err := s.DeleteNodesAndHostsForGroup(ctx, tx, dep, user); err != nil {
		return nil, err
	}

	// Update DB record , marking is as terminated
	if err := tx.Model(&dep).Update(SimulationDeployment{
		DeploymentStatus: simTerminated.ToPtr(),
	}).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	logger(ctx).Debug("Mock ShutdownSimulationAsync - successfully removed groupId: " + groupId)

	// Get fresh simulation record from DB and return it
	dep, err = GetSimulationDeployment(tx, groupId)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	return dep, nil
}

// DeleteNodesAndHostsForGroup starts the shutdown of all the kubernates nodes
// and associated Hosts (instances) of a given Cloudsim Group Id.
func (s *MockService) DeleteNodesAndHostsForGroup(ctx context.Context,
	tx *gorm.DB, dep *SimulationDeployment, user *users.User) (interface{}, *ign.ErrMsg) {
	logger(ctx).Debug("Mock delete nodes and hosts for groupId: " + *dep.GroupId)
	return nil, nil
}

// GetCloudMachineInstances returns a paginated list with all cloud instances.
func (s *MockService) GetCloudMachineInstances(ctx context.Context, p *ign.PaginationRequest,
	tx *gorm.DB, byStatus *MachineStatus, invertStatus bool, user *users.User,
	owner *string, application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg) {

	logger(ctx).Debug("Mock machine instances list")

	// Create the DB query
	var machines MachineInstances
	q := tx.Model(&MachineInstance{})
	if byStatus != nil {
		if invertStatus {
			q = q.Where("last_known_status != ?", byStatus.ToStringPtr())
		} else {
			q = q.Where("last_known_status = ?", byStatus.ToStringPtr())
		}
	}

	pagination, err := ign.PaginateQuery(q, &machines, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &machines, pagination, nil
}

// SimulationDeploymentList returns a paginated list with all cloudsim simulations
func (s *MockService) SimulationDeploymentList(ctx context.Context,
	p *ign.PaginationRequest, tx *gorm.DB, byStatus *DeploymentStatus,
	invertStatus bool, byErrStatus *ErrorStatus, invertErrStatus bool,
	user *users.User, owner *string, application *string, includeChildren bool) (*SimulationDeployments, *ign.PaginationResult, *ign.ErrMsg) {

	logger(ctx).Debug("Mock deployment list")

	// Create the DB query
	var sims SimulationDeployments
	q := tx.Model(&SimulationDeployment{})
	if byStatus != nil {
		if invertStatus {
			q = q.Where("deployment_status != ?", byStatus.ToPtr())
		} else {
			q = q.Where("deployment_status = ?", byStatus.ToPtr())
		}
	}
	if byErrStatus != nil {
		if invertErrStatus {
			q = q.Where("error_status != ?", byErrStatus.ToStringPtr())
		} else {
			q = q.Where("error_status = ?", byErrStatus.ToStringPtr())
		}
	}

	pagination, err := ign.PaginateQuery(q, &sims, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &sims, pagination, nil
}

// GetSimulationDeployment returns a single simulation deployment based on its groupId
func (s *MockService) GetSimulationDeployment(ctx context.Context, tx *gorm.DB,
	groupId string, user *users.User) (interface{}, *ign.ErrMsg) {

	dep, err := GetSimulationDeployment(tx, groupId)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}
	return dep, nil
}

// GetSimulationLogsForDownload returns the generated logs from a simulation.
func (s *MockService) GetSimulationLogsForDownload(ctx context.Context, tx *gorm.DB,
	groupId string, user *users.User) (*string, *ign.ErrMsg) {

	return sptr("invalid link"), nil
}

// CustomizeSimRequest allows registered Applications to customize the incoming CreateSimulation request.
// Eg. reading specific SubT fields.
func (s *MockService) CustomizeSimRequest(ctx context.Context, r *http.Request, createSim *CreateSimulation) *ign.ErrMsg {
	return nil
}
