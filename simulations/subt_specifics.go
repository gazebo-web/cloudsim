package simulations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	platformManager "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform/manager"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

/*

	Logic specific to SubT, to describe needed Node instances, Pods, services, etc.
	The `SimService` and `NodeManager` will delegate to these functions when dealing
	with the competition specifics.

	Some dev notes about Network Policies:
		We are creating NetworkPolicies to only allow communication between a field-computer pod
		and a communication bridge pod, and between the communication bridge and the gazebo server pod.
		The communication bridge pod communicates with the field-computer pod using ign-transport and
		is in charge of filtering relevant topics for the robot.
		These policies also allow comunication with the cloudsim server (this is needed to support
		'ign-transport' with the cloudsim server).
		With these policies:
		- Pods from different Simulations groups should not be able to communicate.
		- Different field-computer pods belonging to the same simulation should not be able
		to communicate either.
		- field-computer pods should not be able to access internet.
		- The Gzserver Pod can access internet.
		If you notice any of the above not being true, please raise an issue.
*/

// SubT Specifics constants
const (
	subtTagKey string = "SubT"
	// A predefined const to refer to the SubT Platform type.
	// This will be used to provision the Nodes (Nvidia, CPU, etc)
	platformSubT string = "subt"

	// SubT resource type identifiers
	// These identifiers are used to tag AWS and Kubernetes resources related to each type of entity
	subtTypeGazebo        = "gzserver"
	subtTypeCommsBridge   = "comms-bridge"
	subtTypeFieldComputer = "field-computer"

	// A predefined const to refer to the SubT Application type
	// This will be used to know which Pods and services launch.
	applicationSubT                       string = "subt"
	CircuitNIOSHSRConfigA                 string = "NIOSH SR Config A"
	CircuitNIOSHSRConfigB                 string = "NIOSH SR Config B"
	CircuitNIOSHEXConfigA                 string = "NIOSH EX Config A"
	CircuitNIOSHEXConfigB                 string = "NIOSH EX Config B"
	CircuitVirtualStix                    string = "Virtual Stix"
	CircuitTunnelCircuit                  string = "Tunnel Circuit"
	CircuitTunnelPractice1                string = "Tunnel Practice 1"
	CircuitTunnelPractice2                string = "Tunnel Practice 2"
	CircuitTunnelPractice3                string = "Tunnel Practice 3"
	CircuitSimpleTunnel1                  string = "Simple Tunnel 1"
	CircuitSimpleTunnel2                  string = "Simple Tunnel 2"
	CircuitSimpleTunnel3                  string = "Simple Tunnel 3"
	CircuitTunnelCircuitWorld1            string = "Tunnel Circuit World 1"
	CircuitTunnelCircuitWorld2            string = "Tunnel Circuit World 2"
	CircuitTunnelCircuitWorld3            string = "Tunnel Circuit World 3"
	CircuitTunnelCircuitWorld4            string = "Tunnel Circuit World 4"
	CircuitTunnelCircuitWorld5            string = "Tunnel Circuit World 5"
	CircuitUrbanQual                      string = "Urban Qualification"
	CircuitUrbanSimple1                   string = "Urban Simple 1"
	CircuitUrbanSimple2                   string = "Urban Simple 2"
	CircuitUrbanSimple3                   string = "Urban Simple 3"
	CircuitUrbanPractice1                 string = "Urban Practice 1"
	CircuitUrbanPractice2                 string = "Urban Practice 2"
	CircuitUrbanPractice3                 string = "Urban Practice 3"
	CircuitUrbanCircuit                   string = "Urban Circuit"
	CircuitUrbanCircuitWorld1             string = "Urban Circuit World 1"
	CircuitUrbanCircuitWorld2             string = "Urban Circuit World 2"
	CircuitUrbanCircuitWorld3             string = "Urban Circuit World 3"
	CircuitUrbanCircuitWorld4             string = "Urban Circuit World 4"
	CircuitUrbanCircuitWorld5             string = "Urban Circuit World 5"
	CircuitUrbanCircuitWorld6             string = "Urban Circuit World 6"
	CircuitUrbanCircuitWorld7             string = "Urban Circuit World 7"
	CircuitUrbanCircuitWorld8             string = "Urban Circuit World 8"
	CircuitCaveSimple1                    string = "Cave Simple 1"
	CircuitCaveSimple2                    string = "Cave Simple 2"
	CircuitCaveSimple3                    string = "Cave Simple 3"
	CircuitCaveQual                       string = "Cave Qualification"
	CircuitCavePractice1                  string = "Cave Practice 1"
	CircuitCavePractice2                  string = "Cave Practice 2"
	CircuitCavePractice3                  string = "Cave Practice 3"
	CircuitCaveCircuit                    string = "Cave Circuit"
	CircuitCaveCircuitWorld1              string = "Cave Circuit World 1"
	CircuitCaveCircuitWorld2              string = "Cave Circuit World 2"
	CircuitCaveCircuitWorld3              string = "Cave Circuit World 3"
	CircuitCaveCircuitWorld4              string = "Cave Circuit World 4"
	CircuitCaveCircuitWorld5              string = "Cave Circuit World 5"
	CircuitCaveCircuitWorld6              string = "Cave Circuit World 6"
	CircuitCaveCircuitWorld7              string = "Cave Circuit World 7"
	CircuitCaveCircuitWorld8              string = "Cave Circuit World 8"
	CircuitFinalsQual                     string = "Finals Qualification"
	CircuitFinalsPractice1                string = "Finals Practice 1"
	CircuitFinalsPractice2                string = "Finals Practice 2"
	CircuitFinalsPractice3                string = "Finals Practice 3"
	CircuitVirtualStixCircuit             string = "Virtual Stix Circuit"
	CircuitVirtualStixCircuit2            string = "Virtual Stix Circuit 2"
	CircuitFinalsPreliminaryRound         string = "Finals Preliminary Round"
	CircuitFinalsPreliminaryRoundWorld1   string = "Finals Preliminary Round World 1"
	CircuitFinalsPreliminaryRoundWorld2   string = "Finals Preliminary Round World 2"
	CircuitFinalsPreliminaryRoundWorld3   string = "Finals Preliminary Round World 3"
	CircuitFinals                         string = "Final Prize Round"
	CircuitFinalsWorld1                   string = "Finals Prize Round World 1"
	CircuitFinalsWorld2                   string = "Finals Prize Round World 2"
	CircuitFinalsWorld3                   string = "Finals Prize Round World 3"
	CircuitFinalsWorld4                   string = "Finals Prize Round World 4"
	CircuitFinalsWorld5                   string = "Finals Prize Round World 5"
	CircuitFinalsWorld6                   string = "Finals Prize Round World 6"
	CircuitFinalsWorld7                   string = "Finals Prize Round World 7"
	CircuitFinalsWorld8                   string = "Finals Prize Round World 8"
	CircuitSystemsFinalsPreliminaryRound1 string = "Systems Finals Preliminary World 1"
	CircuitSystemsFinalsPreliminaryRound2 string = "Systems Finals Preliminary World 2"
	CircuitSystemsFinalsPrizeRound        string = "Systems Finals Prize Round"

	// Container names
	GazeboServerContainerName    string = "gzserver-container"
	CommsBridgeContainerName     string = "comms-bridge"
	FieldComputerContainerName   string = "field-computer"
	CopyToS3SidecarContainerName string = "copy-to-s3"
)

