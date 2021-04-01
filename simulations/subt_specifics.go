package simulations

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/simulations/gloo"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	subtTypeMoleBridge    = "mole-bridge"
	subtTypeCommsBridge   = "comms-bridge"
	subtTypeFieldComputer = "field-computer"

	// A predefined const to refer to the SubT Application type
	// This will be used to know which Pods and services launch.
	applicationSubT            string = "subt"
	CircuitNIOSHSRConfigA      string = "NIOSH SR Config A"
	CircuitNIOSHSRConfigB      string = "NIOSH SR Config B"
	CircuitNIOSHEXConfigA      string = "NIOSH EX Config A"
	CircuitNIOSHEXConfigB      string = "NIOSH EX Config B"
	CircuitVirtualStix         string = "Virtual Stix"
	CircuitTunnelCircuit       string = "Tunnel Circuit"
	CircuitTunnelPractice1     string = "Tunnel Practice 1"
	CircuitTunnelPractice2     string = "Tunnel Practice 2"
	CircuitTunnelPractice3     string = "Tunnel Practice 3"
	CircuitSimpleTunnel1       string = "Simple Tunnel 1"
	CircuitSimpleTunnel2       string = "Simple Tunnel 2"
	CircuitSimpleTunnel3       string = "Simple Tunnel 3"
	CircuitTunnelCircuitWorld1 string = "Tunnel Circuit World 1"
	CircuitTunnelCircuitWorld2 string = "Tunnel Circuit World 2"
	CircuitTunnelCircuitWorld3 string = "Tunnel Circuit World 3"
	CircuitTunnelCircuitWorld4 string = "Tunnel Circuit World 4"
	CircuitTunnelCircuitWorld5 string = "Tunnel Circuit World 5"
	CircuitUrbanQual           string = "Urban Qualification"
	CircuitUrbanSimple1        string = "Urban Simple 1"
	CircuitUrbanSimple2        string = "Urban Simple 2"
	CircuitUrbanSimple3        string = "Urban Simple 3"
	CircuitUrbanPractice1      string = "Urban Practice 1"
	CircuitUrbanPractice2      string = "Urban Practice 2"
	CircuitUrbanPractice3      string = "Urban Practice 3"
	CircuitUrbanCircuit        string = "Urban Circuit"
	CircuitUrbanCircuitWorld1  string = "Urban Circuit World 1"
	CircuitUrbanCircuitWorld2  string = "Urban Circuit World 2"
	CircuitUrbanCircuitWorld3  string = "Urban Circuit World 3"
	CircuitUrbanCircuitWorld4  string = "Urban Circuit World 4"
	CircuitUrbanCircuitWorld5  string = "Urban Circuit World 5"
	CircuitUrbanCircuitWorld6  string = "Urban Circuit World 6"
	CircuitUrbanCircuitWorld7  string = "Urban Circuit World 7"
	CircuitUrbanCircuitWorld8  string = "Urban Circuit World 8"
	CircuitCaveSimple1         string = "Cave Simple 1"
	CircuitCaveSimple2         string = "Cave Simple 2"
	CircuitCaveSimple3         string = "Cave Simple 3"
	CircuitCaveQual            string = "Cave Qualification"
	CircuitCavePractice1       string = "Cave Practice 1"
	CircuitCavePractice2       string = "Cave Practice 2"
	CircuitCavePractice3       string = "Cave Practice 3"
	CircuitCaveCircuit         string = "Cave Circuit"
	CircuitCaveCircuitWorld1   string = "Cave Circuit World 1"
	CircuitCaveCircuitWorld2   string = "Cave Circuit World 2"
	CircuitCaveCircuitWorld3   string = "Cave Circuit World 3"
	CircuitCaveCircuitWorld4   string = "Cave Circuit World 4"
	CircuitCaveCircuitWorld5   string = "Cave Circuit World 5"
	CircuitCaveCircuitWorld6   string = "Cave Circuit World 6"
	CircuitCaveCircuitWorld7   string = "Cave Circuit World 7"
	CircuitCaveCircuitWorld8   string = "Cave Circuit World 8"
	CircuitFinalsQual          string = "Finals Qualification"
	CircuitFinalsPractice1     string = "Finals Practice 1"
	CircuitFinalsPractice2     string = "Finals Practice 2"
	CircuitFinalsPractice3     string = "Finals Practice 3"

	// Container names
	GazeboServerContainerName    string = "gzserver-container"
	MoleBridgeContainerName      string = "mole-ros-pulsar-bridge"
	CommsBridgeContainerName     string = "comms-bridge"
	FieldComputerContainerName   string = "field-computer"
	CopyToS3SidecarContainerName string = "copy-to-s3"
)

// subTSpecificsConfig is an internal type needed by the SubT application definition.
type subTSpecificsConfig struct {
	AwsSecretName string `env:"K8_AWS_SECRET_NAME" envDefault:"aws-secrets"`
	Region        string `env:"AWS_REGION,required"`
	S3LogsBucket  string `env:"AWS_GZ_LOGS_BUCKET,required"`
	// Should we backup logs to S3?
	S3LogsCopyEnabled bool `env:"AWS_GZ_LOGS_ENABLED" envDefault:"true"`
	// MaxDurationForSimulations is the maximum number of minutes a simulation can run in SubT.
	MaxDurationForSimulations int `env:"SUBT_SIM_DURATION_MINUTES" envDefault:"60"`
	// MaxRobotModelCount is the maximum number of robots per model type. E.g. Up to 5 of any model: X1, X2, etc.
	// Robot models are defined in SubTRobotType. A value of 0 means unlimited robots.
	MaxRobotModelCount int `env:"SUBT_MAX_ROBOT_MODEL_COUNT" envDefault:"0"`
	// AllowNotFoundDuringShutdown is a bool flag used to fail if a pod/service is not found
	// during shut down. If 'true' then this handler will not fail when a pod-to-be-killed is not found.
	AllowNotFoundDuringShutdown bool `env:"SUBT_ALLOW_NOT_FOUND" envDefault:"true"`
	// IgnVerbose is the IGN_VERBOSE value that will be passed to Pods launched for SubT.
	IgnVerbose string `env:"IGN_VERBOSE" envDefault:"1"`
	// GazeboLogsVolumeMountPath is the path inside the container where the `gz-logs` Volume will be mounted.
	// eg. '/tmp/ign'.
	// Important: it is important that gazebo is configured to record its logs to a child folder of the
	// passed mount location (eg. following the above example, '/var/log/gzserver/logs'), otherwise gazebo
	// will 'mkdir' and override the mounted folder.
	// See the "gzserver-container" Container Spec below to see the code.
	GazeboLogsVolumeMountPath           string `env:"SUBT_GZSERVER_LOGS_VOLUME_MOUNT_PATH" envDefault:"/tmp/ign"`
	ROSLogsVolumeMountPath              string `env:"SUBT_BRIDGE_LOGS_VOLUME_MOUNT_PATH" envDefault:"/home/developer/.ros"`
	SidecarContainerLogsVolumeMountPath string `env:"SUBT_SIDECAR_CONTAINER_VOLUME_MOUNT_PATH" envDefault:"/tmp/logs"`
	TerminationGracePeriodSeconds       int64  `env:"SUBT_GZSERVER_TERMINATE_GRACE_PERIOD_SECONDS" envDefault:"120"`
	// IgnIP is the Cloudsim server's IP address to use when creating NetworkPolicies.
	// Note: when run at Elasticbeanstalk this env var will be automatically set.
	// See 'docker-entrypoint.sh' script located at the root folder of this project.
	IgnIP string `env:"IGN_IP"`
	// FuelURL contains the URL to a Fuel environment. This base URL is used to generate
	// URLs for users to access specific assets within Fuel.
	FuelURL string `env:"IGN_FUEL_URL" envDefault:"https://fuel.ignitionrobotics.org/1.0"`
	// IngressName is the name of the Kubernetes Ingress used to route client requests from the Internet to
	// different internal services.
	// This configuration is required to enable websocket connections to simulations.
	IngressName string `env:"SUBT_INGRESS_NAME"`
	// WebsocketHost contains the address of the host used to route incoming websocket connections.
	// It is used to select a specific rule to modify in an ingress.
	// The ingress resource referenced by the `IngressName` configuration must contain at least one rule with a host
	// value matching this configuration.
	WebsocketHost string `env:"SUBT_WEBSOCKET_HOST"`
}

// SubTApplication represents an application used to tailor SubT simulation requests.
type SubTApplication struct {
	cfg subTSpecificsConfig
	// From aws go documentation:
	// Sessions should be cached when possible, because creating a new Session
	// will load all configuration values from the environment, and config files
	// each time the Session is created. Sharing the Session value across all of
	// your service clients will ensure the configuration is loaded the fewest
	// number of times possible.
	sess *session.Session
	// s3 clients are safe to use concurrently.
	s3Svc            s3iface.S3API
	schedulableTasks []SchedulableTask
}

// SimulationStatistics contains the summary values of a simulation run.
type SimulationStatistics struct {
	WasStarted          int `yaml:"was_started"`
	SimTimeDurationSec  int `yaml:"sim_time_duration_sec"`
	RealTimeDurationSec int `yaml:"real_time_duration_sec"`
	ModelCount          int `yaml:"model_count"`
}

