package simulations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"strings"
	"time"
)

// SubTCreateSimulation is a CreateSimulation request extension that adds specfic
// fields and validations for SubT application.
type SubTCreateSimulation struct {
	CreateSimulation `json:"-"`
	// Robot Names . Note the validate tag with the "dive" validation to validate each item
	// TODO Reenable notinblacklist validator for Name
	RobotName []string `json:"robot_name" validate:"gt=0,unique,dive,required,min=2,max=24,alphanum" form:"robot_name"`
	RobotType []string `json:"robot_type" validate:"lenEqFieldLen=RobotName,dive,isrobottype" form:"robot_type"`
	// Override the CreateSimulation Image field
	RobotImage []string `json:"robot_image" validate:"lenEqFieldLen=RobotName" form:"robot_image"`
	Marsupial  []string `json:"marsupial" form:"marsupial"`
	TeamBaseImage string `json:"team_base_image" form:"team_base_image"`
	Circuit    string   `json:"circuit" validate:"required,iscircuit" form:"circuit"`
}

// robotImagesBelongToECROwner checks if the field value is a valid SubT image.
// If an ECR image then it needs to below to the same owner.
func (cs *SubTCreateSimulation) robotImagesBelongToECROwner() bool {
	ownerWithUnderscores := strings.Replace(cs.Owner, " ", "_", -1)
	for _, image := range cs.RobotImage {
		// If it's not an ECR image, continue
		// HACK
		if !strings.Contains(image, "dkr.ecr.") && !strings.Contains(image, ".amazonaws.com") {
			continue
		}

		ss := strings.Split(image, "/")
		teamRepo := strings.Split(ss[len(ss)-1], ":")[0]
		if !strings.EqualFold(ownerWithUnderscores, teamRepo) {
			return false
		}
	}

	return true
}

// SubTRobot is an internal type used to describe a single SubT robot (field-computer) request.
type SubTRobot struct {
	Name    string
	Type    string
	Image   string
	Credits int
}

// SubTMarsupial is an internal type used to describe marsupial vehicles in SubT.
type SubTMarsupial struct {
	Parent string
	Child  string
}

// metadataSubT is a struct use to hold the Metadata information added by SubT to
// show to the user.
type metadataSubT struct {
	Circuit string      `json:"circuit,omitempty"`
	Robots  []SubTRobot `json:"robots,omitempty"`
}

// ExtraInfoSubT is a struct use to hold the Extra information added by SubT to
// SimulationDeployments DB records. If new internal fields are added, they should
// be removed inside GetSimulationDeployment() method.
type ExtraInfoSubT struct {
	Circuit    string      `json:"circuit,omitempty"`
	WorldIndex *int        `json:"world_index,omitempty"`
	Robots     []SubTRobot `json:"robots,omitempty"`
	Marsupials []SubTMarsupial `json:"marsupials,omitempty"`
	// Which "simulation run" number is this? It is computed based on the number of worlds in the circuit and
	// how many time to run them. For multiSims, the RunIndex can be seen as the child index.
	RunIndex *int `json:"run_index,omitempty"`
}

// ReadExtraInfoSubT reads the ExtraInfoSubT from a given simulation deployment.
func ReadExtraInfoSubT(dep *SimulationDeployment) (*ExtraInfoSubT, error) {
	var extra ExtraInfoSubT
	if err := json.Unmarshal([]byte(*dep.Extra), &extra); err != nil {
		return nil, err
	}
	return &extra, nil
}

// ToJSON marshals an ExtraInfoSubT into a json string.
func (e *ExtraInfoSubT) ToJSON() (*string, error) {
	byt, err := json.Marshal(*e)
	if err != nil {
		return nil, err
	}
	return sptr(string(byt)), nil
}

// countSimulationsByCircuit counts the number of simulations submitted by an owner
// to a circuit.
func countSimulationsByCircuit(tx *gorm.DB, owner, circuit string) (*int, error) {
	count := 0
	if err := tx.Model(&SimulationDeployment{}).
		// Only Top Level simulations (ie. not child sims from MultiSims)
		Where("multi_sim != ?", multiSimChild).
		Where("error_status IS NULL").
		Where("owner = ?", owner).
		Where("extra_selector = ?", circuit).
		Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}

// SubTCircuitRules holds the rules associated to a given circuit. Eg which worlds
// to run and how many times.
type SubTCircuitRules struct {
	gorm.Model
	Circuit     *string `gorm:"not null;unique" json:"-"`
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
}

