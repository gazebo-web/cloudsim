package timestamp

import "time"

// ResourceTimestamp provides methods to access the creation and deletion timestamps of a Resource.
type ResourceTimestamp interface {
	// CreationTimestamp is a timestamp representing the server time when this object was created.
	CreationTimestamp() time.Time
	// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
	// graceful deletion is requested.
	DeletionTimestamp() *time.Time
}

type resourceTimestamp struct {
	// creationTimestamp is the ResourceTimestamp.CreationTimestamp.
	creationTimestamp time.Time
	// deletionTimestamp is the ResourceTimestamp.DeletionTimestamp.
	deletionTimestamp *time.Time
}

// CreationTimestamp is a timestamp representing the server time when this object was created.
func (s *resourceTimestamp) CreationTimestamp() time.Time {
	return s.creationTimestamp
}

// DeletionTimestamp is a timestamp at which this resource will be deleted. This field is set by the server when a
// graceful deletion is requested.
func (s *resourceTimestamp) DeletionTimestamp() *time.Time {
	return s.deletionTimestamp
}

// NewResourceTimestamp is used to initialize a new ResourceTimestamp implementation.
func NewResourceTimestamp(creation time.Time, deletion *time.Time) ResourceTimestamp {
	return &resourceTimestamp{
		creationTimestamp: creation,
		deletionTimestamp: deletion,
	}
}
