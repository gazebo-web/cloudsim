package domain

// Entity represents a generic entity.
type Entity interface {
	Parser
	Namer
}

// Parser has methods to parse the entity.
type Parser interface {
	// ParseIn parses the given entity into a certain entity.
	ParseIn(input Entity) error
	// ParseOut parses the given entity and returns the parsed result.
	ParseOut(input interface{}) (Entity, error)
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
