package simulations

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"html"
	"strings"
)

// InstallSubTCustomValidators extends validator.v9 with custom validation functions
// and meta tags for SubT simulations.
func InstallSubTCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("isrobottype", isValidRobotType)
	validate.RegisterValidation("iscircuit", isValidCircuit)
}

// SubTRobotType represents a SubT robot. These robots are expected to be part of the SubT portal in Fuel.
type SubTRobotType struct {
	// Owner is the organization that owns the robot in Fuel.
	Owner string `json:"owner"`
	// Name is the name of the robot shown to users.
	Name string `json:"name"`
	// Model is the robot model name. A single model can contain different sets of sensors.
	Model string `json:"model"`
	// Type contains the name of the robot's model and sensor combo.
	Type string `json:"type"`
	// Credits contains the cost of the robot.
	Credits int `json:"credits"`
	// Thumbnail contains the robot's thumbnail URL.
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
func generateSubTRobotType(cfg *subTSpecificsConfig, owner string, robotName string, robotModel string,
	robotType string, credits int) SubTRobotType {
	return SubTRobotType{
		Owner:     owner,
		Name:      robotName,
		Model:     robotModel,
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
			"X1",
			"X1_SENSOR_CONFIG_1",
			270,
		),
		"X1_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 2",
			"X1",
			"X1_SENSOR_CONFIG_2",
			300,
		),
		"X1_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 3",
			"X1",
			"X1_SENSOR_CONFIG_3",
			320,
		),
		"X1_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 4",
			"X1",
			"X1_SENSOR_CONFIG_4",
			370,
		),
		"X1_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 5",
			"X1",
			"X1_SENSOR_CONFIG_5",
			290,
		),
		"X1_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 6",
			"X1",
			"X1_SENSOR_CONFIG_6",
			380,
		),
		"X1_SENSOR_CONFIG_7": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 7",
			"X1",
			"X1_SENSOR_CONFIG_7",
			420,
		),
		"X1_SENSOR_CONFIG_8": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X1 Config 8",
			"X1",
			"X1_SENSOR_CONFIG_8",
			370,
		),
		"EXPLORER_X1_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			"X1",
			"EXPLORER_X1_SENSOR_CONFIG_1",
			390,
		),
		"EXPLORER_X1_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_X1_SENSOR_CONFIG_2",
			"X1",
			"EXPLORER_X1_SENSOR_CONFIG_2",
			440,
		),
		"X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 1",
			"X2",
			"X2_SENSOR_CONFIG_1",
			150,
		),
		"X2_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 2",
			"X2",
			"X2_SENSOR_CONFIG_2",
			160,
		),
		"X2_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 3",
			"X2",
			"X2_SENSOR_CONFIG_3",
			170,
		),
		"X2_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 4",
			"X2",
			"X2_SENSOR_CONFIG_4",
			180,
		),
		"X2_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 5",
			"X2",
			"X2_SENSOR_CONFIG_5",
			170,
		),
		"X2_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 6",
			"X2",
			"X2_SENSOR_CONFIG_6",
			250,
		),
		"X2_SENSOR_CONFIG_7": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 7",
			"X2",
			"X2_SENSOR_CONFIG_7",
			260,
		),
		"X2_SENSOR_CONFIG_8": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 8",
			"X2",
			"X2_SENSOR_CONFIG_8",
			275,
		),
		"X2_SENSOR_CONFIG_9": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X2 Config 9",
			"X2",
			"X2_SENSOR_CONFIG_9",
			205,
		),
		"ROBOTIKA_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			"X2",
			"ROBOTIKA_X2_SENSOR_CONFIG_1",
			190,
		),
		"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			"X2",
			"SOPHISTICATED_ENGINEERING_X2_SENSOR_CONFIG_1",
			180,
		),
		"SSCI_X2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X2_SENSOR_CONFIG_1",
			"X2",
			"SSCI_X2_SENSOR_CONFIG_1",
			230,
		),
		"X3_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 1",
			"X3",
			"X3_SENSOR_CONFIG_1",
			80,
		),
		"X3_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 2",
			"X3",
			"X3_SENSOR_CONFIG_2",
			90,
		),
		"X3_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 3",
			"X3",
			"X3_SENSOR_CONFIG_3",
			90,
		),
		"X3_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 4",
			"X3",
			"X3_SENSOR_CONFIG_4",
			100,
		),
		"X3_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X3 UAV Config 5",
			"X3",
			"X3_SENSOR_CONFIG_5",
			110,
		),
		"X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 1",
			"X4",
			"X4_SENSOR_CONFIG_1",
			130,
		),
		"X4_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 2",
			"X4",
			"X4_SENSOR_CONFIG_2",
			130,
		),
		"X4_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 3",
			"X4",
			"X4_SENSOR_CONFIG_3",
			130,
		),
		"X4_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 4",
			"X4",
			"X4_SENSOR_CONFIG_4",
			160,
		),
		"X4_SENSOR_CONFIG_5": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 5",
			"X4",
			"X4_SENSOR_CONFIG_5",
			150,
		),
		"X4_SENSOR_CONFIG_6": generateSubTRobotType(
			cfg,
			"openrobotics",
			"X4 UAV Config 6",
			"X4",
			"X4_SENSOR_CONFIG_6",
			140,
		),
		"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			"X4",
			"SOPHISTICATED_ENGINEERING_X4_SENSOR_CONFIG_1",
			160,
		),
		"SSCI_X4_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X4_SENSOR_CONFIG_1",
			"X4",
			"SSCI_X4_SENSOR_CONFIG_1",
			200,
		),
		"SSCI_X4_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"SSCI_X4_SENSOR_CONFIG_2",
			"X4",
			"SSCI_X4_SENSOR_CONFIG_2",
			185,
		),
		"COSTAR_HUSKY_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			"HUSKY",
			"COSTAR_HUSKY_SENSOR_CONFIG_1",
			165,
		),
		"COSTAR_HUSKY_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"COSTAR_HUSKY_SENSOR_CONFIG_2",
			"HUSKY",
			"COSTAR_HUSKY_SENSOR_CONFIG_2",
			190,
		),
		"TEAMBASE": generateSubTRobotType(
			cfg,
			"openrobotics",
			"TEAMBASE",
			"TEAMBASE",
			"TEAMBASE",
			0,
		),
		"CERBERUS_ANYMAL_B_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_ANYMAL_B_SENSOR_CONFIG_1",
			"ANYMAL_B",
			"CERBERUS_ANYMAL_B_SENSOR_CONFIG_1",
			215,
		),
		"CERBERUS_M100_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_M100_SENSOR_CONFIG_1",
			"M100",
			"CERBERUS_M100_SENSOR_CONFIG_1",
			95,
		),
		"ROBOTIKA_FREYJA_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_FREYJA_SENSOR_CONFIG_1",
			"FREYJA",
			"ROBOTIKA_FREYJA_SENSOR_CONFIG_1",
			155,
		),
		"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_1",
			"KLOUBAK",
			"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_1",
			145,
		),
		"MARBLE_HUSKY_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HUSKY_SENSOR_CONFIG_1",
			"HUSKY",
			"MARBLE_HUSKY_SENSOR_CONFIG_1",
			220,
		),
		"MARBLE_HD2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HD2_SENSOR_CONFIG_1",
			"HD2",
			"MARBLE_HD2_SENSOR_CONFIG_1",
			155,
		),
		"MARBLE_QAV500_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_QAV500_SENSOR_CONFIG_1",
			"QAV500",
			"MARBLE_QAV500_SENSOR_CONFIG_1",
			100,
		),
		"EXPLORER_R2_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_R2_SENSOR_CONFIG_1",
			"R2",
			"EXPLORER_R2_SENSOR_CONFIG_1",
			235,
		),
		"EXPLORER_R2_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_R2_SENSOR_CONFIG_2",
			"R2",
			"EXPLORER_R2_SENSOR_CONFIG_2",
			260,
		),
		"EXPLORER_DS1_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_DS1_SENSOR_CONFIG_1",
			"DS1",
			"EXPLORER_DS1_SENSOR_CONFIG_1",
			115,
		),
		"CERBERUS_ANYMAL_B_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_ANYMAL_B_SENSOR_CONFIG_2",
			"ANYMAL_B",
			"CERBERUS_ANYMAL_B_SENSOR_CONFIG_2",
			240,
		),
		"MARBLE_HUSKY_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HUSKY_SENSOR_CONFIG_2",
			"HUSKY",
			"MARBLE_HUSKY_SENSOR_CONFIG_2",
			245,
		),
		"MARBLE_HD2_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HD2_SENSOR_CONFIG_2",
			"HD2",
			"MARBLE_HD2_SENSOR_CONFIG_2",
			180,
		),
		"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_2",
			"KLOUBAK",
			"ROBOTIKA_KLOUBAK_SENSOR_CONFIG_2",
			170,
		),
		"ROBOTIKA_FREYJA_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"ROBOTIKA_FREYJA_SENSOR_CONFIG_2",
			"FREYJA",
			"ROBOTIKA_FREYJA_SENSOR_CONFIG_2",
			170,
		),
		"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_1",
			"OZBOT_ATR",
			"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_1",
			235,
		),
		"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_2",
			"OZBOT_ATR",
			"CSIRO_DATA61_OZBOT_ATR_SENSOR_CONFIG_2",
			260,
		),
		"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_1",
			"ABSOLEM",
			"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_1",
			155,
		),
		"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_2",
			"ABSOLEM",
			"CTU_CRAS_NORLAB_ABSOLEM_SENSOR_CONFIG_2",
			180,
		),
		"CERBERUS_GAGARIN_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_GAGARIN_SENSOR_CONFIG_1",
			"GAGARIN",
			"CERBERUS_GAGARIN_SENSOR_CONFIG_1",
			115,
		),
		"CERBERUS_ANYMAL_C_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_ANYMAL_C_SENSOR_CONFIG_1",
			"ANYMAL_C",
			"CERBERUS_ANYMAL_C_SENSOR_CONFIG_1",
			280,
		),
		"CERBERUS_ANYMAL_C_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_ANYMAL_C_SENSOR_CONFIG_2",
			"ANYMAL_C",
			"CERBERUS_ANYMAL_C_SENSOR_CONFIG_2",
			305,
		),
		"CERBERUS_RMF_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_RMF_SENSOR_CONFIG_1",
			"RMF",
			"CERBERUS_RMF_SENSOR_CONFIG_1",
			55,
		),
		"COSTAR_SHAFTER_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"COSTAR_SHAFTER_SENSOR_CONFIG_1",
			"SHAFTER",
			"COSTAR_SHAFTER_SENSOR_CONFIG_1",
			110,
		),
		"MARBLE_HD2_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HD2_SENSOR_CONFIG_3",
			"HD2",
			"MARBLE_HD2_SENSOR_CONFIG_3",
			165,
		),
		"MARBLE_HD2_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HD2_SENSOR_CONFIG_4",
			"HD2",
			"MARBLE_HD2_SENSOR_CONFIG_4",
			190,
		),
		"MARBLE_HUSKY_SENSOR_CONFIG_3": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HUSKY_SENSOR_CONFIG_3",
			"HD2",
			"MARBLE_HUSKY_SENSOR_CONFIG_3",
			230,
		),
		"MARBLE_HUSKY_SENSOR_CONFIG_4": generateSubTRobotType(
			cfg,
			"openrobotics",
			"MARBLE_HUSKY_SENSOR_CONFIG_4",
			"HD2",
			"MARBLE_HUSKY_SENSOR_CONFIG_4",
			255,
		),
		"CERBERUS_M100_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CERBERUS_M100_SENSOR_CONFIG_2",
			"M100",
			"CERBERUS_M100_SENSOR_CONFIG_2",
			105,
		),
		"EXPLORER_R3_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_R3_SENSOR_CONFIG_1",
			"R3",
			"EXPLORER_R3_SENSOR_CONFIG_1",
			235,
		),
		"EXPLORER_R3_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"EXPLORER_R3_SENSOR_CONFIG_2",
			"R3",
			"EXPLORER_R3_SENSOR_CONFIG_2",
			260,
		),
		"CSIRO_DATA61_DTR_SENSOR_CONFIG_1": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CSIRO_DATA61_DTR_SENSOR_CONFIG_1",
			"DTR",
			"CSIRO_DATA61_DTR_SENSOR_CONFIG_1",
			135,
		),
		"CSIRO_DATA61_DTR_SENSOR_CONFIG_2": generateSubTRobotType(
			cfg,
			"openrobotics",
			"CSIRO_DATA61_DTR_SENSOR_CONFIG_2",
			"DTR",
			"CSIRO_DATA61_DTR_SENSOR_CONFIG_2",
			160,
		),
	}
}

