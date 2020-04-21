package quals

import "github.com/jinzhu/gorm"

// Qualification represents an owner that's been qualified for a circuit.
type Qualification struct {
	gorm.Model
	Circuit string
	Owner string
}
