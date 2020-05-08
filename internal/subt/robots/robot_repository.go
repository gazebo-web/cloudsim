package robots

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
)

type IRepository interface {
	GetAllConfigs() (map[string]RobotConfig, error)
	GetConfigByType(robotType string) (*RobotConfig, error)
}

type Repository struct {
	// DB *gorm.DB
	robotConfigs map[string]RobotConfig
}

func (r *Repository) GetAllConfigs() (map[string]RobotConfig, error) {
	panic("implement me")
}

func (r *Repository) GetConfigByType(robotType string) (*RobotConfig, error) {
	robot, ok := r.robotConfigs[robotType]
	if !ok {
		return nil, errors.New("robot type doesn't exist")
	}
	return &robot, nil
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		robotConfigs: generateRobotsMap(),
	}
	return r
}

func generateRobotsMap() map[string]RobotConfig {
	owner := "openrobotics"
	return map[string]RobotConfig{
		"X1_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"X1 Config 1",
			"X1_SENSOR_CONFIG_1",
			270,
		),
		"X1_SENSOR_CONFIG_2": generateRobotConfig(
			owner,
			"X1 Config 2",
			"X1_SENSOR_CONFIG_2",
			300,
		),
		"X1_SENSOR_CONFIG_3": generateRobotConfig(
			owner,
			"X1 Config 3",
			"X1_SENSOR_CONFIG_3",
			320,
		),
		"X1_SENSOR_CONFIG_4": generateRobotConfig(
			owner,
			"X1 Config 4",
			"X1_SENSOR_CONFIG_4",
			370,
		),
		"X1_SENSOR_CONFIG_5": generateRobotConfig(
			owner,
			"X1 Config 5",
			"X1_SENSOR_CONFIG_5",
			290,
		),
		"X1_SENSOR_CONFIG_6": generateRobotConfig(

			owner,
			"X1 Config 6",
			"X1_SENSOR_CONFIG_6",
			380,
		),
		"EXPLORER_X1_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"EXPLORER_X1_SENSOR_CONFIG_1",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			390,
		),
		"X2_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"X2 Config 1",
			"X2_SENSOR_CONFIG_1",
			150,
		),
		"X2_SENSOR_CONFIG_2": generateRobotConfig(
			owner,
			"X2 Config 2",
			"X2_SENSOR_CONFIG_2",
			160,
		),
		"X2_SENSOR_CONFIG_3": generateRobotConfig(
			owner,
			"X2 Config 3",
			"X2_SENSOR_CONFIG_3",
			170,
		),
		"X2_SENSOR_CONFIG_4": generateRobotConfig(
			owner,
			"X2 Config 4",
			"X2_SENSOR_CONFIG_4",
			180,
		),
		"X2_SENSOR_CONFIG_5": generateRobotConfig(
			owner,
			"X2 Config 5",
			"X2_SENSOR_CONFIG_5",
			170,
		),
		"X2_SENSOR_CONFIG_6": generateRobotConfig(
			owner,
			"X2 Config 6",
			"X2_SENSOR_CONFIG_6",
			250,
		),
		"X2_SENSOR_CONFIG_7": generateRobotConfig(
			owner,
			"X2 Config 7",
			"X2_SENSOR_CONFIG_7",
			260,
		),
		"ROBOTIKA_X2_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			190,
		),
		"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			180,
		),
		"SSCI_X2_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"SSCI_X2_SENSOR_CONFIG_1",
			"SSCI_X2_SENSOR_CONFIG_1",
			230,
		),
		"X3_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"X3 UAV Config 1",
			"X3_SENSOR_CONFIG_1",
			80,
		),
		"X3_SENSOR_CONFIG_2": generateRobotConfig(
			owner,
			"X3 UAV Config 2",
			"X3_SENSOR_CONFIG_2",
			90,
		),
		"X3_SENSOR_CONFIG_3": generateRobotConfig(
			owner,
			"X3 UAV Config 3",
			"X3_SENSOR_CONFIG_3",
			90,
		),
		"X3_SENSOR_CONFIG_4": generateRobotConfig(
			owner,
			"X3 UAV Config 4",
			"X3_SENSOR_CONFIG_4",
			100,
		),
		"X3_SENSOR_CONFIG_5": generateRobotConfig(
			owner,
			"X3 UAV Config 5",
			"X3_SENSOR_CONFIG_5",
			110,
		),
		"X4_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"X4 UAV Config 1",
			"X4_SENSOR_CONFIG_1",
			130,
		),
		"X4_SENSOR_CONFIG_2": generateRobotConfig(
			owner,
			"X4 UAV Config 2",
			"X4_SENSOR_CONFIG_2",
			130,
		),
		"X4_SENSOR_CONFIG_3": generateRobotConfig(
			owner,
			"X4 UAV Config 3",
			"X4_SENSOR_CONFIG_3",
			130,
		),
		"X4_SENSOR_CONFIG_4": generateRobotConfig(
			owner,
			"X4 UAV Config 4",
			"X4_SENSOR_CONFIG_4",
			160,
		),
		"X4_SENSOR_CONFIG_5": generateRobotConfig(
			owner,
			"X4 UAV Config 5",
			"X4_SENSOR_CONFIG_5",
			150,
		),
		"X4_SENSOR_CONFIG_6": generateRobotConfig(
			owner,
			"X4 UAV Config 6",
			"X4_SENSOR_CONFIG_6",
			140,
		),
		"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			160,
		),
		"SSCI_X4_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"SSCI_X4_SENSOR_CONFIG_1",
			"SSCI_X4_SENSOR_CONFIG_1",
			200,
		),
		"SSCI_X4_SENSOR_CONFIG_2": generateRobotConfig(
			owner,
			"SSCI_X4_SENSOR_CONFIG_2",
			"SSCI_X4_SENSOR_CONFIG_2",
			185,
		),
		"COSTAR_HUSKY_SENSOR_CONFIG_1": generateRobotConfig(
			owner,
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			335,
		),
	}
}

func generateRobotConfig(owner, name, robotType string, credits int) RobotConfig {
	return RobotConfig{
		Name:      name,
		Owner:     owner,
		Type:      robotType,
		Credits:   credits,
		Thumbnail: tools.GenerateThumbnailURI("https://fuel.ignitionrobotics.org/1.0", owner, name, 1),
	}
}
