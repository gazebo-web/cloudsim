package summaries

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Service groups a set of methods to handle summaries for simulations.
type Service interface {
	// Save saves the given score, stats and run data for the given simulation's GroupID.
	Save(groupID simulations.GroupID, score *float64, stats simulations.Statistics, runData string) (*simulations.Summary, error)

	// Calculate takes a parent group id and returns the calculate summary for all its children.
	Calculate(groupID simulations.GroupID) (*simulations.Summary, error)
}

// service is the Summary Service implementation.
type service struct {
	db *gorm.DB
}

// Calculate takes a parent group id and returns the aggregate summary for all its children.
func (s *service) Calculate(groupID simulations.GroupID) (*simulations.Summary, error) {
	var summary simulations.Summary

	tableName := s.db.NewScope(simulations.Summary{}).TableName()
	q := s.db.Table(tableName).
		Select(`SUM(score) AS score,
			   AVG(sim_time_duration_sec) AS sim_time_duration_avg,
			   STD(sim_time_duration_sec) AS sim_time_duration_std_dev,
			   AVG(real_time_duration_sec) AS real_time_duration_avg,
			   STD(real_time_duration_sec) AS real_time_duration_std_dev,
			   AVG(model_count) AS model_count_avg,
			   STD(model_count) AS model_count_std_dev,
			   GROUP_CONCAT(group_id SEPARATOR ',') AS sources`).
		Where("group_id LIKE ?", fmt.Sprintf("%%%s%%", groupID)).
		Where("deleted_at IS NULL").
		Scan(&summary)

	if err := q.Error; err != nil {
		return nil, err
	}

	return &summary, nil
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
