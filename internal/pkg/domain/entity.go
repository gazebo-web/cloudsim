package domain

// Entity represents a generic entity.
type Entity interface {
	Namer
}

// Namer has methods to name an entity.
type Namer interface {
	// TableName returns the table/collection name for a certain entity.
	TableName() string
	// SingularName returns the entity's name in singular.
	SingularName() string
	// PluralName returns the entity's name in plural.
	PluralName() string
}
