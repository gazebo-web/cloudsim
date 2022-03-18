package repository

import "errors"

var (
	// ErrNoFilter represents an error when no filter are provided.
	ErrNoFilter = errors.New("no filters provided")
	// ErrNoEntriesUpdated represent an error when no entries were updated in the database
	// after an Update operation.
	ErrNoEntriesUpdated = errors.New("no entries were updated")
	// ErrNoEntriesDeleted represent an error when no entries were deleted in the database
	// after a Delete operation.
	ErrNoEntriesDeleted = errors.New("no entries were deleted")
)

// Repository holds methods to CRUD an entity on a certain persistence layer.
type Repository interface {
	// Create is a bulk operation to create multiple entries with a single operation.
	//	entities: should be a slice of the same data structure implementing Model.
	Create(entities []Model) ([]Model, error)
	// Find filters entries and stores filtered entries in output.
	//	output: will contain the result of the query. It must be a pointer to a slice.
	//	offset: defines the number of results to skip before loading values to output.
	//	limit: defines the maximum number of entries to return. A nil value returns infinite results.
	// 	filters: filter entries by field value.
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	// FindOne filters entries and stores the first filtered entry in output, it must be a pointer to
	// a data structure implementing Model.
	FindOne(output Model, filters ...Filter) error
	// Update updates all model entries that match the provided filters with the given data.
	//	data: must be a map[string]interface{}
	//  filters: filter entries that should be updated.
	Update(data interface{}, filters ...Filter) error
	// Delete removes all the model entries that match filters.
	//  filters: filter entries that should be deleted.
	Delete(filters ...Filter) error
	// Model returns this repository's model.
	Model() Model
}