// SubTCircuits holds the list of available circuits in SubT.
var SubTCircuits = []string{
	CircuitNIOSHSRConfigA,
	CircuitNIOSHSRConfigB,
	CircuitNIOSHEXConfigA,
	CircuitNIOSHEXConfigB,
	CircuitVirtualStix,
	CircuitTunnelPractice1,
	CircuitTunnelPractice2,
	CircuitTunnelPractice3,
	CircuitSimpleTunnel1,
	CircuitSimpleTunnel2,
	CircuitSimpleTunnel3,
	CircuitTunnelCircuitWorld1,
	CircuitTunnelCircuitWorld2,
	CircuitTunnelCircuitWorld3,
	CircuitTunnelCircuitWorld4,
	CircuitTunnelCircuitWorld5,
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
	CircuitCavePractice1,
	CircuitCavePractice2,
	CircuitCavePractice3,
	CircuitCaveCircuit,
	CircuitCaveCircuitWorld1,
	CircuitCaveCircuitWorld2,
	CircuitCaveCircuitWorld3,
	CircuitCaveCircuitWorld4,
	CircuitCaveCircuitWorld5,
	CircuitCaveCircuitWorld6,
	CircuitCaveCircuitWorld7,
	CircuitCaveCircuitWorld8,
	CircuitFinalsQual,
	CircuitFinalsPractice1,
	CircuitFinalsPractice2,
	CircuitFinalsPractice3,
}

// SubTCompetitionCircuits is the list of circuits that are used for competitions.
var SubTCompetitionCircuits = []string{
	CircuitTunnelCircuit,
	CircuitUrbanCircuit,
	CircuitCaveCircuit,
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
