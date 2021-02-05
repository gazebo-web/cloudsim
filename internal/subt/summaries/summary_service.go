package summaries

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Service groups a set of methods to handle summaries for simulations.
type Service interface {
	// Save saves the given score, stats and run data for the given simulation's GroupID.
	Save(groupID simulations.GroupID, score *float64, stats simulations.Statistics, runData string) (*simulations.Summary, error)
}

// service is the Summary Service implementation.
type service struct {
	db *gorm.DB
}

// Save saves the given score, stats and run data for the given simulation's GroupID.
func (s *service) Save(groupID simulations.GroupID, score *float64, stats simulations.Statistics, runData string) (*simulations.Summary, error) {
	summary := simulations.Summary{
		GroupID:                &groupID,
		Score:                  score,
		SimTimeDurationAvg:     float64(stats.SimulationTime),
		SimTimeDurationStdDev:  0,
		RealTimeDurationAvg:    float64(stats.RealTime),
		RealTimeDurationStdDev: 0,
		ModelCountAvg:          float64(stats.ModelCount),
		ModelCountStdDev:       0,
		Sources:                "",
		RunData:                runData,
	}

	if err := s.db.Model(&simulations.Summary{}).Create(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

// NewService initializes a new Service implementation using gorm.
func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
