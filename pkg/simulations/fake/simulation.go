package fake

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/calculator"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulations"
	"time"
)

// fakeSimulation is a fake simulations.Simulation implementation.
type fakeSimulation struct {
	groupID    simulations.GroupID
	status     simulations.Status
	kind       simulations.Kind
	err        *simulations.Error
	image      string
	validFor   time.Duration
	processed  bool
	owner      *string
	creator    string
	platform   *string
	launchedAt *time.Time
	rate       *calculator.Rate
	stoppedAt  *time.Time
}

// GetChargedAt mocks the GetChargedAt method.
func (f *fakeSimulation) GetChargedAt() *time.Time {
	now := time.Now()
	return &now
}

// GetRate returns the rate.
func (f *fakeSimulation) GetRate() calculator.Rate {
	if f.rate != nil {
		return *f.rate
	}
	return calculator.Rate{
		Amount:    0,
		Currency:  "usd",
		Frequency: time.Hour,
	}
}

// GetStoppedAt returns the stopped at field.
func (f *fakeSimulation) GetStoppedAt() *time.Time {
	return f.stoppedAt
}

// GetCost mocks the GetCost method.
func (f *fakeSimulation) GetCost() (uint, calculator.Rate, error) {
	if f.rate == nil {
		return 0, calculator.Rate{}, nil
	}
	return 0, *f.rate, nil
}

// SetRate sets the given rate.
func (f *fakeSimulation) SetRate(rate calculator.Rate) {
	f.rate = &rate
}

// IsProcessed returns if the simulation is processed.
func (f *fakeSimulation) IsProcessed() bool {
	return f.processed
}

// GetOwner returns the simulation's owner.
func (f *fakeSimulation) GetOwner() *string {
	return f.owner
}

// GetCreator returns the simulation's creator.
func (f *fakeSimulation) GetCreator() string {
	return f.creator
}

// GetValidFor returns the valid duration.
func (f *fakeSimulation) GetValidFor() time.Duration {
	return f.validFor
}

// HasStatus returns true if the given status matches with the current status.
func (f *fakeSimulation) HasStatus(status simulations.Status) bool {
	return f.status == status
}

// IsKind returns true if the given kind matches with the current kind.
func (f *fakeSimulation) IsKind(kind simulations.Kind) bool {
	return f.kind == kind
}

// GetImage returns the fake simulation's image.
func (f *fakeSimulation) GetImage() string {
	return f.image
}

// GetLaunchedAt returns the time and date the fake simulation was launched.
func (f *fakeSimulation) GetLaunchedAt() *time.Time {
	return f.launchedAt
}

// GetError returns the fake simulation's error.
// It returns nil if no error has been set.
func (f *fakeSimulation) GetError() *simulations.Error {
	return f.err
}

// GetGroupID returns the fake simulation's group id.
func (f *fakeSimulation) GetGroupID() simulations.GroupID {
	return f.groupID
}

// GetStatus returns the fake simulation's status.
func (f *fakeSimulation) GetStatus() simulations.Status {
	return f.status
}

// SetStatus sets the fake simulation's status to the given status.
func (f *fakeSimulation) SetStatus(status simulations.Status) {
	f.status = status
}

// GetKind returns the simulation's kind.
func (f *fakeSimulation) GetKind() simulations.Kind {
	return f.kind
}

// GetPlatform returns the simulation's platform.
func (f *fakeSimulation) GetPlatform() *string {
	return f.platform
}

// NewSimulation initializes a new fake simulation.
func NewSimulation(groupID simulations.GroupID, status simulations.Status, kind simulations.Kind,
	err *simulations.Error, image string, validFor time.Duration, owner *string,
	launchedAt *time.Time) simulations.Simulation {

	return &fakeSimulation{
		groupID:    groupID,
		status:     status,
		kind:       kind,
		err:        err,
		image:      image,
		validFor:   validFor,
		owner:      owner,
		launchedAt: launchedAt,
	}
}