// subTSpecificsConfig is an internal type needed by the SubT application definition.
type subTSpecificsConfig struct {
	Region string `env:"AWS_REGION,required"`
	// MaxDurationForSimulations is the maximum number of minutes a simulation can run in SubT.
	MaxDurationForSimulations int `env:"SUBT_SIM_DURATION_MINUTES" envDefault:"60"`
	// MaxRobotModelCount is the maximum number of robots per model type. E.g. Up to 5 of any model: X1, X2, etc.
	// Robot models are defined in SubTRobotType. A value of 0 means unlimited robots.
	MaxRobotModelCount int `env:"SUBT_MAX_ROBOT_MODEL_COUNT" envDefault:"0"`
	// FuelURL contains the URL to a Fuel environment. This base URL is used to generate
	// URLs for users to access specific assets within Fuel.
	FuelURL string `env:"IGN_FUEL_URL" envDefault:"https://fuel.ignitionrobotics.org/1.0"`
}

// SubTApplication represents an application used to tailor SubT simulation requests.
type SubTApplication struct {
	cfg              subTSpecificsConfig
	schedulableTasks []SchedulableTask
}

// SimulationStatistics contains the summary values of a simulation run.
type SimulationStatistics struct {
	WasStarted          int `yaml:"was_started"`
	SimTimeDurationSec  int `yaml:"sim_time_duration_sec"`
	RealTimeDurationSec int `yaml:"real_time_duration_sec"`
	ModelCount          int `yaml:"model_count"`
}

// NewSubTApplication creates a SubT application, used to tailor simulation requests.
func NewSubTApplication(ctx context.Context) (*SubTApplication, error) {
	logger(ctx).Info("Creating new SubT application")

	s := SubTApplication{}

	s.cfg = subTSpecificsConfig{}
	// Read configuration from environment
	logger(ctx).Info("Parsing Subt config")
	if err := env.Parse(&s.cfg); err != nil {
		return nil, err
	}

	// Populate the robot config list
	loadSubTRobotTypes(&s.cfg)

	return &s, nil
}

func (sa *SubTApplication) getApplicationName() string {
	return applicationSubT
}

