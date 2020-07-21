# Repositories package
The goal of this document is to explain how to use the generic repository interface. It currently supports the following ORMs:
- Gorm

## Interfaces
This package groups a set of interfaces to `Create`, `Find`, `Update` and `Delete` entities from a certain database. It makes use of the `Entity` interface to define tables that will have the information provided with the repository.

### Entity interface
The entity interface groups a set of methods to define the naming convention of a certain entity in the whole system. It also has a method that returns the table name for that specific entity.
```golang
type Entity interface {
	// TableName returns the table/collection name for a certain entity.
	TableName() string
	// SingularName returns the entity's name in singular.
	SingularName() string
	// PluralName returns the entity's name in plural.
	PluralName() string
}
```

### Repository interface
The repository interface, as described above in this section, includes a set of methods to perform CRUD operations with a certain entity.
It makes use of another interface called Filter in order to select the entities that will be used to perform those CRUD operations.

```golang
type Repository interface {
	Create([]domain.Entity) ([]domain.Entity, error)
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	FindOne(entity domain.Entity, filters ...Filter) error
	Update(data domain.Entity, filters ...Filter) error
	Delete(filters ...Filter) error
	SingularName() string
	PluralName() string
	Model() domain.Entity
}
```

### Filter interface
```golang
type Filter interface {
	Template() string
	Values() []interface{}
}
```

## GORM
We've created a gorm implementation that satisfies the interfaces described before. In the following section you'll find a couple of examples to help you understand how to use this specific implementation with your codebase.

### Setup

### Repository

### Filter

### Map