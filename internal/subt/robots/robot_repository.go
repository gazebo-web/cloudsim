package robots

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/robots"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
)

type IRepository interface {
	GetAll() (map[string]Robot, error)
	GetByType(robotType string) (*Robot, error)
}

type Repository struct {
	// DB *gorm.DB
	robots map[string]Robot
}

func (r *Repository) GetByType(robotType string) (*Robot, error) {
	robot, ok := r.robots[robotType]
	if !ok {
		return nil, errors.New("robot type doesn't exist")
	}
	return &robot, nil
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		robots: generateRobotsMap(),
	}
	return r
}

func generateRobotsMap() map[string]Robot {
	owner := "openrobotics"
	return map[string]Robot{
		"X1_SENSOR_CONFIG_1": generateRobot(
			owner,
			"X1 Config 1",
			"X1_SENSOR_CONFIG_1",
			270,
		),
		"X1_SENSOR_CONFIG_2": generateRobot(
			owner,
			"X1 Config 2",
			"X1_SENSOR_CONFIG_2",
			300,
		),
		"X1_SENSOR_CONFIG_3": generateRobot(
			owner,
			"X1 Config 3",
			"X1_SENSOR_CONFIG_3",
			320,
		),
		"X1_SENSOR_CONFIG_4": generateRobot(
			owner,
			"X1 Config 4",
			"X1_SENSOR_CONFIG_4",
			370,
		),
		"X1_SENSOR_CONFIG_5": generateRobot(
			owner,
			"X1 Config 5",
			"X1_SENSOR_CONFIG_5",
			290,
		),
		"X1_SENSOR_CONFIG_6": generateRobot(

			owner,
			"X1 Config 6",
			"X1_SENSOR_CONFIG_6",
			380,
		),
		"EXPLORER_X1_SENSOR_CONFIG_1": generateRobot(
			owner,
			"EXPLORER_X1_SENSOR_CONFIG_1",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			390,
		),
		"X2_SENSOR_CONFIG_1": generateRobot(
			owner,
			"X2 Config 1",
			"X2_SENSOR_CONFIG_1",
			150,
		),
		"X2_SENSOR_CONFIG_2": generateRobot(
			owner,
			"X2 Config 2",
			"X2_SENSOR_CONFIG_2",
			160,
		),
		"X2_SENSOR_CONFIG_3": generateRobot(
			owner,
			"X2 Config 3",
			"X2_SENSOR_CONFIG_3",
			170,
		),
		"X2_SENSOR_CONFIG_4": generateRobot(
			owner,
			"X2 Config 4",
			"X2_SENSOR_CONFIG_4",
			180,
		),
		"X2_SENSOR_CONFIG_5": generateRobot(
			owner,
			"X2 Config 5",
			"X2_SENSOR_CONFIG_5",
			170,
		),
		"X2_SENSOR_CONFIG_6": generateRobot(
			owner,
			"X2 Config 6",
			"X2_SENSOR_CONFIG_6",
			250,
		),
		"X2_SENSOR_CONFIG_7": generateRobot(
			owner,
			"X2 Config 7",
			"X2_SENSOR_CONFIG_7",
			260,
		),
		"ROBOTIKA_X2_SENSOR_CONFIG_1": generateRobot(
			owner,
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			190,
		),
		"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1": generateRobot(
			owner,
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			180,
		),
		"SSCI_X2_SENSOR_CONFIG_1": generateRobot(
			owner,
			"SSCI_X2_SENSOR_CONFIG_1",
			"SSCI_X2_SENSOR_CONFIG_1",
			230,
		),
		"X3_SENSOR_CONFIG_1": generateRobot(
			owner,
			"X3 UAV Config 1",
			"X3_SENSOR_CONFIG_1",
			80,
		),
		"X3_SENSOR_CONFIG_2": generateRobot(
			owner,
			"X3 UAV Config 2",
			"X3_SENSOR_CONFIG_2",
			90,
		),
		"X3_SENSOR_CONFIG_3": generateRobot(
			owner,
			"X3 UAV Config 3",
			"X3_SENSOR_CONFIG_3",
			90,
		),
		"X3_SENSOR_CONFIG_4": generateRobot(
			owner,
			"X3 UAV Config 4",
			"X3_SENSOR_CONFIG_4",
			100,
		),
		"X3_SENSOR_CONFIG_5": generateRobot(
			owner,
			"X3 UAV Config 5",
			"X3_SENSOR_CONFIG_5",
			110,
		),
		"X4_SENSOR_CONFIG_1": generateRobot(
			owner,
			"X4 UAV Config 1",
			"X4_SENSOR_CONFIG_1",
			130,
		),
		"X4_SENSOR_CONFIG_2": generateRobot(
			owner,
			"X4 UAV Config 2",
			"X4_SENSOR_CONFIG_2",
			130,
		),
		"X4_SENSOR_CONFIG_3": generateRobot(
			owner,
			"X4 UAV Config 3",
			"X4_SENSOR_CONFIG_3",
			130,
		),
		"X4_SENSOR_CONFIG_4": generateRobot(
			owner,
			"X4 UAV Config 4",
			"X4_SENSOR_CONFIG_4",
			160,
		),
		"X4_SENSOR_CONFIG_5": generateRobot(
			owner,
			"X4 UAV Config 5",
			"X4_SENSOR_CONFIG_5",
			150,
		),
		"X4_SENSOR_CONFIG_6": generateRobot(
			owner,
			"X4 UAV Config 6",
			"X4_SENSOR_CONFIG_6",
			140,
		),
		"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1": generateRobot(
			owner,
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			160,
		),
		"SSCI_X4_SENSOR_CONFIG_1": generateRobot(
			owner,
			"SSCI_X4_SENSOR_CONFIG_1",
			"SSCI_X4_SENSOR_CONFIG_1",
			200,
		),
		"SSCI_X4_SENSOR_CONFIG_2": generateRobot(
			owner,
			"SSCI_X4_SENSOR_CONFIG_2",
			"SSCI_X4_SENSOR_CONFIG_2",
			185,
		),
		"COSTAR_HUSKY_SENSOR_CONFIG_1": generateRobot(
			owner,
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			335,
		),
	}
}

func generateRobot(owner, name, robotType string, credits int) Robot {
	return Robot{
		Robot: robots.Robot{
			Name:  name,
			Owner: owner,
			Type:  robotType,
		},
		Credits:   credits,
		Thumbnail: tools.GenerateThumbnailURI("https://fuel.ignitionrobotics.org/1.0", owner, name, 1),
	}
}
