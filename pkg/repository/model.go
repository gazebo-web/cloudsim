package repository

// Model represents a generic entity. A Model is part of the domain layer and is  persisted by a certain Repository.
type Model interface {
	// TableName returns the table/collection name for a certain model.
	TableName() string
}
