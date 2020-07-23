package domain

// Entity represents a generic entity. Entities represent models to be used with a certain Repository.
type Entity interface {
	// TableName returns the table/collection name for a certain entity.
	TableName() string
	// SingularName returns the entity's name in singular.
	SingularName() string
	// PluralName returns the entity's name in plural.
	PluralName() string
}
