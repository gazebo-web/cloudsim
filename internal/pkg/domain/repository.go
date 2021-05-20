package domain

// Repository represents a generic repository layer interface.
type Repository interface {
	Create([]Entity) ([]Entity, error)
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	FindOne(entity Entity, filters ...Filter) error
	Update(data interface{}, filters ...Filter) error
	Delete(filters ...Filter) error
	SingularName() string
	PluralName() string
	Model() Entity
}
