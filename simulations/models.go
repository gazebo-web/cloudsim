package simulations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	subtsim "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strconv"
	"strings"
	"time"
)

// Force SimulationDeployment to implement simulations.Simulation interface.
var (
	_ subtsim.Simulation = (*SimulationDeployment)(nil)
)

// SimulationDeployment represents a cloudsim simulation .
type SimulationDeployment struct {
	// Override default GORM Model fields
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Added 2 milliseconds to DeletedAt field
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`
	// Timestamp in which this simulation was stopped/terminated.
	StoppedAt *time.Time `gorm:"type:timestamp(3) NULL" json:"stopped_at,omitempty"`
	// Represents the maximum time this simulation should live. After that time
	// it will be eligible for automatic termination.
	// It is a time.Duration (stored as its string representation).
	ValidFor *string `json:"valid_for,omitempty"`
	// The owner of this deployment (must exist in UniqueOwners). Can be user or org.
	// Also added to the name_owner unique index
	Owner *string `json:"owner,omitempty"`
	// The username of the User that created this resource (usually got from the JWT)
	Creator *string `json:"creator,omitempty"`
	// Private - True to make this a private resource
	Private *bool `json:"private,omitempty"`
	// When shutting down simulations, stop EC2 instances instead of terminating them. Requires admin privileges.
	StopOnEnd *bool `json:"stop_on_end"`
	// The user defined Name for the simulation.
	Name *string `json:"name,omitempty"`
	// The docker image url to use for the simulation (usually for the Field Computer)
	Image *string `json:"image,omitempty" form:"image"`
	// GroupID - Simulation Unique identifier
	// All k8 pods and services (or other created resources) will share this groupID
	GroupID *string `gorm:"not null;unique" json:"group_id"`
	// ParentGroupID (optional) holds the GroupID of the parent simulation record.
	// It is used with requests for multi simulations (multiSims), where a single
	// user request spawns multiple simulation runs based on a single template.
	ParentGroupID *string `json:"parent"`
	// MultiSim holds which role this simulation plays within a multiSim deployment.
	// Values should be of type MultiSimType.
	MultiSim int
	// A value from DeploymentStatus constants
	DeploymentStatus *int `json:"status,omitempty"`
	// A value from ErrorStatus constants
	ErrorStatus *string `json:"error_status,omitempty"`
	// NOTE: statuses should be updated in sequential DB Transactions. ie. one status per TX.
	Platform    *string `json:"platform,omitempty" form:"platform"`
	Application *string `json:"application,omitempty" form:"application"`
	// TODO: both fields Extra and ExtraSelector should be a separate table, specific to
	//   each Application, and with a reference back to this simulation ID.
	// A free form string field to store extra details, usually associated to the
	// chosen Application. Eg. SubT would store here the different robot names, types
	// and images.
	Extra *string `gorm:"size:999999" json:"extra,omitempty"`
	// A extra string field to store a selector that can help specific Applications
	// to filter simulations (eg. SQL WHERE). SubT could store the circuit here.
	ExtraSelector *string `json:"-"`
	// TODO: This is a field specific to SubT. As such this is a temporary field
	//  that should be included in the same separate table where Extra and
	//  ExtraSelector should reside.
	// Contains the names of all robots in the simulation in a comma-separated list.
	Robots *string `gorm:"size:1000" json:"robots"`
	// TODO: This is a field specific to SubT. This is a temporary field that should be
	//  extracted from the SimulationDeployment struct.
	Held bool `json:"held"`
	// Processed indicates that this simulation has been post-processed.
	// Used to avoid post-processing simulations more than once.
	Processed bool `json:"-"`
	// AuthorizationToken contains a security token used to let external services authorize requests related to this
	// simulation.
	// This token is currently used to establish connections with the simulation's websocket server.
	AuthorizationToken *string `json:"-"`
	// Score has the simulation's score. It's updated when the simulations finishes and gets processed.
	Score *float64 `json:"score,omitempty"`
}

// GetRunIndex returns the simulation's run index.
func (dep *SimulationDeployment) GetRunIndex() int {
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return 0
	}

	if extra.RunIndex == nil {
		return 0
	}

	return *extra.RunIndex
}

// GetWorldIndex returns the simulation's world index.
func (dep *SimulationDeployment) GetWorldIndex() int {
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return 0
	}

	if extra.WorldIndex == nil {
		return 0
	}

	return *extra.WorldIndex
}

// IsProcessed returns true if the SimulationDeployment has been processed.
func (dep *SimulationDeployment) IsProcessed() bool {
	return dep.Processed
}

// GetOwner returns the SimulationDeployment's Owner.
func (dep *SimulationDeployment) GetOwner() *string {
	return dep.Owner
}

// GetCreator returns the SimulationDeployment's Creator. It returns an empty string if no creator has been assigned.
func (dep *SimulationDeployment) GetCreator() string {
	if dep.Creator == nil {
		return ""
	}
	return *dep.Creator
}

// GetPlatform returns the SimulationDeployment's Platform.
func (dep *SimulationDeployment) GetPlatform() *string {
	return dep.Platform
}

// GetName returns the SimulationsDeployment's Name. t returns an empty string if no name has been assigned.
func (dep *SimulationDeployment) GetName() string {
	if dep.Name == nil {
		return ""
	}
	return *dep.Name
}

// GetGroupID returns the SimulationDeployment's GroupID.
func (dep *SimulationDeployment) GetGroupID() simulations.GroupID {
	return simulations.GroupID(*dep.GroupID)
}

// GetStatus returns the SimulationDeployment's DeploymentStatus.
func (dep *SimulationDeployment) GetStatus() simulations.Status {
	switch *dep.DeploymentStatus {
	case simPending.ToInt():
		return simulations.StatusPending
	case simLaunchingNodes.ToInt():
		return simulations.StatusLaunchingInstances
	case simLaunchingPods.ToInt():
		return simulations.StatusLaunchingPods
	case simParentLaunching.ToInt():
		return simulations.StatusUnknown
	case simParentLaunchingWithErrors.ToInt():
		return simulations.StatusUnknown
	case simRunning.ToInt():
		return simulations.StatusRunning
	case simTerminateRequested.ToInt():
		return simulations.StatusTerminateRequested
	case simDeletingPods.ToInt():
		return simulations.StatusDeletingPods
	case simDeletingNodes.ToInt():
		return simulations.StatusDeletingNodes
	case simTerminatingInstances.ToInt():
		return simulations.StatusTerminatingInstances
	case simTerminated.ToInt():
		return simulations.StatusTerminated
	case simRejected.ToInt():
		return simulations.StatusRejected
	case simSuperseded.ToInt():
		return simulations.StatusSuperseded
	default:
		return simulations.StatusUnknown
	}
}

// HasStatus checks that the SimulationDeployment's DeploymentStatus is equal to the given status.
func (dep *SimulationDeployment) HasStatus(status simulations.Status) bool {
	return dep.GetStatus() == status
}

// SetStatus sets the SimulationDeployment's DeploymentStatus to the given status.
func (dep *SimulationDeployment) SetStatus(status simulations.Status) {
	dep.setStatus(status)
}

// GetKind returns the SimulationDeployment's Kind. It parses the MultiSim field into a Kind.
func (dep *SimulationDeployment) GetKind() simulations.Kind {
	return simulations.Kind(dep.MultiSim)
}

// IsKind checks that the SimulationDeployment
func (dep *SimulationDeployment) IsKind(kind simulations.Kind) bool {
	return dep.GetKind() == kind
}

// GetError returns the SimulationDeployment's ErrorStatus
func (dep *SimulationDeployment) GetError() *simulations.Error {
	if dep.ErrorStatus == nil {
		return nil
	}
	err := simulations.Error(*dep.ErrorStatus)
	return &err
}

// GetImage returns the SimulationDeployment's image.
func (dep *SimulationDeployment) GetImage() string {
	if dep.Image == nil {
		return ""
	}
	return *dep.Image
}

// GetValidFor returns the SimulationDeployment's ValidFor parsed as time.Duration.
func (dep *SimulationDeployment) GetValidFor() time.Duration {
	if dep.ValidFor == nil {
		return 0
	}
	d, err := time.ParseDuration(*dep.ValidFor)
	if err != nil {
		return 0
	}
	return d
}

// GetTrack returns the SimulationDeployment's circuit name.
func (dep *SimulationDeployment) GetTrack() string {
	info, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return ""
	}
	return info.Circuit
}

// GetToken returns the SimulationDeployment's websocket authorization token.
func (dep *SimulationDeployment) GetToken() *string {
	return dep.AuthorizationToken
}

// GetRobots parses the robots from the extra field and returns them as a slice of robots.
func (dep *SimulationDeployment) GetRobots() []simulations.Robot {
	info, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return nil
	}
	result := make([]simulations.Robot, len(info.Robots))
	for i, robot := range info.Robots {
		r := new(SubTRobot)
		*r = robot
		result[i] = r
	}
	return result
}

// GetMarsupials parses the extra field and returns the marsupials.
func (dep *SimulationDeployment) GetMarsupials() []simulations.Marsupial {
	info, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return nil
	}
	result := make([]simulations.Marsupial, len(info.Marsupials))
	for i, marsupial := range info.Marsupials {
		m := new(SubTMarsupial)
		*m = marsupial
		result[i] = m
	}
	return result
}

// NewSimulationDeployment creates and initializes a simulation deployment struct.
// TODO: Receive a DTO struct as a parameter
func NewSimulationDeployment() (*SimulationDeployment, error) {
	// Generate an auth token
	authToken, err := generateToken(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authorization token: %s", err.Error())
	}

	// Initialize the deployment
	dep := &SimulationDeployment{
		AuthorizationToken: &authToken,
	}

	return dep, nil
}

// GetSimulationDeployment gets a simulation deployment record by its GroupID
// Fails if not found.
func GetSimulationDeployment(tx *gorm.DB, groupID string) (*SimulationDeployment, error) {
	var dep SimulationDeployment
	if err := tx.Model(&SimulationDeployment{}).Where("group_id = ?", groupID).First(&dep).Error; err != nil {
		return nil, err
	}
	return &dep, nil
}

// GetSimulationDeploymentsByOwner gets a list of simulation deployment records
// filtered by a given owner. The returned set will only contain simulations whose
// deploymentStatus is within the given statuses range.
func GetSimulationDeploymentsByOwner(tx *gorm.DB, owner string, statusFrom,
	statusTo DeploymentStatus) (*SimulationDeployments, error) {

	var deps SimulationDeployments
	if err := tx.Model(&SimulationDeployment{}).Where("owner = ?", owner).
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).Find(&deps).Error; err != nil {
		return nil, err
	}
	return &deps, nil
}

// GetSimulationDeploymentsByCircuit gets a list of simulation deployments for a given circuit.
// The returned set will only contain simulations whose deploymentStatus is within the given statuses range.
func GetSimulationDeploymentsByCircuit(tx *gorm.DB, circuit string, statusFrom,
	statusTo DeploymentStatus, held *bool) (*SimulationDeployments, error) {

	var deps SimulationDeployments
	q := tx.Model(&SimulationDeployment{}).
		Where("extra_selector = ?", circuit).
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Where("multi_sim = ? OR multi_sim = ?", multiSimSingle, multiSimParent)

	if held != nil {
		q = q.Where("held = ?", held)
	}

	if err := q.Find(&deps).Error; err != nil {
		return nil, err
	}

	return &deps, nil
}

// GetChildSimulationDeployments returns the child simulation deployments of a given
// parent simulation deployment. The returned set will only contain children simulations whose
// deploymentStatus is within the given statuses range, and with NO ErrorStatus.
func GetChildSimulationDeployments(tx *gorm.DB, dep *SimulationDeployment, statusFrom,
	statusTo DeploymentStatus) (*SimulationDeployments, error) {

	var deps SimulationDeployments
	if err := tx.Model(&SimulationDeployment{}).
		Where("parent_group_id = ?", *dep.GroupID).
		Where("multi_sim = ?", multiSimChild).
		Where("error_status IS NULL").
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&deps).Error; err != nil {
		return nil, err
	}
	return &deps, nil
}

// GetParentSimulationDeployments returns the "parent" simulation deployments from
// Multi Sim deployments.  The returned set will only contain simulations whose
// deploymentStatus is within the given statuses range, and with NO ErrorStatus.
func GetParentSimulationDeployments(tx *gorm.DB, statusFrom,
	statusTo DeploymentStatus, validErrors []ErrorStatus) (*SimulationDeployments, error) {

	var deps SimulationDeployments
	if err := tx.Model(&SimulationDeployment{}).
		Where("multi_sim = ?", multiSimParent).
		Where("(error_status IS NULL OR error_status IN (?))", validErrors).
		Where("deployment_status BETWEEN ? AND ?", int(statusFrom), int(statusTo)).
		Find(&deps).Error; err != nil {
		return nil, err
	}
	return &deps, nil
}

// GetParentSimulation returns the DB record corresponding to the parent simulation
func GetParentSimulation(tx *gorm.DB, dep *SimulationDeployment) (*SimulationDeployment, error) {
	var parent SimulationDeployment
	if err := tx.Model(&SimulationDeployment{}).
		Where("group_id = ?", *dep.ParentGroupID).
		Find(&parent).Error; err != nil {
		return nil, err
	}
	return &parent, nil
}

// toJSON marshals a SimulationD ing.
func (dep *SimulationDeployment) toJSON() (*string, error) {
	byt, err := json.Marshal(*dep)
	if err != nil {
		return nil, err
	}
	return sptr(string(byt)), nil
}

// Clone clones a SimulationDeployment
func (dep *SimulationDeployment) Clone() *SimulationDeployment {
	clone := *dep

	// Clear default GORM Model fields
	clone.ID = uint(0)
	clone.CreatedAt = time.Time{}
	clone.UpdatedAt = time.Time{}
	clone.StoppedAt = nil
	clone.DeletedAt = nil

	return &clone
}

// UpdateHeldStatus returns an error if the SimulationDeployment held field failed to update.
func (dep *SimulationDeployment) UpdateHeldStatus(tx *gorm.DB, state bool) error {
	dep.Held = state
	if err := tx.Save(&dep).Error; err != nil {
		return err
	}
	return nil
}

// UpdateProcessed sets the given state in the Processed value.
// Returns an error if the SimulationDeployment Processed field failed to update.
func (dep *SimulationDeployment) UpdateProcessed(tx *gorm.DB, state bool) error {
	dep.Processed = state
	if err := tx.Save(&dep).Error; err != nil {
		return err
	}
	return nil
}

// UpdateScore is used to update the score of the current simulation while it's being processed.
// Returns an error if the SimulationDeployment Score field failed to update.
func (dep *SimulationDeployment) UpdateScore(tx *gorm.DB, score *float64) *ign.ErrMsg {
	dep.Score = score
	if err := tx.Model(&SimulationDeployment{}).Where("id = ?", dep.ID).Update("score", score).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	return nil
}

// IsRunning returns true if the SimulationDeployment can be considered "running".
// Running goes from the moment a simulation is scheduled to run (ie. Pending state) up to until
// its termination is requested (ie. 'terminate requested' state).
func (dep *SimulationDeployment) IsRunning() bool {
	// If its status is related to the termination phase or greater (after terminate status), then it
	// is not running.
	if *dep.DeploymentStatus >= int(simTerminateRequested) {
		return false
	}
	// If the simulation status is held, then the simulation is not running.
	if dep.Held {
		return false
	}
	// If it has an error, then it cannot be running
	if dep.ErrorStatus != nil {
		return false
	}
	// otherwise we assume it's running
	return true
}

func (dep *SimulationDeployment) isMultiSim() bool {
	return dep.isMultiSimParent() || dep.isMultiSimChild()
}

func (dep *SimulationDeployment) isMultiSimParent() bool {
	return dep.MultiSim == int(multiSimParent)
}

func (dep *SimulationDeployment) isMultiSimChild() bool {
	return dep.MultiSim == int(multiSimChild)
}

// updateCompoundStatuses updates the Parent's DeploymentStatus and ErrorStatus fields
// based on the status of its children. It is assumed that the receiver
// is a Parent in a Multi Simulation.
func (dep *SimulationDeployment) updateCompoundStatuses(tx *gorm.DB) *ign.ErrMsg {

	if !dep.isMultiSimParent() {
		return nil
	}

	// Get this parent's children simulations (all of them)
	var children SimulationDeployments
	if err := tx.Model(&SimulationDeployment{}).Where("parent_group_id = ?", *dep.GroupID).
		Where("multi_sim = ?", multiSimChild).Find(&children).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	if len(children) == 0 {
		// No children yet. Leave the parent status as is
		return nil
	}

	var errorStatus *ErrorStatus
	status := int(simRejected)
	statusFromChildWithError := int(simRejected)

	// A flag used to know if at least one children has ErrorStatus
	childWithError := false
	// A flag used to know if at least one children is running
	childRunning := false

	// Iterate over children to 'compute' the multi-sim status.
	// dev note: this for-loop will find the most severe Error status, in case of any
	// errors. Otherwise it will result in the status of the newest (youngest) child
	// simulation.
	for _, child := range children {
		if child.IsRunning() {
			// the child is actually running (and not waiting in the pending queue)
			childRunning = true
		}

		var childErrorStatus *ErrorStatus
		if child.ErrorStatus != nil {
			childWithError = true
			childErrorStatus = ErrorStatusFrom(*child.ErrorStatus)
		}
		if childErrorStatus == nil {
			// If there is no error status yet, then just update the deployment status
			status = Min(status, *child.DeploymentStatus)
		} else if isMoreSevere(errorStatus, childErrorStatus) {
			// Simulations with Error statuses have priority over others
			errorStatus = childErrorStatus
			statusFromChildWithError = *child.DeploymentStatus
		}
	}

	if childRunning {
		// If there are children running then we don't want to mark the parent with error.
		errorStatus = nil
		if !childWithError {
			status = Max(status, int(simParentLaunching))
		} else {
			status = int(simParentLaunchingWithErrors)
		}
	} else if childWithError {
		// If no children is running and there were errors, mark the parent with the
		// most severe error (and its associated status)
		status = statusFromChildWithError
	}

	// Update DB record with computed statuses
	var errorStrPtr *string
	if errorStatus != nil {
		errorStrPtr = errorStatus.ToStringPtr()
	}
	if err := tx.Model(&dep).Update(SimulationDeployment{
		DeploymentStatus: &status,
		ErrorStatus:      errorStrPtr,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.DeploymentStatus = &status
	if errorStatus != nil {
		dep.ErrorStatus = errorStrPtr
	}

	return nil
}

// recordStop marks the simulation as stopped and also updates the DB.
func (dep *SimulationDeployment) recordStop(tx *gorm.DB) *ign.ErrMsg {
	val := time.Now()
	if err := tx.Model(&dep).Update(SimulationDeployment{
		StoppedAt: &val,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.StoppedAt = &val
	return nil
}

// updateValidFor updates this SimulationDeployment's ValidFor field in the database.
func (dep *SimulationDeployment) updateValidFor(tx *gorm.DB, validFor time.Duration) *ign.ErrMsg {
	validForStr := validFor.String()
	if err := tx.Model(&dep).Update(SimulationDeployment{
		ValidFor: &validForStr,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.ValidFor = &validForStr
	return nil
}

// updateSimDepStatus updates this SimulationDeployment's DeploymentStatus in the database.
func (dep *SimulationDeployment) updateSimDepStatus(tx *gorm.DB, st DeploymentStatus) *ign.ErrMsg {
	val := st.ToPtr()
	if err := tx.Model(&dep).Update(SimulationDeployment{
		DeploymentStatus: val,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.DeploymentStatus = val
	return nil
}

// updatePlatform updates this SimulationDeployment's Platform in the database.
func (dep *SimulationDeployment) updatePlatform(tx *gorm.DB, platform string) *ign.ErrMsg {
	if err := tx.Model(&dep).Update(SimulationDeployment{
		Platform: &platform,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.Platform = &platform

	return nil
}

// setErrorStatus sets an error status to the simulation deployment and updates the DB.
func (dep *SimulationDeployment) setErrorStatus(tx *gorm.DB, st ErrorStatus) *ign.ErrMsg {
	val := st.ToStringPtr()
	if err := tx.Model(&dep).Update(SimulationDeployment{
		ErrorStatus: val,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.ErrorStatus = val
	return nil
}

// assertSimDepStatus checks that the SimulationStatus is in a given valid status. If
// it has a different status OR if it is in an error status, then this function
// will fail. It is an 'assert'
func (dep *SimulationDeployment) assertSimDepStatus(st DeploymentStatus) *ign.ErrMsg {
	if dep.ErrorStatus != nil {
		return ign.NewErrorMessageWithArgs(ign.ErrorInvalidSimulationStatus, nil, []string{"error status", *dep.ErrorStatus})
	}
	if *dep.DeploymentStatus != int(st) {
		ds := DeploymentStatus(*dep.DeploymentStatus)
		return ign.NewErrorMessageWithArgs(ign.ErrorInvalidSimulationStatus, nil, []string{ds.String()})
	}
	return nil
}

// MultiSimType represents which rol plays a simulation within a multiSim deployment
type MultiSimType int

const (
	// multiSimSingle represents a "single" simulation. Meaning it didn't spawn any other simulations.
	// This is the default.
	multiSimSingle MultiSimType = iota
	// multiSimParent is used to tag the main simulation request, which was used to spawn the actual simulations.
	multiSimParent
	// multiSimChild is used to tag those SimulationDeployment records that represent actual simulation runs
	// spawned as part of a multiSim launch.
	multiSimChild
)

// MarkAsMultiSimParent marks this simulationDeployment as the parent in a multiSimulation.
// Parent simulationDeployment usually don't launch simulations themselves but instead
// group a set of child simulations.
func (dep *SimulationDeployment) MarkAsMultiSimParent(tx *gorm.DB) *ign.ErrMsg {
	if err := tx.Model(&dep).Update(SimulationDeployment{
		MultiSim: int(multiSimParent),
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	dep.MultiSim = int(multiSimParent)
	return nil
}

// MarkAsMultiSimChild marks this SimulationDeployment to be a child of the given simulation.
func (dep *SimulationDeployment) MarkAsMultiSimChild(tx *gorm.DB, parent *SimulationDeployment) *ign.ErrMsg {
	// Add or Update
	if err := tx.Where("group_id = ?", *dep.GroupID).Assign(SimulationDeployment{
		MultiSim:      int(multiSimChild),
		ParentGroupID: parent.GroupID,
	}).FirstOrCreate(&dep).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	return nil
}

func (dep *SimulationDeployment) setStatus(status simulations.Status) {
	dep.DeploymentStatus = convertStatus(status).ToPtr()
}

func convertStatus(status simulations.Status) DeploymentStatus {
	switch status {
	case simulations.StatusPending:
		return simPending
	case simulations.StatusLaunchingInstances:
		return simLaunchingNodes
	case simulations.StatusLaunchingPods:
		return simLaunchingPods
	case simulations.StatusRunning:
		return simRunning
	case simulations.StatusTerminateRequested:
		return simTerminateRequested
	case simulations.StatusDeletingPods:
		return simDeletingPods
	case simulations.StatusDeletingNodes:
		return simDeletingNodes
	case simulations.StatusTerminatingInstances:
		return simTerminatingInstances
	case simulations.StatusTerminated:
		return simTerminated
	case simulations.StatusRejected:
		return simRejected
	case simulations.StatusSuperseded:
		return simSuperseded
	default:
		return simPending
	}
}

// SimulationDeployments is a slice of SimulationDeployment
type SimulationDeployments []SimulationDeployment

// DeploymentStatus are the possible status values of SimulationDeployments
type DeploymentStatus int

const (
	// Dev note: these statuses are compared by their ordinals. Any Status after
	// the 'simTerminated' status should be for Terminated simulations. Any status
	// before the 'simTerminated' is used for running (or about to run) simulations.

	// When defining new statuses, follow this rule:
	// Whenever you want to add a new status, add the status values of the statuses
	// that should come before and after it and then divide by 2. This makes it so
	// that the status value is placed in the middle of the two values, and leaves
	// space for new statuses to be placed in the future.

	simPending        DeploymentStatus = 0
	simLaunchingNodes DeploymentStatus = 10
	simLaunchingPods  DeploymentStatus = 20
	// simParentLaunching is only used for Parent simulations when some of their children
	// are still launching and there wasn't any errors so far.
	simParentLaunching DeploymentStatus = 25
	// simParentLaunchingWithErrors is only used for Parent simulations when some of their children
	// finished with errors and some are still launching/running.
	simParentLaunchingWithErrors DeploymentStatus = 28
	simRunning                   DeploymentStatus = 30
	// simRunningWithErrors is used for Parent simulations when some of their children
	// finished with errors and some are still running.
	// @deprecated do not use.
	simRunningWithErrors    DeploymentStatus = 40
	simTerminateRequested   DeploymentStatus = 50
	simDeletingPods         DeploymentStatus = 60
	simDeletingNodes        DeploymentStatus = 70
	simTerminatingInstances DeploymentStatus = 80
	simTerminated           DeploymentStatus = 90
	simRejected             DeploymentStatus = 100
	simSuperseded           DeploymentStatus = 110
)

// Corresponding string value for a Role
var depStatusStr = map[DeploymentStatus]string{
	simPending:                   "Pending",
	simLaunchingNodes:            "LaunchingNodes",
	simLaunchingPods:             "LaunchingPods",
	simParentLaunching:           "Launching",
	simParentLaunchingWithErrors: "RunningWithErrors",
	simRunning:                   "Running",
	simRunningWithErrors:         "RunningWithErrorsDoNotUse",
	simTerminateRequested:        "ToBeTerminated",
	simDeletingPods:              "DeletingPods",
	simDeletingNodes:             "DeletingNodes",
	simTerminatingInstances:      "TerminatingInstances",
	simTerminated:                "Terminated",
	simRejected:                  "Rejected",
	simSuperseded:                "Superseded",
}

// Eq function will compare for equality with an int based Status.
func (ds DeploymentStatus) Eq(st int) bool {
	return int(ds) == st
}

// String function will return the string version of the status.
func (ds DeploymentStatus) String() string {
	return depStatusStr[ds]
}

// ToPtr returns a pointer to int of this status value
func (ds DeploymentStatus) ToPtr() *int {
	i := int(ds)
	return &i
}

// ToInt returns the int value of this status value.
func (ds DeploymentStatus) ToInt() int {
	return int(ds)
}

// DeploymentStatusFrom returns the DeploymentStatus value corresponding to the
// given string. It will return nil if not found.
func DeploymentStatusFrom(str string) *DeploymentStatus {
	for k, v := range depStatusStr {
		if strings.EqualFold(v, str) {
			// Create a new DeploymentStatus to avoid sharing the map key (same pointer) with all callers
			ds := DeploymentStatus(int(k))
			return &ds
		}
	}
	return nil
}

// ErrorStatus are possible status values of SimulationDeployment ErrorStatus field.
type ErrorStatus string

const (
	simErrorWhenInitializing ErrorStatus = "InitializationFailed"
	simErrorWhenTerminating  ErrorStatus = "TerminationFailed"
	// Set when there was a second error during error handling. Marking for human review
	simErrorAdminReview ErrorStatus = "AdminReview"
	// Set when the simulation did not start due to a rejection by the SimService
	simErrorRejected ErrorStatus = "Rejected"
	// simErrorServerRestart is set by the server initialization process when it finds
	// Simulation Deployments left with intermediate statuses (either starting or terminating).
	// Having this error means that the server suffered a shutdown in the middle of a start
	// or terminate operation.
	simErrorServerRestart ErrorStatus = "ServerRestart"
	// Set when there was an error during log upload. Marking for human review
	simErrorFailedToUploadLogs ErrorStatus = "FailedToUploadLogs"
)

// weight returns the relative weight of an error status. It is used for internal
// comparisons.
func (es ErrorStatus) weight() int {
	switch es {
	case simErrorWhenInitializing, simErrorWhenTerminating:
		return 0
	case simErrorRejected:
		return 1
	case simErrorServerRestart:
		return 2
	case simErrorAdminReview, simErrorFailedToUploadLogs:
		return 5
	}
	// shouldn't be here
	panic("shouldn't be here")
}

// isMoreSevere checks which one of the errorStatus arguments is more severe.
// Returns true if es2 is more severe than es1.
func isMoreSevere(es1, es2 *ErrorStatus) bool {
	if es2 == nil {
		return false
	}
	if es1 == nil {
		return true
	}
	return es2.weight() > es1.weight()
}

// ToStringPtr returns a pointer to string of this status value
func (es ErrorStatus) ToStringPtr() *string {
	str := string(es)
	return &str
}

// ErrorStatusFrom returns the ErrorStatus value corresponding to the
// given string. It will return nil if not found.
func ErrorStatusFrom(str string) *ErrorStatus {
	switch strings.ToLower(str) {
	case strings.ToLower(string(simErrorWhenInitializing)):
		s := ErrorStatus(simErrorWhenInitializing)
		return &s
	case strings.ToLower(string(simErrorWhenTerminating)):
		s := ErrorStatus(simErrorWhenTerminating)
		return &s
	case strings.ToLower(string(simErrorAdminReview)):
		s := ErrorStatus(simErrorAdminReview)
		return &s
	case strings.ToLower(string(simErrorRejected)):
		s := ErrorStatus(simErrorRejected)
		return &s
	case strings.ToLower(string(simErrorServerRestart)):
		s := ErrorStatus(simErrorServerRestart)
		return &s
	case strings.ToLower(string(simErrorFailedToUploadLogs)):
		s := ErrorStatus(simErrorFailedToUploadLogs)
		return &s
	}
	return nil
}

// MachineInstance is a host/instance launched by Cloudsim.
// This structure is used by the ec2_machines module
type MachineInstance struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`

	InstanceID      *string `json:"instance_id" gorm:"not null;unique"`
	LastKnownStatus *string `json:"status,omitempty"`
	// Cloudsim Group Id
	GroupID *string `json:"group_id"`
	// Applicaton to which this machine belongs to
	Application *string `json:"application,omitempty"`
}

