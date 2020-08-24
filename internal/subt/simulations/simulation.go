package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Simulation represents a set of aggregated values of a Cloudsim Simulation for SubT.
type Simulation struct {
	gorm.Model
	Base *simulations.Simulation `gorm:"foreignkey:Sim" json:"-"`
	// Simulation unique identifier
	GroupID *string `gorm:"not null;unique" json:"-"`
	// Simulation score
	Score *float64 `gorm:"not null" json:"score"`
	// Simulation run info
	SimTimeDurationSec  int `gorm:"not null" json:"sim_time_duration_sec"`
	RealTimeDurationSec int `gorm:"not null" json:"real_time_duration_sec"`
	ModelCount          int `gorm:"not null" json:"model_count"`
}

func (Simulation) TableName() string {
	return "subt_simulations"
}

func (s *Simulation) Input() *simulations.Simulation {
	return s.Base
}

func (s *Simulation) Output() *simulations.Simulation {
	return s.Base
}

func (s *Simulation) ChildInput() *Simulation {
	return s
}
