package quals

import "github.com/jinzhu/gorm"

// Qualification represents an owner that's been qualified for a circuit.
type Qualification struct {
	gorm.Model
	Circuit string
	Owner   string
}

// TableName defines the table name for the Qualification entity.
func (Qualification) TableName() string {
	return "subt_qualifications"
}
