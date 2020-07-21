package repositories

import "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"

// gormMap is a specific gorm-related map implementation.
type gormMap struct{
	tableName string
	singularName string
	pluralName string
	Map map[string]interface{}
}

// TableName returns the entity's table name.
func (g gormMap) TableName() string {
	return g.tableName
}

// SingularName returns the entity's name in singular.
func (g gormMap) SingularName() string {
	return g.singularName
}

// PluralName returns the entity's name in plural.
func (g gormMap) PluralName() string {
	return g.pluralName
}

// NewGormMap initializes a new gormMap that implements the Entity interface.
func NewGormMap(input map[string]interface{}, entity domain.Entity) domain.Entity {
	return &gormMap{
		tableName: entity.TableName(),
		singularName: entity.SingularName(),
		pluralName: entity.PluralName(),
		Map: input,
	}
}