// GetPendingCircuitRules gets a list of circuits that are scheduled for competition
func GetPendingCircuitRules(tx *gorm.DB) (*[]SubTCircuitRules, error) {
	var circuits []SubTCircuitRules
	if err := tx.Model(&SubTCircuitRules{}).Where("competition_date >= ?", time.Now()).Find(&circuits).Error; err != nil {
		return nil, err
	}
	return &circuits, nil
}

// GetCircuitRules gets the rules for a given circuit
func GetCircuitRules(tx *gorm.DB, circuit string) (*SubTCircuitRules, error) {
	var c SubTCircuitRules
	if err := tx.Model(&SubTCircuitRules{}).Where("circuit = ?", circuit).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// SubTQualifiedParticipant represents a qualification from a certain owner to participate in a circuit.
type SubTQualifiedParticipant struct {
	gorm.Model
	Circuit string `gorm:"not null" json:"circuit"`
	Owner   string `gorm:"not null" json:"owner"`
}

// IsOwnerQualifiedForCircuit returns true when an owner is qualified for certain circuit
// In any other cases, it returns false.
func IsOwnerQualifiedForCircuit(tx *gorm.DB, owner, circuit string, username string) bool {
	if globals.Permissions.IsSystemAdmin(username) {
		return true
	}

	var r *SubTCircuitRules
	var err error
	if r, err = GetCircuitRules(tx, circuit); err != nil {
		return false
	}

	if r.RequiresQualification == nil {
		return true
	}

	if !(*r.RequiresQualification) {
		return true
	}

	var q SubTQualifiedParticipant
	return tx.Model(&SubTQualifiedParticipant{}).
		Where("circuit = ? AND owner = ?", circuit, owner).
		First(&q).Error == nil
}

// SimulationDeploymentsSubTValue holds SubT-specific values associated to a given simulation deployment.
// E.g. specific simulation score, summary values, etc.
type SimulationDeploymentsSubTValue struct {
	gorm.Model
	SimulationDeployment *SimulationDeployment `gorm:"foreignkey:SimDep" json:"-"`
	// Simulation unique identifier
	GroupID *string `gorm:"not null;unique" json:"-"`
	// Simulation score
	Score *float64 `gorm:"not null" json:"score"`
	// Simulation run info
	SimTimeDurationSec  int `gorm:"not null" json:"sim_time_duration_sec"`
	RealTimeDurationSec int `gorm:"not null" json:"real_time_duration_sec"`
	ModelCount          int `gorm:"not null" json:"model_count"`
}

// AggregatedSubTSimulationValues contains the total score and average statistics of a group of simulation deployments.
// These simulations are typically all the child simulations of a multi-sim.
type AggregatedSubTSimulationValues struct {
	Score                  float64 `json:"-"`
	SimTimeDurationAvg     float64 `json:"sim_time_duration_avg"`
	SimTimeDurationStdDev  float64 `json:"sim_time_duration_std_dev"`
	RealTimeDurationAvg    float64 `json:"real_time_duration_avg"`
	RealTimeDurationStdDev float64 `json:"real_time_duration_std_dev"`
	ModelCountAvg          float64 `json:"model_count_avg"`
	ModelCountStdDev       float64 `json:"model_count_std_dev"`
	Sources                string  `json:"-"`
}

// GetAggregatedSubTSimulationValues returns the total score and average and standard deviation of a group of
// simulations.
func GetAggregatedSubTSimulationValues(tx *gorm.DB, simDep *SimulationDeployment) (*AggregatedSubTSimulationValues, error) {
	var values AggregatedSubTSimulationValues
	if simDep.isMultiSimChild() {
		return nil, errors.New("cannot aggregate values for multisim children")
	}

	tableName := tx.NewScope(SimulationDeploymentsSubTValue{}).TableName()
	if err := tx.Table(tableName).
		Select(`SUM(score) AS score,
			   AVG(sim_time_duration_sec) AS sim_time_duration_avg,
			   STD(sim_time_duration_sec) AS sim_time_duration_std_dev,
			   AVG(real_time_duration_sec) AS real_time_duration_avg,
			   STD(real_time_duration_sec) AS real_time_duration_std_dev,
			   AVG(model_count) AS model_count_avg,
			   STD(model_count) AS model_count_std_dev,
			   GROUP_CONCAT(group_id SEPARATOR ',') AS sources`).
		Where("group_id LIKE ?", fmt.Sprintf("%%%s%%", *simDep.GroupID)).
		Where("deleted_at IS NULL").
		Scan(&values).
		Error; err != nil {
		return nil, err
	}

	return &values, nil
}