// updateMachineStatus updates the status of a given machine
func (m *MachineInstance) updateMachineStatus(tx *gorm.DB, st MachineStatus) *ign.ErrMsg {
	statusStr := st.ToStringPtr()
	if err := tx.Model(&m).Update(MachineInstance{
		LastKnownStatus: statusStr,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	m.LastKnownStatus = statusStr
	return nil
}

// MachineInstances is a MachineInstance slice
type MachineInstances []MachineInstance

func (m MachineInstances) getInstanceIDs() []*string {
	instanceIDs := make([]*string, len(m))
	for i, machine := range m {
		instanceIDs[i] = machine.InstanceID
	}

	return instanceIDs
}

// updateMachineStatuses updates the status of this set of machines.
func (m MachineInstances) updateMachineStatuses(ctx context.Context, tx *gorm.DB, st MachineStatus) *ign.ErrMsg {
	logger := logger(ctx)
	if m == nil {
		logger.Error("Attempted to update machine statuses with nil MachineInstances")
		return ign.NewErrorMessage(ign.ErrorUnexpected)
	} else if len(m) == 0 {
		logger.Warning("Attempted to update machine statuses for MachineInstances with length 0")
		return ign.NewErrorMessage(ign.ErrorUnexpected)
	}

	statusStr := st.ToStringPtr()

	if err := tx.Model(&m).Update(MachineInstance{
		LastKnownStatus: statusStr,
	}).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	for _, machine := range m {
		machine.LastKnownStatus = statusStr
	}

	return nil
}

// MachineStatus is a status of Host machines/instances
type MachineStatus string

const (
	macInitializing MachineStatus = "initializing"
	macRunning      MachineStatus = "running"
	macTerminating  MachineStatus = "terminating"
	macTerminated   MachineStatus = "terminated"
	macError        MachineStatus = "error"
)

// MachineStatusFrom returns the MachineStatus value corresponding to the
// given string. It will return nil if not found.
func MachineStatusFrom(str string) *MachineStatus {
	switch strings.ToLower(str) {
	case strings.ToLower(string(macInitializing)):
		s := MachineStatus(str)
		return &s
	case strings.ToLower(string(macRunning)):
		s := MachineStatus(str)
		return &s
	case strings.ToLower(string(macTerminating)):
		s := MachineStatus(str)
		return &s
	case strings.ToLower(string(macTerminated)):
		s := MachineStatus(str)
		return &s
	case strings.ToLower(string(macError)):
		s := MachineStatus(str)
		return &s
	}
	return nil
}

// ToStringPtr returns a pointer to string of this status value
func (ms MachineStatus) ToStringPtr() *string {
	str := string(ms)
	return &str
}

// GetMachine gets a machine instance record by its instanceID
// Fails if not found.
func GetMachine(tx *gorm.DB, instanceID string) (*MachineInstance, error) {
	var m MachineInstance
	if err := tx.Model(&MachineInstance{}).Where("instance_id = ?", instanceID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateSimulation contains information about a simulation creation request.
type CreateSimulation struct {
	// TODO Reenable notinblacklist validator for Name
	Name  string `json:"name" validate:"required,min=3,alphanum" form:"name"`
	Owner string `json:"owner" form:"owner"`
	// The docker image(s) that will be used for the Field Computer(s)
	Image       []string `json:"image" form:"image"`
	Platform    string   `json:"platform" form:"platform"`
	Application string   `json:"application" form:"application"`
	Private     *bool    `json:"private" validate:"omitempty" form:"private"`
	// When shutting down simulations, stop EC2 instances instead of terminating them. Requires admin privileges.
	StopOnEnd *bool `json:"stop_on_end" validate:"omitempty" form:"stop_on_end"`
	// Extra: it is expected that this field will be set by the Application logic
	// and not From form values.
	Extra *string `form:"-"`
	// ExtraSelector: it is expected that this field will be set by the Application logic
	// and not from Form values.
	ExtraSelector *string `form:"-"`
	// TODO: This is a field specific to SubT. As such this is a temporary field
	//  that should be included in the same separate table where Extra and
	//  ExtraSelector should reside.
	// Contains the names of all robots in the simulation in a comma-separated list.
	Robots *string `form:"-"`
}

// CustomRuleType defines the type for circuit custom rules
type CustomRuleType string

// List of rule types
const (
	// Maximum number of submissions allowed for a specific circuit
	MaxSubmissions CustomRuleType = "max_submissions"
)

// CircuitCustomRule holds custom rules for a specific combination of owner
// (user or organization) or circuit. Rules contain arbitrary values that can
// be used to configure specific aspects of a circuit/application (e.g.
// max_submissions - A custom rule for a specific owner to allow for
// extra submissions in a specific circuit). Rules for several owners or
// circuits can be defined by creating rules with NULL values in either fields.
// Rules with NULL values have less priority than rules with values, with the
// following priority: owner, circuit. This means that a rule with NULL circuit
// and owner will apply to ALL circuits for ALL owners, but any rule with either
// circuit or owner will override this general rule.
//
type CircuitCustomRule struct {
	gorm.Model
	Owner    *string        `json:"owner"`
	Circuit  *string        `json:"circuit" validate:"iscircuit"`
	RuleType CustomRuleType `gorm:"not null" json:"rule_type" validate:"isruletype"`
	Value    string         `gorm:"not null" json:"value"`
}

// CircuitCustomRules is a slice of CircuitCustomRule
type CircuitCustomRules []CircuitCustomRule

// GetCircuitCustomRule returns the rule value for a specific circuit and owner.
func GetCircuitCustomRule(tx *gorm.DB, circuit string, owner string, rule CustomRuleType) (*CircuitCustomRule, error) {
	var c CircuitCustomRule
	if err := tx.Model(&CircuitCustomRule{}).
		// Allow empty values for rules
		Where("rule_type = ?", rule).
		Where("circuit = ? OR circuit IS NULL", circuit).
		Where("owner = ? OR owner IS NULL", owner).
		// Set rule priority
		Order("owner DESC, circuit DESC").
		Limit(1).
		First(&c).
		Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// getRemainingSubmissions returns the number of remaining submissions compared to the Circuit custom rules, if any.
// The result can be a negative value, to indicate how large is the existing difference between the rule and the number
// of previously submitted simulations.
func getRemainingSubmissions(tx *gorm.DB, circuit string, owner string) (*int, error) {
	// Check if the owner has a specific number of max submissions, fallback
	// to unlimited submissions otherwise.
	limit := 0
	customRule, err := GetCircuitCustomRule(tx, circuit, owner, MaxSubmissions)
	if err == nil {
		limit, _ = strconv.Atoi(customRule.Value)
	}

	// Get the number of submissions by the owner
	count, err := countSimulationsByCircuit(tx, owner, circuit)
	if err != nil {
		return nil, err
	}

	// If there's no limit in place, simply return nil
	if limit == 0 {
		return nil, nil
	}

	remaining := limit - *count

	return &remaining, nil
}

// PodLog describes a pod log from kubernetes
type PodLog struct {
	Log string `json:"log"`
}

// LogGateway describes a response from the logs gateway handler.
type LogGateway struct {
	Path   string `json:"path"`
	IsFile bool   `json:"is_file"`
}

// SchedulableTask describes a task that can be given to a scheduler to be run at a certain time.
type SchedulableTask struct {
	Fn   func()
	Date time.Time
}
