package domain

// DTO represents a generic Data Transfer Object.
type DTO interface {
	// Value converts the given DTO into a struct accepted by the entity's repository layer.
	Value() (interface{}, error)
}
