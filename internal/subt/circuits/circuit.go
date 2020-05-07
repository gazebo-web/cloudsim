package circuits

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Circuit struct {
	gorm.Model
	Name        *string `gorm:"not null;unique" json:"-"`
	Image       *string `json:"-"`
	BridgeImage *string `json:"-"`
	Worlds      *string `gorm:"size:2048" json:"-"`
	Times       *string `json:"-"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	WorldStatsTopics *string `gorm:"size:2048" json:"-"`
	// Topic used to track when the simulation officially starts and ends
	WorldWarmupTopics *string `gorm:"size:2048" json:"-"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	WorldMaxSimSeconds *string `json:"-"`
	// A comma separated list of seed numbers. Each seed will be used with each world run.
	// As an example, if field "Worlds" contains 3 worlds and "times" contains "1,2,2", then
	// there should be 5 seeds.
	Seeds           *string    `json:"-"`
	MaxCredits      *int       `json:"-"`
	CompetitionDate *time.Time `json:"-"`
	// If this field is set to true, every team that has qualified for this circuit
	// must be added to the table sub_t_qualified_participant.
	// All the participants that were not added to the qualified participants table will be rejected when submitting
	// a new simulation for this circuit.
	RequiresQualification *bool `json:"-"`
	Enabled               bool  `json:"-"`
}

func (Circuit) TableName() string {
	return "subt_circuits"
}
