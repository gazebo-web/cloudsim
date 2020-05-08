package robots

import (
	"github.com/jinzhu/gorm"
)

type IRepository interface {
	GetAllConfigs() ([]RobotConfig, error)
	GetConfigByType(robotType string) (*RobotConfig, error)
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetAllConfigs() ([]RobotConfig, error) {
	panic("implement me")
}

func (r *Repository) GetConfigByType(robotType string) (*RobotConfig, error) {
	var config RobotConfig
	err := r.DB.Model(&RobotConfig{}).Where("type = ?", robotType).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		DB:db,
	}
	return r
}