func (sa *SubTApplication) getPlatformName() string {
	return platformSubT
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getGazeboPodName returns the name of the Gazebo pod for a simulation.
func (sa *SubTApplication) getGazeboPodName(podNamePrefix string) string {
	return fmt.Sprintf("%s-gzserver", podNamePrefix)
}

// getCommsBridgePodName returns the name of the comms bridge pod for a specific field-computer in a simulation.
func (sa *SubTApplication) getCommsBridgePodName(podNamePrefix string, robotIdentifier string) string {
	return fmt.Sprintf("%s-comms-%s", podNamePrefix, robotIdentifier)
}

// getFieldComputerPodName returns the name of a specific field-computer in a simulation.
func (sa *SubTApplication) getFieldComputerPodName(podNamePrefix string, robotIdentifier string) string {
	return fmt.Sprintf("%s-fc-%s", podNamePrefix, robotIdentifier)
}

func (sa *SubTApplication) getCopyPodName(targetPodName string) string {
	return fmt.Sprintf("%s-copy", targetPodName)
}

// getGazeboPodName returns the name of the Gazebo pod for a simulation.
func (sa *SubTApplication) getWebsocketServiceName(podNamePrefix string) string {
	return fmt.Sprintf("%s-websocket", podNamePrefix)
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getSimulationIngressPath returns the request path for the websocket server of a simulation.
func (sa *SubTApplication) getSimulationIngressPath(groupID string) string {
	return fmt.Sprintf("/simulations/%s", groupID)
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getSimulationSummaryFilename returns the filename of a simulation summary.
func (sa *SubTApplication) getSimulationSummaryFilename(groupID string) string {
	return fmt.Sprintf("%s-summary.json", groupID)
}

// getGazeboLogsFilename returns the filename of the Gazebo logs for a specific
// simulation.
func (sa *SubTApplication) getGazeboLogsFilename(groupID string) string {
	return fmt.Sprintf("%s.tar.gz", groupID)
}

// getRobotROSLogsFilename returns the filename of the ROS logs for a specific
// robot in a simulation.
func (sa *SubTApplication) getRobotROSLogsFilename(groupID string, robotName string) string {
	return fmt.Sprintf("%s-fc-%s-commsbridge.tar.gz", groupID, strings.ToLower(robotName))
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getRobotIdentifierFromNames returns the robot identifier of a robot in a simulation.
// You should call this function when you have a simulation deployment, since it has a property called robotNames.
func (sa *SubTApplication) getRobotIdentifierFromNames(robotNames *string, robotName string) (*string, *ign.ErrMsg) {
	robots := strings.Split(*robotNames, ",")
	for i, rn := range robots {
		if rn == robotName {
			identifier := sptr(fmt.Sprintf("rbt%d", i+1))
			return identifier, nil
		}
	}
	err := NewErrorMessage(ErrorRobotIdentifierNotFound)
	return nil, err
}

// getRobotIdentifierFromList returns the robot identifier of a robot in a simulation.
// You should call this function when you don't have a simulation deployment yet, but you have the actual list of robots.
func (sa *SubTApplication) getRobotIdentifierFromList(robotList []SubTRobot, robotName string) (*string, *ign.ErrMsg) {
	for i, robot := range robotList {
		if robot.Name == robotName {
			identifier := sptr(fmt.Sprintf("rbt%d", i+1))
			return identifier, nil
		}
	}
	err := NewErrorMessage(ErrorRobotIdentifierNotFound)
	return nil, err
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// customizeSimulationRequest performs operations to a simulation request in order to be
// executed by SubT application.
func (sa *SubTApplication) customizeSimulationRequest(ctx context.Context,
	s *Service, r *http.Request, tx *gorm.DB, createSim *CreateSimulation, username string) *ign.ErrMsg {
	var subtSim SubTCreateSimulation
	var rules *SubTCircuitRules
	var creditsSum int
	var err error

	if em := ParseStruct(&subtSim, r, true); em != nil {
		return em
	}
	// Create the slice of robots and then serialize into the createSim's "Extra" field
	robots := make([]SubTRobot, 0)
	robotNames := make([]string, 0)
	robotModelCount := make(map[string]int, 0)
	for i, rn := range subtSim.RobotName {
		robotType := SubTRobotTypes[subtSim.RobotType[i]]

		robot := SubTRobot{
			Name:    rn,
			Type:    subtSim.RobotType[i],
			Image:   subtSim.RobotImage[i],
			Credits: robotType.Credits,
		}
		creditsSum += robot.Credits

		// Limit the number of robots per model if a limit is set
		if sa.cfg.MaxRobotModelCount > 0 {
			robotModelCount[robotType.Model]++
			if robotModelCount[robotType.Model] > sa.cfg.MaxRobotModelCount {
				msg := fmt.Sprintf("too many robots of model %s", robotType.Model)
				return NewErrorMessageWithBase(ErrorRobotModelLimitReached, errors.New(msg))
			}
		}

		robots = append(robots, robot)
		robotNames = append(robotNames, robot.Name)
	}

	marsupials := make([]SubTMarsupial, 0)
	// Process the marsupial parameters.
	for _, mar := range subtSim.Marsupial {
		// A marsupial pair is specified as a string of the form "<parent>:<child>"
		parts := strings.Split(mar, ":")

		// Make sure there is both a parent and a child.
		if len(parts) != 2 {
			return NewErrorMessageWithBase(ErrorInvalidMarsupialSpecification, err)
		}

		// Try to find the parent and child in the set of robots.
		var foundParent = false
		var foundChild = false
		for _, robot := range robots {
			if robot.Name == parts[0] {
				foundParent = true
			}
			if robot.Name == parts[1] {
				foundChild = true
			}
		}
		// Make sure both the parent and child were found.
		if !foundParent || !foundChild {
			return NewErrorMessageWithBase(ErrorInvalidMarsupialSpecification, err)
		}

		marsupial := SubTMarsupial{
			Parent: parts[0],
			Child:  parts[1],
		}
		marsupials = append(marsupials, marsupial)
	}

	rules, err = GetCircuitRules(tx, subtSim.Circuit)
	if err != nil {
		return NewErrorMessageWithBase(ErrorCircuitRuleNotFound, err)
	}

	// Perform some additional checks if the user is not a system admin
	if !s.userAccessor.IsSystemAdmin(username) {
		if rules.MaxCredits != nil {
			if creditsSum > *rules.MaxCredits {
				return NewErrorMessage(ErrorCreditsExceeded)
			}
		}

		if !subtSim.robotImagesBelongToECROwner() {
			return NewErrorMessage(ErrorInvalidRobotImage)
		}

		if !sa.isQualified(subtSim.Owner, subtSim.Circuit, username) {
			return NewErrorMessage(ErrorNotQualified)
		}

		if isSubmissionDeadlineReached(*rules) {
			return NewErrorMessage(ErrorSubmissionDeadlineReached)
		}
	}

	extra := &ExtraInfoSubT{
		Circuit:    subtSim.Circuit,
		Robots:     robots,
		Marsupials: marsupials,
	}
	createSim.ExtraSelector = &subtSim.Circuit

	if createSim.Extra, err = extra.ToJSON(); err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorMarshalJSON, err)
	}
	// TODO: Robots are temporarily placed in the create simulation object.
	//  Ideally this should be accessed directly from the extra params or
	//  wherever that information is being stored (another table, etc.).
	createSim.Robots = sptr(strings.Join(robotNames, ","))

	return nil
}

// isQualified returns a boolen if the owner is qualified to participate in the circuit.
// If the username matches with a system admin, this function will return true as well.
func (sa *SubTApplication) isQualified(owner, circuit string, username string) bool {
	return IsOwnerQualifiedForCircuit(globals.Server.Db, owner, circuit, username)
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// spawnChildSimulationDeployments. By default, we launch a single simulation from a CreateSimulation request.
// But we allow specific ApplicationTypes (eg. SubT) to spawn multiple simulations
// from a single request. When that happens, we call those "child simulations"
// and they will be grouped by the same parent simulation's groupID.
func (sa *SubTApplication) spawnChildSimulationDeployments(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) ([]*SimulationDeployment, *ign.ErrMsg) {

	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	result := []*SimulationDeployment{}

	// Determine if the simulation is a multisim.
	// Multisims will run a set of worlds and multiple times each of them.
	// The set of worlds and times to run each of them will be stored in the DB so they
	// are not known in advance by SubT participant teams.
	rules, err := GetCircuitRules(tx, extra.Circuit)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}
	times, err := StrToIntSlice(*rules.Times)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}
	multisim := len(times) > 1 || times[0] > 1

	// If the simulation is a multisim, then create and return the set of child simulations
	if multisim {
		// Create the child simulations
		worlds := ign.StrToSlice(*rules.Worlds)
		var worldIdx int
		var childIdx int
		for worldIdx = range worlds {
			for j := 0; j < times[worldIdx]; j++ {
				childIdx++
				childSim := dep.Clone()
				childSim.GroupID = sptr(fmt.Sprintf("%s-c-%d", *dep.GroupID, childIdx))

				// Set new auth token to authorize external services
				token, err := generateToken(nil)
				if err != nil {
					return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
				}
				childSim.AuthorizationToken = &token

				// Create a clone of the parent's extra info and set it to the child sim.
				newExtra := *extra
				newExtra.WorldIndex = &worldIdx
				newExtra.RunIndex = intptr(childIdx - 1)
				if childSim.Extra, err = newExtra.ToJSON(); err != nil {
					return nil, ign.NewErrorMessageWithBase(ign.ErrorMarshalJSON, err)
				}
				result = append(result, childSim)
			}
		}

		return result, nil
	}
	// This is a normal single simulation. No child simulations.
	return nil, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// GetSchedulableTasks returns a slice of schedulable tasks to be registered on cloudsim's scheduler.
func (sa *SubTApplication) GetSchedulableTasks(ctx context.Context, s *Service, tx *gorm.DB) []SchedulableTask {
	sa.schedulableTasks = append(sa.schedulableTasks, sa.launchCircuitOnCompetitionDay(ctx, s, tx)...)
	return sa.schedulableTasks
}

// launchCircuitOnCompetitionDay returns an array of schedulable tasks to deploy a set of simulations on a competition day.
// It creates a set of tasks from any circuit that has a date assigned to its SubT rule.
func (sa *SubTApplication) launchCircuitOnCompetitionDay(ctx context.Context, s *Service, tx *gorm.DB) []SchedulableTask {
	rules, err := GetPendingCircuitRules(tx)
	if err != nil {
		return []SchedulableTask{}
	}

	var tasks []SchedulableTask
	for _, rule := range *rules {
		if rule.CompetitionDate == nil {
			continue
		}

		logger(ctx).Info(
			fmt.Sprintf("Scheduling [%s] simulations to run on [%s (%f seconds)].",
				*rule.Circuit,
				*rule.CompetitionDate,
				(*rule.CompetitionDate).Sub(time.Now()).Seconds(),
			),
		)

		// The circuit value needs to be stored in a variable to avoid sharing
		// it between closures
		circuit := *rule.Circuit
		task := SchedulableTask{
			Fn: func() {
				logger(ctx).Info(fmt.Sprintf("Launching scheduled simulations for [%s].", circuit))
				s.DeployHeldCircuitSimulations(ctx, tx, circuit)
			},
			Date: *rule.CompetitionDate,
		}

		tasks = append(tasks, task)
	}

	return tasks
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// checkCanShutdownSimulation allows specific applications to decide if a given user
// can shutdown a simulation.
func (sa *SubTApplication) checkCanShutdownSimulation(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment, user *users.User) (bool, *ign.ErrMsg) {
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return false, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Members of 'subt' Organization (ie. Competition Admins) are the only ones
	// that can shutdown competition simulations.
	circuits := []string{
		CircuitTunnelCircuit,
		CircuitUrbanCircuit,
	}
	if StrSliceContains(extra.Circuit, circuits) {
		return s.userAccessor.CanPerformWithRole(sptr(applicationSubT), *user.Username, per.Member)
	}

	return true, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// checkValidNumberOfSimulations checks if the given owner hasn't gone beyond the
// maximum number of allowed simulations for the circuit.
func (sa *SubTApplication) checkValidNumberOfSimulations(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {

	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Get the number of remaining submissions allowed
	remaining, err := getRemainingSubmissions(tx, extra.Circuit, *dep.Owner)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Dev note: we check 'after' creating the record in the DB to make
	// sure that in case of a race condition then both records are added with pending state
	// and one of those (or both) can be rejected immediately.
	// So, we need to check for "< 0"
	if remaining != nil && *remaining < 0 {
		errStr := fmt.Sprintf("Subt - The Owner [%s] has reached the Max simulations submission limit for Circuit [%s]", *dep.Owner, extra.Circuit)
		logger(ctx).Info(errStr)
		newErr := errors.New(errStr)
		return NewErrorMessageWithBase(ErrorCircuitSubmissionLimitReached, newErr)
	}

	// All OK
	return nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getGazeboWorldStatsTopicAndLimit returns the topic to subscribe to get notifications about the simulation
// state (eg. /world/default/stats) and time, as well as the limit of simulation seconds, if any.
func (sa *SubTApplication) getGazeboWorldStatsTopicAndLimit(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (string, int, error) {
	// Parse the SubT extra info required for this Simulation
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return "", 0, err
	}
	// Now get the Circuit rules for this simulation
	rules, err := GetCircuitRules(tx, extra.Circuit)
	if err != nil {
		return "", 0, err
	}

	// Read the list of stats topics (for each world)
	worldStatsTopics := ign.StrToSlice(*rules.WorldStatsTopics)
	statsTopic := worldStatsTopics[0]
	if len(worldStatsTopics) > 1 {
		statsTopic = worldStatsTopics[*extra.WorldIndex]
	}

	// Read the list of "max simulation seconds" for each Circuit world
	simSecondsList, err := StrToIntSlice(*rules.WorldMaxSimSeconds)
	if err != nil {
		return "", 0, err
	}
	maxSeconds := simSecondsList[0]
	if len(simSecondsList) > 1 {
		maxSeconds = simSecondsList[*extra.WorldIndex]
	}
	// Simulations will automatically be shutdown by Gazebo once their time expires.
	// In order to prevent a simulation from running forever in case Gazebo is
	// not able to shut it down, the maximum simulation duration is limited to
	// 1.5 times the simulation time.
	maxSeconds = maxSeconds * 3 / 2

	return statsTopic, maxSeconds, nil
}

// getGazeboWorldWarmupTopic returns the topic to subscribe to get notifications about the simulation
// being ready to start and finish.
func (sa *SubTApplication) getGazeboWorldWarmupTopic(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (string, error) {
	// Parse the SubT extra info required for this Simulation
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return "", err
	}
	// Now get the Circuit rules for this simulation
	rules, err := GetCircuitRules(tx, extra.Circuit)
	if err != nil {
		return "", err
	}

	if rules.WorldWarmupTopics == nil {
		// No topics configured. Just return nil
		return "", nil
	}

	// Read the list of warmup topics (for each world)
	worldWarmupTopics := ign.StrToSlice(*rules.WorldWarmupTopics)

	warmupTopic := worldWarmupTopics[0]
	if len(worldWarmupTopics) > 1 {
		warmupTopic = worldWarmupTopics[*extra.WorldIndex]
	}

	return warmupTopic, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getSimulationWebsocketAddress returns a simulation's websocket server address as well a the authorization token.
func (sa *SubTApplication) getSimulationWebsocketAddress(ctx context.Context, s *Service, tx *gorm.DB,
	store store.Store, dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	// The simulation must be running to be able to connect to the websocket server
	if *dep.DeploymentStatus != int(simRunning) {
		return nil, ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus)
	}

	host := store.Orchestrator().IngressHost()
	path := store.Ignition().GetWebsocketPath(dep.GetGroupID())

	return &WebsocketAddressResponse{
		Token:   *dep.AuthorizationToken,
		Address: fmt.Sprintf("%s/%s", host, path),
	}, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getSimulationLogsForDownload returns a link to the GZ logs that were saved in S3.
func (sa *SubTApplication) getSimulationLogsForDownload(ctx context.Context, tx *gorm.DB, p platform.Platform,
	dep *SimulationDeployment, robotName *string) (*string, *ign.ErrMsg) {

	// In SubT, we return a summary generated from all children simulations for
	// multi-sims. For single sims, we should return ROS logs for a specific
	// robot if a robot name is specified or the complete Gazebo logs otherwise.
	var fileName string
	if dep.isMultiSim() {
		fileName = sa.getSimulationSummaryFilename(*dep.GroupID)
	} else if robotName != nil {
		fileName = sa.getRobotROSLogsFilename(*dep.GroupID, *robotName)
	} else {
		fileName = sa.getGazeboLogsFilename(*dep.GroupID)
	}

	bucket := p.Store().Ignition().LogsBucket()
	ownerNameEscaped := url.PathEscape(*dep.Owner)
	folderPath := fmt.Sprintf("/gz-logs/%s/%s/", ownerNameEscaped, *dep.GroupID)
	filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
	logger(ctx).Debug(fmt.Sprintf("SubT App - Fetching generating link to fetch logs from S3 bucket [%s] with path [%s]\n", bucket, filePath))

	url, err := p.Storage().GetURL(bucket, filePath, 5*time.Minute)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	return &url, nil
}

// getSimulationsLiveLogs returns logs from a running simulation.
// In case the simulation is a multisim, this method will return an *AggregatedSubTSimulationValues
// If the simulation a single simulation, it will return a PodLog directly from Kubernetes.
func (sa *SubTApplication) getSimulationLiveLogs(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment, robotName *string, lines int64) (interface{}, *ign.ErrMsg) {

	// This block covers the summary case indicated in the documentation above.
	if dep.isMultiSim() {
		summary, err := GetAggregatedSubTSimulationValues(tx, dep)
		if err != nil {
			return nil, NewErrorMessageWithBase(ErrorFailedToGetLiveLogs, err)
		}
		return summary, nil
	}

	var podName string
	var container string
	podPrefix := getSimulationPodNamePrefix(*dep.GroupID)

	if robotName != nil {
		identifier, err := sa.getRobotIdentifierFromNames(dep.Robots, *robotName)
		if err != nil {
			return nil, err
		}
		podName = sa.getCommsBridgePodName(podPrefix, *identifier)
		container = CommsBridgeContainerName
	} else {
		podName = sa.getGazeboPodName(podPrefix)
		container = GazeboServerContainerName
	}

	// Get platform
	p, err := platformManager.GetSimulationPlatform(s.platforms, dep)
	if err != nil {
		return nil, NewErrorMessageWithBase(ErrorFailedToGetLiveLogs, err)
	}

	// Get logs
	res := resource.NewResource(
		podName,
		p.Store().Orchestrator().Namespace(),
		nil,
	)
	reader := p.Orchestrator().Pods().Reader(res)
	log, err := reader.Logs(container, lines)
	if err != nil {
		return nil, NewErrorMessageWithBase(ErrorFailedToGetLiveLogs, err)
	}

	podLog := PodLog{log}

	return &podLog, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getIngress retrieves an Ingress from the cluster.
func getIngress(ctx context.Context, kc kubernetes.Interface, namespace string,
	ingressName string) (*v1beta1.Ingress, error) {
	return kc.ExtensionsV1beta1().Ingresses(namespace).Get(ingressName, metav1.GetOptions{})
}

// updateIngress updates an Ingress resource in the cluster.
func updateIngress(ctx context.Context, kc kubernetes.Interface, namespace string,
	ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {

	ingress, err := kc.ExtensionsV1beta1().Ingresses(namespace).Update(ingress)
	if err != nil {
		errMsg := "failed to update ingress"
		logger(ctx).Error(errMsg, err)
		return nil, errors.New(errMsg)
	}

	return ingress, nil
}

// getIngressRule gets a host's rule from an Ingress resource.
// The `host` parameter is used to select the rule from which to remove paths.
// If there is more than one rule for a host, only the first rule will be returned.
// If `host` is nil, the first rule with an empty host field will be returned.
func getIngressRule(ctx context.Context, ingress *v1beta1.Ingress,
	host *string) (*v1beta1.HTTPIngressRuleValue, error) {
	// Set host default value if nil
	noHost := sptr("")
	if host == nil {
		host = noHost
	}

	// Find the target rule
	var rule *v1beta1.HTTPIngressRuleValue
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.Host == *host {
			rule = ingressRule.IngressRuleValue.HTTP
			return rule, nil
		}
	}
	// If no rule was found return an error
	if host == noHost {
		host = sptr("nil")
	}
	return nil, fmt.Errorf("ingress rule for host %s was not found", *host)
}

// upsertIngressRule inserts or updates a set of paths into an Ingress rule.
// The `host` parameter is used to select the rule from which to remove paths. If there is more than one rule for a
// host, only the first rule will be modified. If `host` is nil, a rule with no host will be modified.
func upsertIngressRule(ctx context.Context, kc kubernetes.Interface, namespace string, ingressName string,
	host *string, paths ...*v1beta1.HTTPIngressPath) (*v1beta1.Ingress, error) {

	// Get the ingress from the cluster
	ingress, err := getIngress(ctx, kc, namespace, ingressName)
	if err != nil {
		return nil, err
	}

	// Extract the host rule from the ingress resource
	rule, err := getIngressRule(ctx, ingress, host)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		// Try to find and update the path
		updated := false
		for i, rulePath := range rule.Paths {
			if rulePath.Path == path.Path {
				updated = true
				if path != nil {
					rule.Paths[i] = *path
				}
				break
			}
		}
		// No path was updated, create a new one
		if !updated && path != nil {
			rule.Paths = append(rule.Paths, *path)
		}
	}

	// Apply updated rule
	return updateIngress(ctx, kc, namespace, ingress)
}

// removeIngressRule removes a set of paths from an Ingress rule.
// Note that the Kubernetes spec requires ingress rules to have at least one path. Attempting to remove a rule's only
// path will fail.
// The `host` parameter is used to select the rule from which to remove paths. If there is more than one rule for a
// host, only the first rule will be modified. If `host` is nil, a rule with no host will be modified.
func removeIngressRule(ctx context.Context, kc kubernetes.Interface, namespace string,
	ingressName string, host *string, paths ...string) (*v1beta1.Ingress, error) {

	// Get the ingress from the cluster
	ingress, err := getIngress(ctx, kc, namespace, ingressName)
	if err != nil {
		return nil, err
	}

	// Extract the host rule from the ingress resource
	rule, err := getIngressRule(ctx, ingress, host)
	if err != nil {
		return nil, err
	}

	// Remove paths
	for _, path := range paths {
		for i, rulePath := range rule.Paths {
			if rulePath.Path == path {
				pathsLen := len(rule.Paths)
				if pathsLen > 1 {
					rule.Paths[i] = rule.Paths[pathsLen-1]
				}
				rule.Paths = rule.Paths[:pathsLen-1]
				break
			}
		}
	}

	// Apply updated rule
	return updateIngress(ctx, kc, namespace, ingress)
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// processScore creates the competition_scores entry for a simulation and updates its score.
func processScore(ctx context.Context, userAccessor useracc.Service, tx *gorm.DB,
	dep *SimulationDeployment, score *float64) *ign.ErrMsg {
	logger(ctx).Info(
		fmt.Sprintf("processScore - Creating competition_scores entry for simulation [%s]", *dep.GroupID),
	)
	if em := userAccessor.AddScore(dep.GroupID, dep.Application, dep.ExtraSelector, dep.Owner,
		score, dep.GroupID); em != nil {
		logMsg := fmt.Sprintf(
			"processScore - Could not create competition_scores entry for simulation [%s].",
			*dep.GroupID,
		)
		logger(ctx).Error(logMsg, em)
		return em
	}
	if em := dep.UpdateScore(tx, score); em != nil {
		logMsg := fmt.Sprintf(
			"processScore - Could not create competition_scores entry for simulation [%s].",
			*dep.GroupID,
		)
		logger(ctx).Error(logMsg, em)
		return em
	}
	return nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// getMaxDurationForSimulation returns the max allowed duration for a simulation,
// before marking it for shutdown.
func (sa *SubTApplication) getMaxDurationForSimulation(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) time.Duration {
	return time.Duration(sa.cfg.MaxDurationForSimulations) * time.Minute
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// uploadToS3LogBucket uploads a file to a simulation log folder.
func (sa *SubTApplication) uploadToS3SimulationLogBucket(p platform.Platform, dep *SimulationDeployment,
	filename string, file []byte) *ign.ErrMsg {

	input := storage.UploadInput{
		Bucket:        p.Store().Ignition().LogsBucket(),
		Key:           path.Join(GetS3SimulationLogKey(dep), filename),
		File:          bytes.NewReader(file),
		ContentLength: int64(len(file)),
		ContentType:   http.DetectContentType(file),
	}
	if err := p.Storage().Upload(input); err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	return nil
}

// uploadSimulationSummary uploads the simulation summary to the log bucket of the simulation.
// Teams can later download this summary through a generated link.
func (sa *SubTApplication) uploadSimulationSummary(p platform.Platform, simDep *SimulationDeployment,
	summary *AggregatedSubTSimulationValues) *ign.ErrMsg {

	// Prepare data
	values, err := json.Marshal(summary)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorMarshalJSON, err)
	}

	// Upload file
	fileName := sa.getSimulationSummaryFilename(*simDep.GroupID)
	if em := sa.uploadToS3SimulationLogBucket(p, simDep, fileName, values); em != nil {
		return em
	}

	return nil
}

// updateMultiSimStatuses updates the status of a Multi-Sim Parent and executes application-specific logic based on the
// state of its children.
func (sa *SubTApplication) updateMultiSimStatuses(ctx context.Context, tx *gorm.DB, userAccessor useracc.Service,
	p platform.Platform, simDep *SimulationDeployment) *ign.ErrMsg {
	// Note: simDep is a Parent in a multi-sim

	// Only proceed if the simulation terminated successfully. Get the aggregated values from all children
	if simDep.IsRunning() || simDep.ErrorStatus != nil || simDep.Processed || simDep.Held {
		return nil
	}

	// Get the score for the simulation. Parent simulation scores are based on the performance of its children.
	summary, err := GetAggregatedSubTSimulationValues(tx, simDep)
	if err != nil {
		logger(ctx).Error("Error computing aggregated values for simulation: "+*simDep.GroupID, err)
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Create the score entry
	if !globals.DisableScoreGeneration {
		if em := processScore(ctx, userAccessor, tx, simDep, &summary.Score); em != nil {
			return em
		}
	}

	// Create and upload the parent summary to S3
	if p.Store().Ignition().LogsCopyEnabled() {
		logger(ctx).Info(
			fmt.Sprintf("updateMultiSimStatuses - Uploading simulation summary for simulation [%s]", *simDep.GroupID),
		)
		if em := sa.uploadSimulationSummary(p, simDep, summary); em != nil {
			logMsg := fmt.Sprintf(
				"updateMultiSimStatuses - Could not upload simulation summary to S3 for simulation [%s].",
				*simDep.GroupID,
			)
			logger(ctx).Error(logMsg, em)
			return em
		}
	}

	// Send an email with the summary to the competitor
	if !globals.DisableSummaryEmails {
		// TODO Use the simulation platform for this
		SendSimulationSummaryEmail(p.EmailSender(), simDep, *summary, nil)
	}

	if !simDep.Processed {
		if err := simDep.UpdateProcessed(tx, true); err != nil {
			return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
		}
	}

	return nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// invalidateSimulation is invoked by the sim_service when a simulation is about
// to be restarted. The old simulation run should be invalidated. This usually
// means soft-deleting some data from DB.
func (sa *SubTApplication) invalidateSimulation(ctx context.Context, tx *gorm.DB,
	simDep *SimulationDeployment) error {

	// we just soft delete the SimulationDeploymentsSubTValue corresponding to the
	// given simulation
	if err := tx.Where("group_id = ?", *simDep.GroupID).
		Delete(&SimulationDeploymentsSubTValue{}).Error; err != nil {
		return err
	}
	return nil
}

// getCompetitionRobots returns the list of available robots configurations for SubT circuits.
func (sa *SubTApplication) getCompetitionRobots() (interface{}, *ign.ErrMsg) {
	return SubTRobotTypes, nil
}

// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////

// ValidateSimulationLaunch returns an error if there is an error on the validation process.
// Any function that needs to validate a simulation right before pushing to the queue should be appended here.
func (sa *SubTApplication) ValidateSimulationLaunch(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {
	if err := sa.checkHeldSimulation(ctx, tx, dep); err != nil {
		return err
	}

	if err := sa.checkSupersededSimulation(ctx, *dep.GroupID, *dep.DeploymentStatus); err != nil {
		return err
	}

	return nil
}

// checkHeldSimulation is a validator that returns an error if the simulations is being held.
// It should be used inside the ValidateSimulationLaunch before pushing a simulation to the queue.
func (sa *SubTApplication) checkHeldSimulation(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {
	if dep.Held {
		logger(ctx).Warning(fmt.Sprintf("checkHeldSimulation - Cannot run a held simulation (Group ID: %s)", *dep.GroupID))
		return NewErrorMessage(ErrorLaunchHeldSimulation)
	}
	return nil
}

// checkSupersededSimulation is a validator that returns an error if the given status equals to superseded.
func (sa *SubTApplication) checkSupersededSimulation(ctx context.Context, groupID string, status int) *ign.ErrMsg {
	if simSuperseded.Eq(status) {
		logger(ctx).Warning(fmt.Sprintf("checkSupersededSimulation - Cannot run a Superseded simulation (Group ID: %s)", groupID))
		return NewErrorMessage(ErrorLaunchSupersededSimulation)
	}
	return nil
}

// simulationIsHeld returns true if the simulation needs to be held. In any other cases, it returns false.
// It checks if the simulation is part of a certain circuit that has not reached its competition day yet.
func (sa *SubTApplication) simulationIsHeld(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) bool {
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		logger(ctx).Warning(fmt.Sprintf("simulationIsHeld - Cannot read extra field from simulation %s", *dep.GroupID))
		return false
	}

	rules, err := GetCircuitRules(tx, extra.Circuit)
	if err != nil {
		logger(ctx).Warning(fmt.Sprintf("simulationIsHeld - Cannot get rules for the circuit %s", extra.Circuit))
		return false
	}

	if rules.CompetitionDate == nil {
		logger(ctx).Debug(fmt.Sprintf("simulationIsHeld - Competition date for %s is null", extra.Circuit))
		return false
	}

	delta := time.Now().Sub(*rules.CompetitionDate).Seconds()
	if delta >= -1 {
		return false
	}
	return true
}

// isSubmissionDeadlineReached checks if a certain circuit has reached its submission deadline.
// It only returns true if the deadline is set and has been reached, in any other case it returns false.
func isSubmissionDeadlineReached(circuit SubTCircuitRules) bool {
	return circuit.SubmissionDeadline != nil && circuit.SubmissionDeadline.Before(time.Now())
}

// IsCompetitionCircuit checks if the given circuit is a competition circuit.
// This is used to check if the given circuit is a Tunnel, Urban or Cave circuit.
func IsCompetitionCircuit(circuit string) bool {
	for _, c := range SubTCompetitionCircuits {
		if c == circuit {
			return true
		}
	}
	return false
}
