package fake

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// fakeSimulation is a fake simulations.Simulation implementation.
type fakeSimulation struct {
	groupID   simulations.GroupID
	status    simulations.Status
	kind      simulations.Kind
	err       *simulations.Error
	image     string
	validFor  time.Duration
	processed bool
	owner     *string
	creator   string
}

func (f *fakeSimulation) GetOwner() *string {
	return f.owner
}

func (f *fakeSimulation) GetCreator() string {
	return f.creator
}

// GetValidFor returns the valid duration.
func (f *fakeSimulation) GetValidFor() time.Duration {
	return f.validFor
}

// IsProcessed returns true if the simulation has been processed.
func (f *fakeSimulation) IsProcessed() bool {
	return f.processed
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

// NewSimulation initializes a new fake simulation.
func NewSimulation(groupID simulations.GroupID, status simulations.Status, kind simulations.Kind, err *simulations.Error, image string, validFor time.Duration) simulations.Simulation {
	return &fakeSimulation{
		groupID:  groupID,
		status:   status,
		kind:     kind,
		err:      err,
		image:    image,
		validFor: validFor,
	}
}