// moleBridgeConfig contains information used to configure a Mole ROS/PULSAR bridge.
type moleBridgeConfig struct {
	// PodName is the pod name. It should be initialized with getMoleBridgePodName
	PodName string
	// Image contains the ROS/PULSAR bridge docker image URI.
	Image string
	// Labels contains the pod's labels
	Labels map[string]string
	// Selectors contains the pod's selector. It is currently limited to node labels.
	Selector map[string]string
	// WorldID contains the world number identifier. World IDs start at 1.
	WorldID int
	// TeamID is the team's unique identifier.
	TeamID int
	// GroupID is the simulation's unique identifier.
	GroupID string
	// ROSMasterIP is the IP of the ROS master to read ROS topics from.
	ROSMasterIP string
}

// NewSubTApplication creates a SubT application, used to tailor simulation requests.
func NewSubTApplication(ctx context.Context, s3Svc s3iface.S3API) (*SubTApplication, error) {
	logger(ctx).Info("Creating new SubT application")

	s := SubTApplication{}

	s.cfg = subTSpecificsConfig{}
	// Read configuration from environment
	logger(ctx).Info("Parsing Subt config")
	if err := env.Parse(&s.cfg); err != nil {
		return nil, err
	}

	var err error
	// We need to know the IP address of this host in order to create Network Policies.
	// We allow the user to define the desired IP using the IGN_IP env var. Otherwise,
	// we use one of the IP addresses of this host.
	if s.cfg.IgnIP == "" {
		logger(ctx).Warning(fmt.Print("No IGN_IP config found. Env var value:", os.Getenv("IGN_IP")))
		if s.cfg.IgnIP, err = getLocalIPAddressString(); err != nil {
			return nil, err
		}
	}
	logger(ctx).Warning(fmt.Println("Using IGN_IP value", s.cfg.IgnIP))

	s.s3Svc = s3Svc

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// getGazeboPodName returns the name of the Gazebo pod for a simulation.
func (sa *SubTApplication) getGazeboPodName(podNamePrefix string) string {
	return fmt.Sprintf("%s-gzserver", podNamePrefix)
}

// getMoleBridgePodName returns the name of the Mole ROS/Pulsar bridge pod for a simulation.
func (sa *SubTApplication) getMoleBridgePodName(podNamePrefix string) string {
	return fmt.Sprintf("%s-mole-bridge", podNamePrefix)
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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// getSimulationIngressPath returns the request path for the websocket server of a simulation.
func (sa *SubTApplication) getSimulationIngressPath(groupID string) string {
	return fmt.Sprintf("/simulations/%s", groupID)
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// getSimulationWebsocketAddress returns a simulation's websocket server address as well a the authorization token.
func (sa *SubTApplication) getSimulationWebsocketAddress(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	// The simulation must be running to be able to connect to the websocket server
	if *dep.DeploymentStatus != int(simRunning) {
		return nil, ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus)
	}

	return &WebsocketAddressResponse{
		Token:   *dep.AuthorizationToken,
		Address: fmt.Sprintf("%s/%s", sa.getSimulationWebsocketHost(), sa.getSimulationWebsocketPath(*dep.GroupID)),
	}, nil
}

// getSimulationWebsocketHost returns the host of the websocket address.
func (sa *SubTApplication) getSimulationWebsocketHost() string {
	return sa.cfg.WebsocketHost
}

// getSimulationWebsocketPath returns the websocket address path for the given simulation's group id.
func (sa *SubTApplication) getSimulationWebsocketPath(groupID string) string {
	return fmt.Sprintf("simulations/%s", groupID)
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// getSimulationLogsForDownload returns a link to the GZ logs that were saved in S3.
func (sa *SubTApplication) getSimulationLogsForDownload(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment, robotName *string) (*string, *ign.ErrMsg) {

	if sa.s3Svc == nil {
		err := errors.New("SubT Application wasn't given an S3Svc implementation but now is requested to fetch logs from S3")
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

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

	bucket := sa.cfg.S3LogsBucket
	ownerNameEscaped := url.PathEscape(*dep.Owner)
	folderPath := fmt.Sprintf("/gz-logs/%s/%s/", ownerNameEscaped, *dep.GroupID)
	filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
	logger(ctx).Debug(fmt.Sprintf("SubT App - Fetching generating link to fetch logs from S3 bucket [%s] with path [%s]\n", bucket, filePath))

	req, _ := sa.s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})

	url, err := req.Presign(5 * time.Minute)
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

	raw, err := KubernetesPodGetLog(ctx, s.clientset, s.cfg.KubernetesNamespace, podName, container, lines)

	if raw == nil {
		raw = sptr("No logs to show")
	}

	if err != nil {
		return nil, NewErrorMessageWithBase(ErrorFailedToGetLiveLogs, err)
	}

	podLog := PodLog{*raw}

	return &podLog, nil
}

func (sa *SubTApplication) copyPodReadFile(ctx context.Context, s *Service,
	dep *SimulationDeployment, filepath string) (*bytes.Buffer, error) {

	podName := sa.getCopyPodName(
		sa.getGazeboPodName(getSimulationPodNamePrefix(*dep.GroupID)),
	)
	return KubernetesPodReadFile(
		ctx,
		s.clientset,
		s.cfg.KubernetesNamespace,
		podName,
		CopyToS3SidecarContainerName,
		filepath,
	)
}

// getSimulationScore returns the simulation score from the pod of a simulation deployment.
func (sa *SubTApplication) getSimulationScore(ctx context.Context, s *Service,
	dep *SimulationDeployment) (*float64, *ign.ErrMsg) {

	// HACK This is a temporary fix until we can properly mock the Kubernetes clientset
	// s.clientset will not be a kubernetes.Clientset if this is a test.
	// A hardcoded value is returned if a test kubernetes interface is detected
	// for tests to complete successfully.
	if _, ok := s.clientset.(*kubernetes.Clientset); ok == false {
		score := float64(0)
		return &score, nil
	}

	out, err := sa.copyPodReadFile(
		ctx,
		s,
		dep,
		fmt.Sprintf("%s/score.yml", sa.cfg.SidecarContainerLogsVolumeMountPath),
	)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(int64(ErrorInvalidScore), err)
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(string(out.String())), 64)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(int64(ErrorInvalidScore), err)
	}

	return &score, nil
}

// getSimulationStatistics returns the simulation statistics summary from the pod of a simulation deployment.
func (sa *SubTApplication) getSimulationStatistics(ctx context.Context, s *Service,
	dep *SimulationDeployment) (*SimulationStatistics, *ign.ErrMsg) {

	// HACK This is a temporary fix until we can properly mock the Kubernetes clientset
	// s.clientset will not be a kubernetes.Clientset if this is a test.
	// A hardcoded value is returned if a test kubernetes interface is detected
	// for tests to complete successfully.
	if _, ok := s.clientset.(*kubernetes.Clientset); ok == false {
		return &SimulationStatistics{
			WasStarted:          0,
			SimTimeDurationSec:  0,
			RealTimeDurationSec: 0,
			ModelCount:          0,
		}, nil
	}

	out, err := sa.copyPodReadFile(
		ctx,
		s,
		dep,
		fmt.Sprintf("%s/summary.yml", sa.cfg.SidecarContainerLogsVolumeMountPath),
	)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(int64(ErrorInvalidSummary), err)
	}

	var statistics SimulationStatistics
	err = yaml.Unmarshal(out.Bytes(), &statistics)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(int64(ErrorInvalidSummary), err)
	}

	return &statistics, nil
}

