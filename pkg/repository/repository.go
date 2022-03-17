package repository

// Repository holds methods to CRUD an entity on a certain persistence layer.
type Repository interface {
	// CreateBulk saves the given set of input entities in a persistence layer.
	CreateBulk(in ...interface{}) error
	// Get writes on out the data of the entity identified by id.
	Get(out interface{}, id uint) error
	// UpdateBulk overwrites the information given as the input on the different entities
	// identified by the given ids.
	UpdateBulk(in interface{}, ids ...uint) error
	// RemoveBulk removes the entities identified by the given ids.
	RemoveBulk(ids ...uint) error
}
