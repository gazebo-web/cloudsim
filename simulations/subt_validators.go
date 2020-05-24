package simulations

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"html"
)

// InstallSubTCustomValidators extends validator.v9 with custom validation functions
// and meta tags for SubT simulations.
func InstallSubTCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("isrobottype", isValidRobotType)
	validate.RegisterValidation("iscircuit", isValidCircuit)
}

// SubTRobotType represents a SubT robot
type SubTRobotType struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Credits   int    `json:"credits"`
	Thumbnail string `json:"thumbnail"`
}

// SubTRobotTypes holds the list of available robot types
var SubTRobotTypes map[string]SubTRobotType

// generateThumbnailURI generates a thumbnail URI for a specific robot.
func generateThumbnailURI(cfg *subTSpecificsConfig, owner string, robotName string, thumbnailNo int) string {
	robotName = html.EscapeString(robotName)
	template := "%s/%s/models/%s/tip/files/thumbnails/%d.png"
	return fmt.Sprintf(template, cfg.FuelURL, owner, robotName, thumbnailNo)
}

// generateSubTRobotType creates a new SubTRobotType. It is setup as a function to allow using
// robot properties to generate a thumbnail.
func generateSubTRobotType(cfg *subTSpecificsConfig, owner string, robotName string, robotType string,
	credits int) SubTRobotType {
	return SubTRobotType{
		Owner:     owner,
		Name:      robotName,
		Type:      robotType,
		Credits:   credits,
		Thumbnail: generateThumbnailURI(cfg, owner, robotName, 1),
	}
}

// loadSubTRobotTypes populates the list of valid robot types. The list is not
// defined using a literal because the application config is required to
// initialize it, and this config is only loaded on startup.
func loadSubTRobotTypes(cfg *subTSpecificsConfig) {
	SubTRobotTypes = map[string]SubTRobotType{
		"X1_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 1",
			"X1_SENSOR_CONFIG_1",
			270,
		),
		"X1_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 2",
			"X1_SENSOR_CONFIG_2",
			300,
		),
		"X1_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 3",
			"X1_SENSOR_CONFIG_3",
			320,
		),
		"X1_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 4",
			"X1_SENSOR_CONFIG_4",
			370,
		),
		"X1_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 5",
			"X1_SENSOR_CONFIG_5",
			290,
		),
		"X1_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 6",
			"X1_SENSOR_CONFIG_6",
			380,
		),
		"EXPLORER_X1_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			390,
		),
		"X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 1",
			"X2_SENSOR_CONFIG_1",
			150,
		),
		"X2_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 2",
			"X2_SENSOR_CONFIG_2",
			160,
		),
		"X2_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 3",
			"X2_SENSOR_CONFIG_3",
			170,
		),
		"X2_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 4",
			"X2_SENSOR_CONFIG_4",
			180,
		),
		"X2_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 5",
			"X2_SENSOR_CONFIG_5",
			170,
		),
		"X2_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 6",
			"X2_SENSOR_CONFIG_6",
			250,
		),
		"X2_SENSOR_CONFIG_7": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 7",
			"X2_SENSOR_CONFIG_7",
			260,
		),
		"ROBOTIKA_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			190,
		),
		"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			180,
		),
		"SSCI_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X2_SENSOR_CONFIG_1",
			"SSCI_X2_SENSOR_CONFIG_1",
			230,
		),
		"X3_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 1",
			"X3_SENSOR_CONFIG_1",
			80,
		),
		"X3_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 2",
			"X3_SENSOR_CONFIG_2",
			90,
		),
		"X3_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 3",
			"X3_SENSOR_CONFIG_3",
			90,
		),
		"X3_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 4",
			"X3_SENSOR_CONFIG_4",
			100,
		),
		"X3_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 5",
			"X3_SENSOR_CONFIG_5",
			110,
		),
		"X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 1",
			"X4_SENSOR_CONFIG_1",
			130,
		),
		"X4_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 2",
			"X4_SENSOR_CONFIG_2",
			130,
		),
		"X4_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 3",
			"X4_SENSOR_CONFIG_3",
			130,
		),
		"X4_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 4",
			"X4_SENSOR_CONFIG_4",
			160,
		),
		"X4_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 5",
			"X4_SENSOR_CONFIG_5",
			150,
		),
		"X4_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 6",
			"X4_SENSOR_CONFIG_6",
			140,
		),
		"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			160,
		),
		"SSCI_X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X4_SENSOR_CONFIG_1",
			"SSCI_X4_SENSOR_CONFIG_1",
			200,
		),
		"SSCI_X4_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X4_SENSOR_CONFIG_2",
			"SSCI_X4_SENSOR_CONFIG_2",
			185,
		),
		"COSTAR_HUSKY_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			335,
		),
		"TEAMBASE": generateSubTRobotType(
			cfg,
			"openrobotics",
			"TEAMBASE",
			"TEAMBASE",
			0,
		),
	}
}

// SubTCircuits holds the list of available circuits in SubT.
var SubTCircuits = []string{
	CircuitVirtualStix,
	CircuitTunnelPractice1,
	CircuitTunnelPractice2,
	CircuitTunnelPractice3,
	CircuitSimpleTunnel1,
	CircuitSimpleTunnel2,
	CircuitSimpleTunnel3,
	CircuitUrbanQual,
	CircuitUrbanSimple1,
	CircuitUrbanSimple2,
	CircuitUrbanSimple3,
	CircuitUrbanPractice1,
	CircuitUrbanPractice2,
	CircuitUrbanPractice3,
	CircuitUrbanCircuitWorld1,
	CircuitUrbanCircuitWorld2,
	CircuitUrbanCircuitWorld3,
	CircuitUrbanCircuitWorld4,
	CircuitUrbanCircuitWorld5,
	CircuitUrbanCircuitWorld6,
	CircuitUrbanCircuitWorld7,
	CircuitUrbanCircuitWorld8,
	CircuitCaveSimple1,
	CircuitCaveSimple2,
	CircuitCaveSimple3,
	CircuitCaveQual,
}

// isValidRobotType checks if the field value is a valid Robot Type.
func isValidRobotType(fl validator.FieldLevel) bool {
	_, ok := SubTRobotTypes[strings.ToUpper(fl.Field().String())]
	return ok
}

// isValidCircuit checks if the field value is a valid SubT Circuit.
func isValidCircuit(fl validator.FieldLevel) bool {
	return StrSliceContains(fl.Field().String(), SubTCircuits)
}