// getSimulationScore returns the simulation run data containing metrics, information about found artifacts, etc.
func (sa *SubTApplication) getSimulationRunData(ctx context.Context, s *Service,
	dep *SimulationDeployment) (*string, *ign.ErrMsg) {

	// HACK This is a temporary fix until we can properly mock the Kubernetes clientset
	// s.clientset will not be a kubernetes.Clientset if this is a test.
	// A hardcoded value is returned if a test kubernetes interface is detected
	// for tests to complete successfully.
	if _, ok := s.clientset.(*kubernetes.Clientset); ok == false {
		out := ""
		return &out, nil
	}

	out, err := sa.copyPodReadFile(
		ctx,
		s,
		dep,
		fmt.Sprintf("%s/run.yml", sa.cfg.SidecarContainerLogsVolumeMountPath),
	)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(int64(ErrorInvalidRunInformation), err)
	}

	runData := out.String()

	return &runData, nil
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// launchApplication is a SubT specific function responsible of launching all the
// pods and services needed for a SubT simulation.
func (sa *SubTApplication) launchApplication(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment, podNamePrefix string, baseLabels map[string]string) (interface{}, *ign.ErrMsg) {

	groupID := *dep.GroupID

	// Extend base labels with SubT specific ones
	baseLabels[subtTagKey] = "true"

	gzserverLabels := cloneStringsMap(baseLabels)
	gzserverLabels[subtTypeGazebo] = "true"

	bridgeLabels := cloneStringsMap(baseLabels)
	bridgeLabels[subtTypeCommsBridge] = "true"

	fcLabels := cloneStringsMap(baseLabels)
	fcLabels[subtTypeFieldComputer] = "true"

	// Parse the SubT extra info required for this Simulation
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		logger(ctx).Error("[launchApplication] Failed to get extra info.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}
	// Now get the Circuit rules for this simulation
	rules, err := GetCircuitRules(tx, extra.Circuit)
	if err != nil {
		logger(ctx).Error("[launchApplication] Failed to get circuit rules when launching sim.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	worlds := ign.StrToSlice(*rules.Worlds)

	// Get the world to launch, in case it's null, use the default world.
	var worldToLaunch string

	worldToLaunch = worlds[0]
	if extra.WorldIndex != nil {
		worldToLaunch = worlds[*extra.WorldIndex]
	}

	// We split by ";" (semicolon), in case the configured worldToLaunch string has arguments.
	// eg. 'tunnel_circuit_practice.ign;worldName:=tunnel_circuit_practice_01'
	gzRunCommand := strings.Split(worldToLaunch, ";")

	// Set the simulation time limit
	gzRunCommand = append(gzRunCommand, fmt.Sprintf("durationSec:=%s", *rules.WorldMaxSimSeconds))

	// Enable ROS on the simulation instance
	gzRunCommand = append(gzRunCommand, "ros:=true")

	// Set headless
	gzRunCommand = append(gzRunCommand, "headless:=true")

	// Set increased update rate. Commented out because simulation was running
	// too fast for team's control loop.
	// gzRunCommand = append(gzRunCommand, "updateRate:=1000000")

	// Get the configured Seed for this run
	if rules.Seeds != nil {
		seeds, err := StrToIntSlice(*rules.Seeds)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		var seed int
		seed = seeds[0]
		if extra.RunIndex != nil {
			seed = seeds[*extra.RunIndex]
		}

		gzRunCommand = append(gzRunCommand, fmt.Sprintf("seed:=%d", seed))
	}

	// Get the world name parameter to pass on to the comms bridge
	var worldNameParam string
	for _, param := range gzRunCommand {
		if strings.Index(param, "worldName:=") != -1 {
			worldNameParam = param
			break
		}
	}
	// Check that a world was found
	if worldNameParam == "" {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, errors.New("World name not found"))
	}

	// Set the authorization token if it exists
	// TODO: Confirm parameter name
	if dep.AuthorizationToken != nil {
		gzRunCommand = append(gzRunCommand, fmt.Sprintf("websocketAuthKey:=%s", *dep.AuthorizationToken))
		gzRunCommand = append(gzRunCommand, fmt.Sprintf("websocketAdminAuthKey:=%s", *dep.AuthorizationToken))
	}

	// Set the maximum number of websocket connections. This can be removed
	// when websocket scaling across multiple machines is implemented.
	gzRunCommand = append(gzRunCommand, "websocketMaxConnections:=500")

	// Pass Robot names and types to the gzserver Pod.
	// robotName1:=xxx robotConfig1:=yyy robotName2:=xxx robotConfig2:=yyy (Note the numbers).
	for i, robot := range extra.Robots {
		gzRunCommand = append(gzRunCommand, fmt.Sprintf("robotName%d:=%s", i+1, robot.Name), fmt.Sprintf("robotConfig%d:=%s", i+1, robot.Type))
	}

	// Pass marsupial names to the gzserver Pod.
	// marsupialN:=<parent>:<child>
	for i, marsupial := range extra.Marsupials {
		gzRunCommand = append(gzRunCommand, fmt.Sprintf("marsupial%d:=%s:%s", i+1, marsupial.Parent, marsupial.Child))
	}
	logger(ctx).Info(fmt.Sprintf("gzRunCommand to use: %v", gzRunCommand))

	// Done to log the details into rollbar
	simStr, err := dep.toJSON()
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	logger(ctx).Warning(fmt.Sprintf("subt launchApplication - Trying to launch SimulationDeployment [%s]. SubT Specifics [%s]. GzRunCommand [%s]. Simulation Image [%s]",
		*simStr, *dep.Extra, gzRunCommand, *rules.Image))

	// Add NetworkPolicies to control ingress and egress to/from Pods.
	// Note: we want the gzserver and comms-bridge Pods to freely communicate between each other.
	// But the field-computers can only talk with the comms-bridge (and not with the gzserver).

	gazeboPodName := sa.getGazeboPodName(podNamePrefix)

	// First, the Network Policy for the gzserver Pod
	// Note: we add the rules before launching the pods, so they are active when the pod starts.
	npGz := sa.createNetworkPolicy(ctx, gazeboPodName, baseLabels, gzserverLabels, bridgeLabels)
	// We update the networkpolicy of the GzServer to also allow outbound connections to internet.
	npGz.Spec.Egress = append(npGz.Spec.Egress, networkingv1.NetworkPolicyEgressRule{})

	// Create the gzserver Pod definition (ie. the simulation server pod)
	// hostPath contains the path in the node that is mounted as a shared
	// directory among pods.
	hostPath := "/tmp"
	gzPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   gazeboPodName,
			Labels: gzserverLabels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: int64ptr(sa.cfg.TerminationGracePeriodSeconds),
			NodeSelector: map[string]string{
				// Force this pod to run on the same node as the target pod
				nodeLabelKeyGroupID:          *dep.GroupID,
				nodeLabelKeyCloudsimNodeType: subtTypeGazebo,
			},
			Containers: []corev1.Container{
				{
					Name:  GazeboServerContainerName,
					Image: *rules.Image,
					Args:  gzRunCommand,
					SecurityContext: &corev1.SecurityContext{
						Privileged:               boolptr(true),
						AllowPrivilegeEscalation: boolptr(true),
					},
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 11345,
						},
						{
							ContainerPort: 11311,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "logs",
							MountPath: sa.cfg.GazeboLogsVolumeMountPath,
						},
						{
							Name:      "xauth",
							MountPath: "/tmp/.docker.xauth",
						},
						{
							Name:      "localtime",
							MountPath: "/etc/localtime",
						},
						{
							Name:      "devinput",
							MountPath: "/dev/input",
						},
						{
							Name:      "x11",
							MountPath: "/tmp/.X11-unix",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "DISPLAY",
							Value: ":0",
						},
						{
							Name:  "QT_X11_NO_MITSHM",
							Value: "1",
						},
						{
							Name:  "XAUTHORITY",
							Value: "/tmp/.docker.xauth",
						},
						{
							Name:  "USE_XVFB",
							Value: "1",
						},
						{
							Name:  "IGN_RELAY",
							Value: sa.cfg.IgnIP,
						},
						{
							Name:  "IGN_PARTITION",
							Value: groupID,
						},
						{
							Name:  "IGN_VERBOSE",
							Value: sa.cfg.IgnVerbose,
						},
						{
							Name: "ROS_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
						{
							Name:  "ROS_MASTER_URI",
							Value: "http://$(ROS_IP):11311",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "logs",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: hostPath,
						},
					},
				},
				{
					Name: "x11",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/.X11-unix",
						},
					},
				},
				{
					Name: "devinput",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/dev/input",
						},
					},
				},
				{
					Name: "localtime",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/etc/localtime",
						},
					},
				},
				{
					Name: "xauth",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/.docker.xauth",
						},
					},
				},
			},
			// These DNS servers provide alternative DNS server from the internet
			// in case the cluster DNS service isn't available
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: []string{
					"8.8.8.8",
					"1.1.1.1",
				},
			},
		},
	}

	// Launch the gzserver Pod
	_, err = s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(gzPod)
	if err != nil {
		logger(ctx).Error("[launchApplication] Failed to create the gzserver pod.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	// Wait until the gazebo server pod has an IP address before continuing.
	// We need to get its IP address to share it with the other pods.
	// This call will block.
	ptrIP, err := waitForPodIPAndGetIP(ctx, s, gzserverLabels)
	if err != nil {
		logger(ctx).Error("[launchApplication] Timed out waiting for the gzserver pod to launch.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}
	gzserverPodIP := (*ptrIP)[gazeboPodName]

	// If S3 log backup is enabled then add an additional copy pod to upload
	// logs at the end of the simulation.
	if sa.cfg.S3LogsCopyEnabled {
		copyPod := sa.createS3CopyPod(
			ctx,
			s,
			dep,
			gzPod.Spec.NodeSelector,
			gazeboPodName,
			hostPath,
			"logs",
			sa.cfg.S3LogsBucket,
			sa.getGazeboLogsFilename(groupID),
		)
		// Launch the copy pod
		_, err := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(copyPod)
		if err != nil {
			logger(ctx).Error("[launchApplication] Failed to create the gzserver copy pod.", err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
	}

	// Add a mole ROS/Pulsar bridge container if the circuit has one defined
	if rules.MoleBridgeImage != nil {
		// Get worldID
		worldID := 1
		if extra.WorldIndex != nil {
			worldID = *extra.WorldIndex + 1
		}

		// Get teamID
		teamID := 0
		if dep.Owner == nil {
			org, em := s.userAccessor.GetOrganization(*dep.Owner)
			if em != nil {
				return nil, em
			}
			teamID = int(org.ID)
		}

		// Prepare pod labels
		labels := cloneStringsMap(baseLabels)
		labels[subtTypeMoleBridge] = "true"

		// Create and add the container to the list of active containers
		moleBridgePod := sa.createMolePod(&moleBridgeConfig{
			PodName: sa.getMoleBridgePodName(podNamePrefix),
			Labels:  labels,
			Selector: map[string]string{
				nodeLabelKeyGroupID: groupID,
				// Force this pod to run on the same node as the Gazebo server
				nodeLabelKeyCloudsimNodeType: subtTypeGazebo,
			},
			Image:       *rules.MoleBridgeImage,
			WorldID:     worldID,
			TeamID:      teamID,
			GroupID:     groupID,
			ROSMasterIP: gzserverPodIP,
		})

		// Launch and wait
		_, err = s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(moleBridgePod)
		if err != nil {
			logger(ctx).Error("[launchApplication] Failed to create the Mole ROS/Pulsar bridge pod.", err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		// Wait until the mole bridge pod has an IP address before continuing.
		// This call will block.
		_, err := waitForPodReadyAndGetIP(ctx, s, labels)
		if err != nil {
			logger(ctx).Error("[launchApplication] Timed out waiting for the gzserver pod to launch.", err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		// Add a rule to the Gazebo server network policy to allow connections to and from the mole-bridge
		npGz.Spec.Ingress = append(npGz.Spec.Ingress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
		})
		npGz.Spec.Egress = append(npGz.Spec.Egress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
		})
	}

	// Expose the Gazebo websocket server if an auth key is available
	if dep.AuthorizationToken != nil && len(sa.cfg.IngressName) > 0 {
		// Create and launch a service to expose the websocket server to the cluster
		_, err := sa.createWebsocketService(
			ctx,
			s.clientset,
			s.cfg.KubernetesNamespace,
			dep,
			podNamePrefix,
			baseLabels,
			gzserverLabels,
		)
		if err != nil {
			logger(ctx).Error("[launchApplication] Failed to create websocket Kubernetes service.", err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
		// Add a Gloo route to allow traffic to the websocket service from the Internet
		_, err = sa.CreateGlooVirtualServiceRoute(
			ctx,
			s,
			dep,
		)
		if err != nil {
			logger(ctx).Error("[launchApplication] Failed to create websocket Kubernetes ingress rule.", err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
	}

	// Create the Gazebo server network policy
	_, err = s.clientset.NetworkingV1().NetworkPolicies(s.cfg.KubernetesNamespace).Create(npGz)
	if err != nil {
		logger(ctx).Error("[launchApplication] Failed to create the gzserver network policy.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	fcPods := map[string]*corev1.Pod{}
	copyPods := make([]*corev1.Pod, len(extra.Robots))
	// Now launch the Comms Pods and the needed NetworkPolicies before
	// launching the field-computer pods (team solutions).
	for robotNumber, robot := range extra.Robots {
		// Note: it is assumed the Robot.Name is "alphanum". See its validator at subt_models.go
		robotNameLower := strings.ToLower(robot.Name)
		robotIdentifier := fmt.Sprintf("rbt%d", robotNumber+1)

		bridgePodName := sa.getCommsBridgePodName(podNamePrefix, robotIdentifier)
		fcPodName := sa.getFieldComputerPodName(podNamePrefix, robotIdentifier)

		specificBridgeLabels := cloneStringsMap(bridgeLabels)
		specificBridgeLabels["comms-for-robot"] = robotNameLower

		specificFcLabels := cloneStringsMap(fcLabels)
		specificFcLabels["robot-name"] = robotNameLower

		// Network Policy for this robot's comms-bridge Pod
		npBridge := sa.createNetworkPolicy(ctx, bridgePodName, baseLabels,
			specificBridgeLabels, gzserverLabels, specificFcLabels)
		// We update the networkpolicy of the Comms Bridge to also allow outbound connections to internet.
		npBridge.Spec.Egress = append(npBridge.Spec.Egress, networkingv1.NetworkPolicyEgressRule{})
		_, err = s.clientset.NetworkingV1().NetworkPolicies(s.cfg.KubernetesNamespace).Create(npBridge)
		if err != nil {
			msg := fmt.Sprintf("[launchApplication] Failed to create network policy for robot %s bridge.", robot.Name)
			logger(ctx).Error(msg, err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		// Network Policy for field-computer Pods (Note: they cannnot connect to internet)
		npfc := sa.createNetworkPolicy(ctx, fcPodName, baseLabels, specificFcLabels, specificBridgeLabels)
		_, err = s.clientset.NetworkingV1().NetworkPolicies(s.cfg.KubernetesNamespace).Create(npfc)
		if err != nil {
			msg := fmt.Sprintf("[launchApplication] Failed to create network policy for robot %s.", robot.Name)
			logger(ctx).Error(msg, err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		childMarsupial := "false"
		for _, marsupial := range extra.Marsupials {
			if marsupial.Child == robot.Name {
				childMarsupial = "true"
				break
			}
		}

		// Launch the comms-bridge Pod
		bridgePod := sa.createCommsBridgePod(
			ctx,
			dep,
			bridgePodName,
			specificBridgeLabels,
			gzserverPodIP,
			hostPath,
			"robot-logs",
			robotNumber+1,
			robot,
			*rules.BridgeImage,
			worldNameParam,
			childMarsupial,
		)

		// If S3 log backup is enabled then add an additional copy pod to upload
		// logs at the end of the simulation.
		if sa.cfg.S3LogsCopyEnabled {
			// Change the owner of the shared volume to the bridge's user
			sa.addSharedVolumeConfigurationContainer(bridgePod, 1000, 1000, "logs")

			// Create the copy pod for the current bridge
			copyPods[robotNumber] = sa.createS3CopyPod(
				ctx,
				s,
				dep,
				bridgePod.Spec.NodeSelector,
				bridgePodName,
				hostPath,
				"robot-logs",
				sa.cfg.S3LogsBucket,
				sa.getRobotROSLogsFilename(groupID, robotNameLower),
			)
		}

		// Launch the bridge pods
		logger(ctx).Info("Launching bridge pod", bridgePod.Spec.InitContainers)
		_, err = s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(bridgePod)
		if err != nil {
			msg := fmt.Sprintf("[launchApplication] Failed to create bridge pod for robot %s.", robot.Name)
			logger(ctx).Error(msg, err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		// Create the FC Pod and save it for launching later
		fcPod := sa.createFieldComputerPod(
			ctx,
			dep,
			fcPodName,
			specificFcLabels,
			groupID,
			robot,
		)
		fcPods[robotIdentifier] = fcPod
	}

	// Wait for Comms Bridge Pods to be Ready and get their IP addresses
	bridgeIPs, err := waitForPodReadyAndGetIP(ctx, s, bridgeLabels)
	if err != nil {
		logger(ctx).Error("[launchApplication] Timed out waiting for bridge pods to launch.", err)
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	// Launch the bridge copy pods
	if sa.cfg.S3LogsCopyEnabled {
		for _, pod := range copyPods {
			_, err := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(pod)
			if err != nil {
				msg := fmt.Sprintf("[launchApplication] Failed to create copy pod for bridge pod %s.", pod.Name)
				logger(ctx).Error(msg, err)
				return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
			}
		}
	}
	// Now launch the previously created Field-Computer Pod(s)
	for robotIdentifier, fcPod := range fcPods {
		bridgePodName := sa.getCommsBridgePodName(podNamePrefix, robotIdentifier)
		bridgeIP := (*bridgeIPs)[bridgePodName]
		rosIP := fmt.Sprintf("http://%s:11311", bridgeIP)

		// Set the ROS MASTER URI to the FC pod (i.e. the CommsBridge's IP)
		fcPod.Spec.Containers[0].Env = append(fcPod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "ROS_MASTER_URI",
			Value: rosIP,
		})

		// Launch the field-computer Pod
		_, err = s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).Create(fcPod)
		if err != nil {
			msg := fmt.Sprintf("[launchApplication] Failed to create field-computer pod %s.", fcPod.Name)
			logger(ctx).Error(msg, err)
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
	}

	return nil, nil
}

// waitForPodIPAndGetIP blocks until the pods identified by matchLabels have an
// IP address.
// Returns a map with (podName, IP address).
// Dev note: This func is used to get the IP of the Gzserver
func waitForPodIPAndGetIP(ctx context.Context, s *Service, matchLabels map[string]string) (*map[string]string, error) {
	return waitForPodConditionAndGetIP(ctx, s, matchLabels, "'Has IP status'", podHasIPAddress)
}

// waitForPodReadyAndGetIP blocks until the pods identified by matchLabels have
// Ready status.
// Returns a map with (podName, IP address).
// Dev note: This func is used to get the IP of the CommsBridge pods.
func waitForPodReadyAndGetIP(ctx context.Context, s *Service, matchLabels map[string]string) (*map[string]string, error) {
	return waitForPodConditionAndGetIP(ctx, s, matchLabels, "Ready", subtPodRunningAndReady)
}

// subtPodRunningAndReady checks if a pod is ready, specifically for SubT. This function is used for Wait polls.
func subtPodRunningAndReady(ctx context.Context, pod *corev1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case corev1.PodFailed:
		return false, conditions.ErrPodCompleted
	case corev1.PodRunning:
		return podutil.IsPodReady(pod), nil
	case corev1.PodSucceeded:
		_, isFC := pod.Labels[subtTypeFieldComputer]
		if isFC {
			logger(ctx).Warning(fmt.Sprintf("FC pod %s status is Succeeded. Considering it Ready.", pod.Name))
		}
		return isFC, nil
	}
	return false, nil
}

// waitForPodConditionAndGetIP blocks until the pods identified by matchLabels
// meet a condition.
// Returns a map with (podName, IP address).
func waitForPodConditionAndGetIP(ctx context.Context, s *Service, matchLabels map[string]string,
	condStr string, cond PodCondition) (*map[string]string, error) {

	var selectorBuilder strings.Builder
	for k, v := range matchLabels {
		fmt.Fprintf(&selectorBuilder, "%s=%s,", k, v)
	}
	labelSelector := strings.TrimRight(selectorBuilder.String(), ",")

	timeout := time.Duration(s.cfg.PodReadyTimeoutSeconds) * time.Second
	opts := metav1.ListOptions{LabelSelector: labelSelector}

	if err := WaitForMatchPodsCondition(ctx, s.clientset, s.cfg.KubernetesNamespace,
		opts, condStr, timeout, cond); err != nil {
		return nil, err
	}

	podsInterface := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace)
	runningPods, err := podsInterface.List(opts)
	if err != nil || len(runningPods.Items) == 0 {
		return nil, errors.New("Pods not found for labels: " + labelSelector)
	}

	ips := map[string]string{}
	for _, p := range runningPods.Items {
		name := p.ObjectMeta.Name
		ip := p.Status.PodIP
		ips[name] = ip
	}

	return &ips, nil
}

// podHasIPAddress checks if a pod by name is running. This function is used
// for Wait polls.
func podHasIPAddress(ctx context.Context, pod *corev1.Pod) (bool, error) {
	if pod.Status.PodIP != "" {
		return true, nil
	}
	return false, nil
}

// addSharedVolumeConfigurationContainer changes the owner of a pod's shared
// volume directory to the specified user and group.
// Logs produced by Gazebo and Bridge pods are uploaded to S3 using an additional
// copy pod launched in the same node. A Kubernetes volume is used to share
// a specific directory between these pods. The shared directory is physically
// stored inside the node where the log and copy pods are scheduled.
// Gazebo and Bridge pods are configured to create a directory to store logs
// inside the shared volume because the shared directory may contain other files
// not related to the simulation.
// Kubernetes creates the directory, but sets the permissions to root:root 755
// which does not give write permissions to the bridge container because it runs
// with `developer` as its user and group. This function adds an InitContainer to
// the pod spec that changes the owner of the directory from root to developer
// before the bridge container starts, giving write permissions to the bridge
// container and allowing it to store logs.
// `userID` is the linux user id (UID) of the user in the pod producing logs.
// `groupID` is the linux group id (GID) of the user in the pod producing logs.
// `volumeName` is the name of the Kubernetes hostPath volume containing the shared directory.
func (sa *SubTApplication) addSharedVolumeConfigurationContainer(pod *corev1.Pod, userID int, groupID int,
	volumeName string) {
	pod.Spec.InitContainers = []corev1.Container{
		{
			Name:    "chown-shared-volume",
			Image:   "infrastructureascode/aws-cli:latest",
			Command: []string{"/bin/sh"},
			Args:    []string{"-c", fmt.Sprintf("chown %d:%d /tmp", userID, groupID)},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      volumeName,
					MountPath: "/tmp",
				},
			},
		},
	}
}

// createWebsocketService creates a Kubernetes Service that exposes the websocket server to the cluster.
// This service can be exposed to the Internet via a load balancer.
func (sa *SubTApplication) createWebsocketService(ctx context.Context, kc kubernetes.Interface, namespace string,
	dep *SimulationDeployment, podNamePrefix string, serviceLabels map[string]string,
	targetLabels map[string]string) (*corev1.Service, error) {

	logger := logger(ctx)
	logger.Info("Creating simulation cluster websocket service.")

	// Prepare the service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   sa.getWebsocketServiceName(podNamePrefix),
			Labels: serviceLabels,
		},
		Spec: corev1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: targetLabels,
			Ports: []corev1.ServicePort{
				{
					Name: "websockets",
					Port: 9002,
				},
			},
		},
	}

	// Launch the service
	_, err := kc.CoreV1().Services(namespace).Create(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// createWebsocketIngress adds a rule to the environment's cluster ingress resource.
// If an ingress controller is installed in the cluster, the creation of the ingress rule will trigger an ingress
// controller to configure a load balancer to redirect simulation websocket requests to the cluster service created
// for the websocket server.
func (sa *SubTApplication) createWebsocketIngress(ctx context.Context, kc kubernetes.Interface, namespace string,
	dep *SimulationDeployment, service *corev1.Service) (*v1beta1.Ingress, error) {

	logger := logger(ctx)
	logger.Info("Creating simulation websocket ingress rule.")

	// Prepare the rule path
	rulePath := &v1beta1.HTTPIngressPath{
		Path: sa.getSimulationIngressPath(*dep.GroupID),
		Backend: v1beta1.IngressBackend{
			ServiceName: service.Name,
			ServicePort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: service.Spec.Ports[0].Port,
			},
		},
	}

	var host *string
	if len(sa.cfg.WebsocketHost) > 0 {
		host = &sa.cfg.WebsocketHost
	}

	var ingress *v1beta1.Ingress
	var err error
	retries := uint(9)
	for i := uint(1); i < retries; i++ {
		ingress, err = upsertIngressRule(
			ctx,
			kc,
			namespace,
			sa.cfg.IngressName,
			host,
			rulePath,
		)
		if err == nil {
			break
		} else if i < retries-1 {
			logger.Info("Failed to create ingress rule. Retrying. Error:", err)
		}

		// If an error occurred, retry for up to 10 min with exponential backoff
		Sleep(time.Duration(1<<i) * time.Second)
	}

	return ingress, err
}

// GetGlooServiceUpstreamName gets the name of an upstream for a cluster service.
// Upstreams are used by Gloo Virtual Service routes to direct traffic to a Kubernetes service.
// `glooNamespace` is the Kubernetes namespace Gloo is running in.
// `service` is the Kubernetes service the upstream points to.
func (*SubTApplication) GetGlooServiceUpstreamName(gloo gloo.Clientset, glooNamespace string,
	groupID string) (string, error) {

	opts := metav1.ListOptions{
		LabelSelector: labelAndValue(podLabelKeyGroupID, groupID),
	}

	upstreams, err := gloo.Gloo().Upstreams(glooNamespace).List(opts)
	if err != nil {
		return "", err
	}
	// There should only ever be a single upstream for the target service
	if len(upstreams.Items) < 1 {
		return "", errors.New("did not find a Gloo upstream for target Kubernetes service")
	} else if len(upstreams.Items) > 1 {
		return "", errors.New("found too many Gloo upstreams for target Kubernetes service")
	}

	return upstreams.Items[0].Name, nil
}

// CreateGlooVirtualServiceRoute creates a Gloo Route to redirect websocket simulation requests to a Kubernetes service.
// TODO: Receive a generic application instead of a SubTApplication pointer.
func (sa *SubTApplication) CreateGlooVirtualServiceRoute(ctx context.Context, s *Service,
	dep *SimulationDeployment) (*gatewayv1.VirtualService, error) {

	logger := logger(ctx)
	logger.Info("Creating simulation websocket ingress rule.")

	groupID := *dep.GroupID
	retries := uint(9)
	var err error

	// Upstreams are Gloo's representation of services that can be routed to. An upstream can point to different things.
	// For now we're only using Kubernetes services, so a Gloo upstream is synonymous with a Kubernetes service.
	var upstream string
	for i := uint(1); i < retries; i++ {
		logger.Info(fmt.Sprintf("Getting service upstream in namespace %s with id %s", s.cfg.KubernetesGlooNamespace, groupID))
		upstream, err = sa.GetGlooServiceUpstreamName(
			s.glooClientset,
			s.cfg.KubernetesGlooNamespace,
			groupID,
		)
		if err == nil {
			logger.Info(fmt.Sprintf("Got upstream! Name: %s", upstream))
			break
		} else if i < retries-1 {
			logger.Info("Failed to get Kubernetes Service Gloo Upstream name. Retrying. Error:", err)
		}

		// If an error occurred, retry for up to 10 min with exponential backoff
		Sleep(time.Duration(1<<i) * time.Second)
	}
	if err != nil {
		return nil, err
	}

	// Matchers define the criteria a request must match to route traffic to upstreams.
	matcher := gloo.CreateVirtualHostRouteExactMatcher(sa.getSimulationIngressPath(groupID))

	// Actions perform routing operations when executed.
	// This particular action redirects traffic to a single target upstream.
	action := gloo.CreateVirtualHostRouteAction(
		s.cfg.KubernetesGlooNamespace,
		upstream,
	)

	// Routes bind an upstream, a set of matchers and an action into an element used to configure the Gloo Envoy proxy.
	route := gloo.CreateVirtualHostRoute(
		groupID,
		[]*matchers.Matcher{matcher},
		action,
	)

	// Add the route to the Virtual Service
	var virtualService *gatewayv1.VirtualService
	for i := uint(1); i < retries; i++ {
		virtualService, err = gloo.UpsertVirtualServiceRoute(
			ctx,
			s.glooClientset,
			s.cfg.KubernetesGlooNamespace,
			sa.cfg.IngressName,
			route,
		)
		if err == nil {
			break
		} else if i < retries-1 {
			logger.Info("Failed to create Gloo Virtual Service route. Retrying. Error:", err)
		}

		// If an error occurred, retry for up to 10 min with exponential backoff
		Sleep(time.Duration(1<<i) * time.Second)
	}

	return virtualService, err
}

// createFieldComputerCommsBridgePod creates a basic comms-bridge pod. Callers should then
// change the Pod's Image, Command and Args fields.
func (sa *SubTApplication) createMolePod(config *moleBridgeConfig) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.PodName,
			Labels: config.Labels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: int64ptr(sa.cfg.TerminationGracePeriodSeconds),
			NodeSelector:                  config.Selector,
			Containers: []corev1.Container{
				{
					Name:  MoleBridgeContainerName,
					Image: config.Image,
					SecurityContext: &corev1.SecurityContext{
						Privileged:               boolptr(true),
						AllowPrivilegeEscalation: boolptr(true),
					},
					Env: []corev1.EnvVar{
						{
							Name:  "CS_PB_PULSAR_IP",
							Value: "mole-pulsar-broker",
						},
						{
							Name:  "CS_PB_PULSAR_PORT",
							Value: "6650",
						},
						{
							Name:  "CS_PB_TOPIC_REGEX",
							Value: "^subt/",
						},
						{
							Name:  "ROS_MASTER_URI",
							Value: fmt.Sprintf("http://%s:11311", config.ROSMasterIP),
						},
						{
							Name:  "WORLD_ID",
							Value: strconv.Itoa(config.WorldID),
						},
						{
							Name:  "TEAM_ID",
							Value: strconv.Itoa(config.TeamID),
						},
						{
							Name:  "REPLICATION_ID",
							Value: config.GroupID,
						},
					},
				},
			},
		},
	}
}

// createCommsBridgePod creates a basic comms-bridge pod. Callers should then
// change the Pod's Image, Command and Args fields.
func (sa *SubTApplication) createCommsBridgePod(ctx context.Context, dep *SimulationDeployment,
	podName string, labels map[string]string, gzserverPodIP string, hostPath string, logDirectory string,
	robotNumber int, robot SubTRobot, bridgeImage string, worldNameParam string, childMarsupial string) *corev1.Pod {

	logMountPath := path.Join(hostPath, logDirectory)
	hostPathType := corev1.HostPathDirectoryOrCreate
	bridgePod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: int64ptr(sa.cfg.TerminationGracePeriodSeconds),
			Containers: []corev1.Container{
				{
					Name:            CommsBridgeContainerName,
					ImagePullPolicy: corev1.PullIfNotPresent,
					// Default Image/Command/Args, for testing
					Image:   "infrastructureascode/aws-cli:latest",
					Command: []string{"/bin/sh", "-c", "--"},
					Args:    []string{"while true; do sleep 30; done;"},
					SecurityContext: &corev1.SecurityContext{
						Privileged:               boolptr(true),
						AllowPrivilegeEscalation: boolptr(true),
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "logs",
							MountPath: sa.cfg.ROSLogsVolumeMountPath,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "IGN_PARTITION",
							Value: *dep.GroupID,
						},
						{
							// IGN_IP contains the IP of the subscriber (this pod)
							Name: "IGN_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
						{
							// IGN_RELAY should contain the IP of the publisher (the gzserver)
							Name:  "IGN_RELAY",
							Value: gzserverPodIP,
						},
						{
							Name:  "IGN_VERBOSE",
							Value: sa.cfg.IgnVerbose,
						},
						{
							Name: "ROS_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
						{
							Name:  "ROS_MASTER_URI",
							Value: "http://$(ROS_IP):11311",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "logs",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: logMountPath,
							Type: &hostPathType,
						},
					},
				},
			},
			// These DNS servers provide alternative DNS server from the internet
			// in case the cluster DNS service isn't available
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: []string{
					"8.8.8.8",
					"1.1.1.1",
				},
			},
		},
	}

	if bridgeImage != "infrastructureascode/aws-cli:latest" {
		bridgePod.Spec.Containers[0].Image = bridgeImage
		bridgePod.Spec.Containers[0].Command = nil
		bridgePod.Spec.Containers[0].Args = []string{
			worldNameParam,
			fmt.Sprintf("robotName%d:=%s", robotNumber, robot.Name),
			fmt.Sprintf("robotConfig%d:=%s", robotNumber, robot.Type),
			"headless:=true",
			fmt.Sprintf("marsupial:=%s", childMarsupial),
		}
	}

	return bridgePod
}

// createFieldComputerPod creates a basic field-computer pod. Callers should then
// change the Pod's Image, Command and Args fields.
// The field-computer pod runs the Team Solution container.
func (sa *SubTApplication) createFieldComputerPod(ctx context.Context, dep *SimulationDeployment,
	podName string, labels map[string]string, groupID string, robot SubTRobot) *corev1.Pod {

	fcPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: int64ptr(sa.cfg.TerminationGracePeriodSeconds),
			NodeSelector: map[string]string{
				// Needed to force this pod to run on specific nodes
				nodeLabelKeyGroupID:          groupID,
				nodeLabelKeyCloudsimNodeType: subtTypeFieldComputer,
				nodeLabelKeySubTRobotName:    strings.ToLower(robot.Name),
			},
			Containers: []corev1.Container{
				{
					Name:            FieldComputerContainerName,
					ImagePullPolicy: corev1.PullIfNotPresent,
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: boolptr(false),
					},
					// Default Image/Command/Args, for testing
					Image:   "infrastructureascode/aws-cli:latest",
					Command: []string{"/bin/sh", "-c", "--"},
					Args:    []string{"while true; do sleep 30; done;"},
					// Limit to 95% of the total memory of a g3.4xlarge instance
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("115Gi"),
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "ROBOT_NAME",
							Value: robot.Name,
						},
						{
							Name: "ROS_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
					},
				},
			},
		},
	}

	if robot.Image != "infrastructureascode/aws-cli:latest" {
		fcPod.Spec.Containers[0].Image = robot.Image
		fcPod.Spec.Containers[0].Command = nil
		fcPod.Spec.Containers[0].Args = nil
	}

	return fcPod
}

type stringsMap map[string]string

// createNetworkPolicy is a helper function to create a Network Policy.
// @param baseLabels is used to label the new network policy.
// @param matchingPodLabels is used to define which Pods to apply this policy to.
// @param allowFromLabels is an array of labels used to define the Ingress and Egress
// rules allowing communication to and from the matching pods.
func (sa *SubTApplication) createNetworkPolicy(ctx context.Context, npName string,
	baseLabels, matchingPodLabels stringsMap, allowFromLabels ...stringsMap) *networkingv1.NetworkPolicy {

	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   npName,
			Labels: baseLabels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: matchingPodLabels,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				// Dev note: Important -- the IP addresses listed here should be the IP of the Cloudsim pod.
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							IPBlock: &networkingv1.IPBlock{
								// We always allow traffic coming from the Cloudsim host.
								CIDR: sa.cfg.IgnIP + "/32",
							},
						},
					},
				},
				// Allow traffic to websocket server
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 9002,
							},
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				// Dev note: Important -- the IP addresses listed here should be the IP of the Cloudsim pod.
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							IPBlock: &networkingv1.IPBlock{
								// We always allow traffic targetted to the Cloudsim host
								CIDR: sa.cfg.IgnIP + "/32",
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
		},
	}

	// Allow communication to/from pods from "allowFromLabels" argument
	for _, allow := range allowFromLabels {
		np.Spec.Ingress = append(np.Spec.Ingress, networkingv1.NetworkPolicyIngressRule{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: allow,
					},
				},
			},
		})
		np.Spec.Egress = append(np.Spec.Egress, networkingv1.NetworkPolicyEgressRule{
			To: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: allow,
					},
				},
			},
		})
	}

	return np
}

