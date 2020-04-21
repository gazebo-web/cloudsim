package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type SubTSimulation struct {
	gorm.Model
	Simulation *simulations.Simulation `gorm:"foreignkey:SimDep" json:"-"`
	// Simulation unique identifier
	GroupID *string `gorm:"not null;unique" json:"-"`
	// Simulation score
	Score *float64 `gorm:"not null" json:"score"`
	// Simulation run info
	SimTimeDurationSec  int `gorm:"not null" json:"sim_time_duration_sec"`
	RealTimeDurationSec int `gorm:"not null" json:"real_time_duration_sec"`
	ModelCount          int `gorm:"not null" json:"model_count"`
}

func (SubTSimulation) TableName() string {
	return "subt_simulations"
}