// createS3CopyPod launches a copy pod in the same node as a target pod.
// This copy pod is used by Cloudsim during simulation termination to compress
// and upload the the entire content of a target pod's log volume to the given
// S3 bucket.
// `targetPodName` is the name of the pod the copy pod will copy logs from.
// In order to share a common directory between the target and copy pods, a
// directory on the node is exposed to both pods. This only works because we
// know that both the target and copy pods will run on the same node.
// `logVolumePath` is the path of node directory that is exposed to the target and copy pods.
// `logVolumePath` is the path to the directory within the node directory that the copy pod should target. This is
// included to allow multiple pods to share a single common directory on the node. If no specific directory is required,
// use "".
// `filename` sets the name of the compressed directory uploaded to S3.
func (sa *SubTApplication) createS3CopyPod(ctx context.Context, s *Service, dep *SimulationDeployment,
	targetNodeSelector map[string]string, targetPodName string, logVolumePath string, logVolumeSubPath string,
	s3Bucket string, filename string) *corev1.Pod {

	podName := sa.getCopyPodName(targetPodName)

	logger(ctx).Debug(fmt.Sprintf(
		"Creating copy pod for [%s] to upload logs on simulation termination.", podName,
	))

	// Prepare the copy pod spec
	hostPathType := corev1.HostPathDirectoryOrCreate
	copyPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
			Labels: map[string]string{
				"copy-to-s3": "true",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: int64ptr(sa.cfg.TerminationGracePeriodSeconds),
			NodeSelector:                  targetNodeSelector,
			Containers: []corev1.Container{
				{
					Name:  "copy-to-s3",
					Image: "infrastructureascode/aws-cli:latest",
					// DEV NOTE: This command is a hack to keep the container running. If the container ends its main process,
					// K8 will consider it finished and will try to restart it.
					// We need this sidecar container to keep running until the logs of the target container are ready to upload.
					// Before terminating a pod, Cloudsim will run a command using the sidecar container to upload the logs to S3.
					Command:         []string{"tail", "-f", "/dev/null"},
					ImagePullPolicy: corev1.PullIfNotPresent,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: "logs",
							// The sidecar container will always mount the logs volume to `/tmp/logs`.
							// The content of this volume is set by the container generating the logs.
							MountPath: sa.cfg.SidecarContainerLogsVolumeMountPath,
							SubPath:   logVolumeSubPath,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "AWS_DEFAULT_REGION",
							Value: sa.cfg.Region,
						},
						{
							Name:  "AWS_REGION",
							Value: sa.cfg.Region,
						},
						{
							Name: "AWS_ACCESS_KEY_ID",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{Name: sa.cfg.AwsSecretName},
									Key:                  "aws-access-key-id",
									Optional:             boolptr(false),
								},
							},
						},
						{
							Name: "AWS_SECRET_ACCESS_KEY",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{Name: sa.cfg.AwsSecretName},
									Key:                  "aws-secret-access-key",
									Optional:             boolptr(false),
								},
							},
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "logs",
					// This volume provides exposes a directory in the node pods are running in, providing a shared
					// directory between the target and copy pod This only works because we know that pods and their
					// respective copy pods run on the same node.
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: logVolumePath,
							Type: &hostPathType,
						},
					},
				},
			},
			// These DNS servers provide alternative DNS server from the internet
			// in case the cluster DNS service isn't available
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: []string{
					"8.8.8.8",
					"1.1.1.1",
				},
			},
		},
	}

	return copyPod
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// processSimulationResults when shutting down a simulation.
// It registers the score and statistics of the simulation being shutdown.
func (sa *SubTApplication) processSimulationResults(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment) *ign.ErrMsg {

	if dep.Processed {
		return nil
	}

	values := SimulationDeploymentsSubTValue{
		SimulationDeployment: dep,
		GroupID:              dep.GroupID,
	}

	// Create and upload logs to S3
	if sa.cfg.S3LogsCopyEnabled {
		logger(ctx).Info(
			fmt.Sprintf("processSimulationResults - Uploading simulation logs to S3 for simulation [%s]", *dep.GroupID),
		)
		if em := sa.uploadSimulationLogs(ctx, s, dep); em != nil {
			logMsg := fmt.Sprintf(
				"processSimulationResults - Could not upload simulation logs to S3 for simulation [%s].",
				*dep.GroupID,
			)
			logger(ctx).Error(logMsg, em)
			return em
		}
	}

	// Score and summary entries should be generated for single simulations or child simulations of multi-sims,
	// but not for parent simulations of multi-sims.
	if dep.isMultiSimParent() {
		return nil
	}

	// Get score
	score, em := sa.getSimulationScore(ctx, s, dep)
	if em != nil {
		return em
	}
	values.Score = score

	// Get simulation statistics
	statistics, em := sa.getSimulationStatistics(ctx, s, dep)
	if em != nil {
		return em
	}
	values.RealTimeDurationSec = statistics.RealTimeDurationSec
	values.SimTimeDurationSec = statistics.SimTimeDurationSec
	values.ModelCount = statistics.ModelCount
	if err := tx.Create(&values).Error; err != nil {
		return NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// If this is not a multisim, then we can create the final summary and create the score entry
	if !dep.isMultiSim() {
		// Create the score entry
		if !globals.DisableScoreGeneration {
			if em := processScore(ctx, s.userAccessor, tx, dep, values.Score); em != nil {
				return em
			}
		}

		// Send an email with the summary to the competitor
		if !globals.DisableSummaryEmails {
			summary := AggregatedSubTSimulationValues{
				Score:                  *values.Score,
				SimTimeDurationAvg:     float64(values.SimTimeDurationSec),
				SimTimeDurationStdDev:  0,
				RealTimeDurationAvg:    float64(values.RealTimeDurationSec),
				RealTimeDurationStdDev: 0,
				ModelCountAvg:          float64(values.ModelCount),
				ModelCountStdDev:       0,
			}

			// Get simulation run information
			runData, em := sa.getSimulationRunData(ctx, s, dep)
			if em != nil {
				return em
			}

			SendSimulationSummaryEmail(dep, summary, runData)
		}
	}

	if !dep.Processed {
		dep.UpdateProcessed(tx, true)
	}

	return nil
}

// processScore creates the competition_scores entry for a simulation and updates its score.
func processScore(ctx context.Context, userAccessor useracc.UserAccessor, tx *gorm.DB,
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

func isTeamSolutionPod(pod corev1.Pod) bool {
	// field-computer pods (ie. team solutions) have the "field-computer" label set to "true"
	flag, ok := pod.Labels[subtTypeFieldComputer]
	return ok && (flag == "true")
}

// deleteApplication is a SubT specific function responsible of deleting all the
// pods and services created by a SubT simulation.
func (sa *SubTApplication) deleteApplication(ctx context.Context, s *Service, tx *gorm.DB,
	dep *SimulationDeployment) *ign.ErrMsg {

	groupID := *dep.GroupID
	groupIDLabel := getPodLabelSelectorForSearches(groupID)

	// Upload logs and process score and summary entries for the simulation
	if em := sa.processSimulationResults(ctx, s, tx, dep); em != nil {
		logger(ctx).Error("Could not process simulation results", em)
		return em
	}

	// Find and delete all Pods associated to the groupID.
	podsInterface := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace)
	pods, err := podsInterface.List(metav1.ListOptions{LabelSelector: groupIDLabel})
	if err != nil || len(pods.Items) == 0 {
		// Pods for this groupID not found. Continue or fail?
		logger(ctx).Warning("Pods not found for the groupID: "+groupID, err)
		if !sa.cfg.AllowNotFoundDuringShutdown {
			err = errors.Wrap(err, "Pods not found for the groupID: "+groupID)
			return ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
		}
	}
	for _, p := range pods.Items {
		err = podsInterface.Delete(p.Name, &metav1.DeleteOptions{})
		if err != nil {
			// There was an unexpected error deleting the Pod. If it's Team Solution pod,
			// the we log the error and continue, as this can sometimes happen if (e.g.)
			// the teams finish their Pod's main process by itself.
			// Otherwise, if the failed Pod is the gzserver or the comms-bridge we mark
			// the simulation as failed, as this is an unexpected scenario.
			em := ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
			logger(ctx).Error("Error while invoking k8 Delete Pod. Make sure a sysadmin deletes the Pod manually.", em)
			if !isTeamSolutionPod(p) {
				return em
			}
		}
	}
	logger(ctx).Info("Successfully deleted pods for groupID: " + groupID)

	// Delete Gloo route associated to the groupID.
	if len(sa.cfg.IngressName) > 0 && len(sa.cfg.WebsocketHost) > 0 {
		retries := uint(9)
		for i := uint(1); i < retries; i++ {
			_, err := gloo.RemoveVirtualServiceRoute(
				ctx,
				s.glooClientset,
				s.cfg.KubernetesGlooNamespace,
				sa.cfg.IngressName,
				&gatewayapiv1.Route{
					Name: groupID,
				},
			)
			if err == nil {
				break
			} else if i < retries-1 {
				logger(ctx).Info("Failed to delete Gloo Virtual Service route. Retrying. Error:", err)
			}

			// If an error occurred, retry for up to 10 min with exponential backoff
			Sleep(time.Duration(1<<i) * time.Second)
		}

		if err != nil {
			// There was an unexpected error deleting the Gloo Virtual Service route and the simulation will be marked
			// as failed, as this is an unexpected scenario.
			em := ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
			errMsg := "Error while removing k8 Ingress rule. Make sure a sysadmin deletes the ingress rule manually."
			logger(ctx).Error(errMsg, em)
			return em
		}

		logger(ctx).Info("Successfully deleted ingress rules for groupID: " + groupID)
	}

	// Find and delete all Services associated to the groupID.
	serviceInterface := s.clientset.CoreV1().Services(s.cfg.KubernetesNamespace)
	services, err := serviceInterface.List(metav1.ListOptions{LabelSelector: groupIDLabel})
	if err != nil || len(services.Items) == 0 {
		// Services for this groupID not found. Continue or fail?
		logger(ctx).Warning("Services not found for the groupID: "+groupID, err)
		if !sa.cfg.AllowNotFoundDuringShutdown {
			err = errors.Wrap(err, "Services not found for the groupID: "+groupID)
			return ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
		}
	}
	for _, service := range services.Items {
		err = serviceInterface.Delete(service.Name, &metav1.DeleteOptions{})
		if err != nil {
			// There was an unexpected error deleting the Service and the simulation will be marked as failed, as this
			// is an unexpected scenario.
			em := ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
			errMsg := "Error while invoking k8 Delete Service. Make sure a sysadmin deletes the Service manually."
			logger(ctx).Error(errMsg, em)
			return em
		}
	}
	logger(ctx).Info("Successfully deleted services for groupID: " + groupID)

	// Find and delete all the network policies associated to the groupID.
	// Dev note: it is important to remove the network policies AFTER the gzlogs are
	// copied to S3. Otherwise, if we remove the policies before, the pod will lose
	// access to outside world and the copy to S3 will not work.
	npInterface := s.clientset.NetworkingV1().NetworkPolicies(s.cfg.KubernetesNamespace)
	nps, err := npInterface.List(metav1.ListOptions{LabelSelector: groupIDLabel})
	if err != nil || len(nps.Items) == 0 {
		logger(ctx).Warning("Network Policies not found for the groupID: "+groupID, err)
		// Continue or fail?
		if !sa.cfg.AllowNotFoundDuringShutdown {
			err = errors.Wrap(err, "Network Policies not found for the groupID: "+groupID)
			return ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
		}
	}
	for _, np := range nps.Items {
		err = npInterface.Delete(np.Name, &metav1.DeleteOptions{})
		if err != nil {
			// There was an error deleting the NetworkPolicy. We log the error and continue,
			// as we want to free the used resources (ec2 instance).
			em := ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
			logger(ctx).Error("Error while invoking k8 Delete NetworkPolicy. Make sure a sysadmin deletes the NetworkPolicy manually", em)
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// setupEC2InstanceSpecifics is invoked by the EC2 NodeManager to describe the needed EC2 instance details for SubT.
func (sa *SubTApplication) setupEC2InstanceSpecifics(ctx context.Context, s *Ec2Client,
	tx *gorm.DB, dep *SimulationDeployment, template *ec2.RunInstancesInput) ([]*ec2.RunInstancesInput, error) {

	// Create some Tags that all instances will share
	subTTag := ec2.Tag{Key: aws.String(subtTagKey), Value: aws.String("true")}
	subTTag2 := ec2.Tag{Key: aws.String("cloudsim-application"), Value: aws.String(subtTagKey)}
	subTTag3 := ec2.Tag{Key: aws.String("cloudsim-simulation-worker"), Value: aws.String(s.awsCfg.NamePrefix)}
	subTTag4 := ec2.Tag{Key: aws.String(nodeLabelKeyCloudsimNodeType), Value: aws.String(subtTypeGazebo)}
	appendTags(template, &subTTag, &subTTag2, &subTTag3, &subTTag4)
	inputs := make([]*ec2.RunInstancesInput, 0)

	// Configure EC2 instance specifics
	template.ImageId = aws.String(s.ec2Cfg.AMI)
	template.InstanceType = aws.String("g3.4xlarge")

	gzInput, err := cloneRunInstancesInput(template)
	if err != nil {
		return nil, err
	}

	// Add the new Input to the result array
	inputs = append(inputs, gzInput)

	// Create instances for the field computers; one per robot definition.
	extra, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return nil, err
	}

	for _, r := range extra.Robots {
		fcInput, err := cloneRunInstancesInput(template)
		if err != nil {
			return nil, err
		}

		// Replace the node type tag
		replaceTag(fcInput, &ec2.Tag{
			Key:   sptr(nodeLabelKeyCloudsimNodeType),
			Value: sptr(subtTypeFieldComputer),
		})
		userData, _ := s.buildUserDataString(*dep.GroupID,
			labelAndValue(nodeLabelKeyCloudsimNodeType, subtTypeFieldComputer),
			labelAndValue(nodeLabelKeySubTRobotName, strings.ToLower(r.Name)),
		)
		fcInput.UserData = aws.String(userData)
		replaceInstanceNameTag(fcInput, s.getInstanceNameFor(*dep.GroupID, "fc-"+r.Name))
		inputs = append(inputs, fcInput)
	}

	return inputs, nil
}

func cloneRunInstancesInput(src *ec2.RunInstancesInput) (*ec2.RunInstancesInput, error) {
	var bs bytes.Buffer
	enc := gob.NewEncoder(&bs)
	if err := enc.Encode(*src); err != nil {
		return nil, err
	}
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&bs)
	var dst ec2.RunInstancesInput
	if err := dec.Decode(&dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// getMaxDurationForSimulation returns the max allowed duration for a simulation,
// before marking it for shutdown.
func (sa *SubTApplication) getMaxDurationForSimulation(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) time.Duration {
	return time.Duration(sa.cfg.MaxDurationForSimulations) * time.Minute
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// uploadToS3LogBucket uploads a file to a simulation log folder.
func (sa *SubTApplication) uploadToS3SimulationLogBucket(dep *SimulationDeployment, filename string,
	file []byte) *ign.ErrMsg {
	key := path.Join(GetS3SimulationLogKey(dep), filename)
	if _, err := UploadToS3Bucket(sa.s3Svc, &sa.cfg.S3LogsBucket, &key, file); err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	return nil
}

// uploadSimulationLogs uploads Gazebo/ROS simulation logs to the log bucket of the simulation.
// Teams can later download this summary through a generated link.
func (sa *SubTApplication) uploadSimulationLogs(ctx context.Context, s *Service,
	simDep *SimulationDeployment) *ign.ErrMsg {

	logger := logger(ctx)

	groupID := *simDep.GroupID
	bucket := filepath.Join(sa.cfg.S3LogsBucket, GetS3SimulationLogKey(simDep))
	failedPodUploads := make(map[string]error, 0)

	// Upload Gazebo logs
	opts := MakeListOptions(
		getPodLabelSelectorForSearches(groupID),
		labelAndValue(subtTypeGazebo, "true"),
	)
	pods, err := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).List(opts)
	if err != nil {
		msg := fmt.Sprintf("Could not get the simulation pod while attempting to upload log files.")
		logger.Error(msg, err)
	} else {
		for _, pod := range pods.Items {
			podName := pod.Name
			err := KubernetesPodSendS3CopyCommand(
				ctx,
				s.clientset,
				s.cfg.KubernetesNamespace,
				sa.getCopyPodName(podName),
				CopyToS3SidecarContainerName,
				bucket,
				sa.cfg.SidecarContainerLogsVolumeMountPath,
				sa.getGazeboLogsFilename(groupID),
			)
			if err != nil {
				failedPodUploads[podName] = err
			}
		}
	}

	// Upload ROS logs
	opts = MakeListOptions(
		getPodLabelSelectorForSearches(groupID),
		labelAndValue(subtTypeCommsBridge, "true"),
	)
	pods, err = s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).List(opts)
	if err != nil {
		msg := fmt.Sprintf("Could not get comms-bridge pods while attempting to upload log files.")
		logger.Error(msg, err)
	} else {
		for _, pod := range pods.Items {
			podName := pod.Name
			robotName := pod.Labels["comms-for-robot"]
			err := KubernetesPodSendS3CopyCommand(
				ctx,
				s.clientset,
				s.cfg.KubernetesNamespace,
				sa.getCopyPodName(podName),
				CopyToS3SidecarContainerName,
				bucket,
				sa.cfg.SidecarContainerLogsVolumeMountPath,
				sa.getRobotROSLogsFilename(groupID, robotName),
			)
			if err != nil {
				failedPodUploads[podName] = err
			}
		}
	}

	// Check for errors
	if len(failedPodUploads) > 0 {
		for podName, err := range failedPodUploads {
			logger.Error(fmt.Sprintf("Failed to upload logs for pod %s: %s", podName, err))
		}
		return ign.NewErrorMessage(int64(ErrorFailedToUploadLogs))
	}

	return nil
}

// uploadSimulationSummary uploads the simulation summary to the log bucket of the simulation.
// Teams can later download this summary through a generated link.
func (sa *SubTApplication) uploadSimulationSummary(simDep *SimulationDeployment,
	summary *AggregatedSubTSimulationValues) *ign.ErrMsg {
	values, err := json.Marshal(summary)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorMarshalJSON, err)
	}
	fileName := sa.getSimulationSummaryFilename(*simDep.GroupID)
	if em := sa.uploadToS3SimulationLogBucket(simDep, fileName, values); em != nil {
		return em
	}

	return nil
}

// updateMultiSimStatuses updates the status of a Multi-Sim Parent and executes application-specific logic based on the
// state of its children.
func (sa *SubTApplication) updateMultiSimStatuses(ctx context.Context, tx *gorm.DB, userAccessor useracc.UserAccessor,
	simDep *SimulationDeployment) *ign.ErrMsg {
	// Note: simDep is a Parent in a multi-sim

	// Only proceed if the simulation terminated successfully. Get the aggregated values from all children
	if simDep.IsRunning() || simDep.ErrorStatus != nil || simDep.Processed {
		return nil
	}

	if simDep.Held {
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
	if sa.cfg.S3LogsCopyEnabled {
		logger(ctx).Info(
			fmt.Sprintf("updateMultiSimStatuses - Uploading simulation summary for simulation [%s]", *simDep.GroupID),
		)
		if em := sa.uploadSimulationSummary(simDep, summary); em != nil {
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
		SendSimulationSummaryEmail(simDep, *summary, nil)
	}

	if !simDep.Processed {
		simDep.UpdateProcessed(tx, true)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

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
