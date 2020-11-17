package simulations

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"gitlab.com/ignitionrobotics/web/cloudsim/queues"
	"gitlab.com/ignitionrobotics/web/cloudsim/simulations/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/transport"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/transport/ign"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	The Simulations Service is in charge of launching and terminating Gazebo simulations. And,
	in case of an error, it is responsible of rolling back the failed operation.

	To do this and handle some concurrency without exhausting the host, it has
	one worker-thread-pool for each main activity (launch, terminate, error handling).
	The `launch` and `terminate` pools can launch 10 concurrent workers (eg. the launcher can
	launch 10 simulations in parallel). The error handler pool only has one worker.

	In addition, the SimService has a background go routine that checks for expired
	simulations and send those to automatic termination.

	The Simulations Service interacts with a given NodeManager to start and terminate
	Nodes in the Kubernetes cluster.
	Some examples of NodeManager implementations are EC2Client and LocalNodes.

	This service also delegates to "Application" and "Platform" specific handlers, so they can
	manage the specific details of the simulations to launch. As an example, a competition
	can have custom request for the Nodes to be launched and the simulation details.

*/

// TODO pending set 1 pod per Node in Affinity or other conf.

// SimService is an interface that supports managing simulation instances.
type SimService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	CustomizeSimRequest(ctx context.Context, r *http.Request, tx *gorm.DB, createSim *CreateSimulation, username string) *ign.ErrMsg
	GetCloudMachineInstances(ctx context.Context, p *ign.PaginationRequest,
		tx *gorm.DB, byStatus *MachineStatus, invertStatus bool, groupID *string, user *users.User, application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg)
	DeleteNodesAndHostsForGroup(ctx context.Context, tx *gorm.DB,
		dep *SimulationDeployment, user *users.User) (interface{}, *ign.ErrMsg)
	GetSimulationDeployment(ctx context.Context, tx *gorm.DB, groupID string,
		user *users.User) (interface{}, *ign.ErrMsg)
	GetSimulationWebsocketAddress(ctx context.Context, tx *gorm.DB, user *users.User,
		groupID string) (interface{}, *ign.ErrMsg)
	GetSimulationLiveLogs(ctx context.Context, tx *gorm.DB, user *users.User, groupID string,
		robotName *string, lines *int64) (interface{}, *ign.ErrMsg)
	GetSimulationLogsForDownload(ctx context.Context, tx *gorm.DB, user *users.User, groupID string,
		robotName *string) (*string, *ign.ErrMsg)
	RegisterApplication(ctx context.Context, app ApplicationType)
	ShutdownSimulationAsync(ctx context.Context, tx *gorm.DB,
		groupID string, user *users.User) (interface{}, *ign.ErrMsg)
	SimulationDeploymentList(ctx context.Context, p *ign.PaginationRequest, tx *gorm.DB, byStatus *DeploymentStatus,
		invertStatus bool, byErrStatus *ErrorStatus, invertErrStatus bool, byCircuit *string, user *users.User,
		application *string, includeChildren bool) (*SimulationDeployments, *ign.PaginationResult, *ign.ErrMsg)
	StartSimulationAsync(ctx context.Context, tx *gorm.DB, createSim *CreateSimulation,
		user *users.User) (interface{}, *ign.ErrMsg)
	LaunchSimulationAsync(ctx context.Context, tx *gorm.DB, groupID string,
		user *users.User) (interface{}, *ign.ErrMsg)
	RestartSimulationAsync(ctx context.Context, tx *gorm.DB, groupID string,
		user *users.User) (interface{}, *ign.ErrMsg)
	GetRemainingSubmissions(ctx context.Context, tx *gorm.DB, user *users.User, circuit *string,
		owner *string) (interface{}, *ign.ErrMsg)
	CustomRuleList(ctx context.Context, p *ign.PaginationRequest, tx *gorm.DB, user *users.User, application *string,
		circuit *string, owner *string, ruleType *CustomRuleType) (*CircuitCustomRules, *ign.PaginationResult, *ign.ErrMsg)
	SetCustomRule(ctx context.Context, tx *gorm.DB, user *users.User, application *string,
		circuit *string, owner *string, ruleType *CustomRuleType, value *string) (*CircuitCustomRule, *ign.ErrMsg)
	DeleteCustomRule(ctx context.Context, tx *gorm.DB, user *users.User, application *string, circuit *string,
		owner *string, ruleType *CustomRuleType) (interface{}, *ign.ErrMsg)
	GetCompetitionRobots(applicationName string) (interface{}, *ign.ErrMsg)
	QueueGetElements(ctx context.Context, user *users.User, page, pageSize *int) ([]interface{}, *ign.ErrMsg)
	QueueMoveElementToFront(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg)
	QueueMoveElementToBack(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg)
	QueueSwapElements(ctx context.Context, user *users.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg)
	QueueRemoveElement(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg)
	QueueCount(ctx context.Context, user *users.User) (interface{}, *ign.ErrMsg)
}

// NodeManager is responsible of creating and removing cloud instances, and
// kubernetes nodes.
// NodeManager is the expected interface to be implemented by Cloudsim NodeManagers.
// Example implementations can be found in `ec2_machines.go` and `local_machines.go`.
type NodeManager interface {
	CloudMachinesList(ctx context.Context, p *ign.PaginationRequest,
		tx *gorm.DB, byStatus *MachineStatus, invertStatus bool, groupID *string, application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg)
	// Requests the NodeManager to terminate the hosts (or instances or VMs) used to run a GroupID.
	// It also updates the MachineInstance DB records with the status of the terminated hosts.
	deleteHosts(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (interface{}, *ign.ErrMsg)
	// Requests the NodeManager to delete involved k8 nodes.
	// It is expected that if the labeled Nodes cannot be found, then this function should return an ErrorLabeledNodeNotFound.
	deleteK8Nodes(ctx context.Context, tx *gorm.DB, groupID string) (interface{}, *ign.ErrMsg)
	// asks the NodeManager to launch a set of nodes to run a simulation
	launchNodes(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (*string, *ign.ErrMsg)
}

const (
	podLabelPodGroup       = "pod-group"
	podLabelKeyGroupID     = "cloudsim-group-id"
	cloudsimTagLabelKey    = "cloudsim"
	cloudsimTagLabel       = "cloudsim=true"
	launcherRelaunchNeeded = "relaunch"
)

// Service is the main struct exported by this Simulations service.
type Service struct {
	// Whether this service will automatically requeue instances that failed with
	// ErrorLaunchingCloudInstanceNotEnoughResources error. True by default.
	AllowRequeuing bool
	// A reference to the kubernetes client
	clientset kubernetes.Interface
	// A reference to the Gloo client
	glooClientset gloo.Clientset
	// A reference to the nodes manager implementation
	hostsSvc NodeManager
	DB       *gorm.DB
	// Workers (ie. Thread Pools)
	launcher     JobPool
	terminator   JobPool
	errorHandler JobPool
	// A buffered channel used to buffer requests to create simulations.
	// Items from this channel will be then used to feed the 'launcher' JobPool.
	launchHandlerQueue *queues.LaunchQueueRepository
	// A buffered channel used to buffer requests to shutdown simulations.
	// Items from this channel will be then used to feed the 'terminator' JobPool.
	terminationHandlerQueue chan string
	// A buffered channel used to buffer requests to handle simulation errors.
	// Items from this channel will be then used to feed the 'error handler' JobPool.
	errorHandlerQueue chan string
	// The base Context from the main application
	baseCtx context.Context
	// The service config. Read from environment variables
	cfg simServConfig
	// A map with the current live RunningSimulations
	runningSimulations map[string]*RunningSimulation
	// A mutex to protect access to read/write operations over the map
	lockRunningSimulations sync.RWMutex
	// Expired simulations cleaning process
	expiredSimulationsTicker *time.Ticker
	expiredSimulationsDone   chan bool
	// MultiSim Parent status updater routine
	multisimStatusUpdater     *time.Ticker
	multisimStatusUpdaterDone chan bool
	applications              map[string]ApplicationType
	// The UserAccessor to check for Users/Orgs permissions
	userAccessor             useracc.UserAccessor
	poolNotificationCallback PoolNotificationCallback
	scheduler                *scheduler.Scheduler
}

// SimServImpl holds the instance of the Simulations Service. It is set at initialization.
var SimServImpl SimService

type simServConfig struct {
	// KubernetesNamespace is the Kubernetes namespace for simulation resources.
	KubernetesNamespace string `env:"SIMSVC_KUBERNETES_NAMESPACE" envDefault:"default"`
	// KubernetesGlooNamespace is the Gloo namespace in the Kubernetes cluster.
	KubernetesGlooNamespace string `env:"SIMSVC_KUBERNETES_GLOO_NAMESPACE" envDefault:"gloo-system"`
	// PoolSizeLaunchSim is the number of worker threads available to launch simulations.
	PoolSizeLaunchSim int `env:"SIMSVC_POOL_LAUNCH_SIM" envDefault:"10"`
	// PoolSizeTerminateSim is the number of worker threads available to terminate simulations.
	PoolSizeTerminateSim int `env:"SIMSVC_POOL_TERMINATE_SIM" envDefault:"10"`
	// PoolSizeErrorHandler is the number of worker threads available to handle simulation errors.
	PoolSizeErrorHandler int `env:"SIMSVC_POOL_ERROR_HANDLER" envDefault:"20"`
	// Timeout in seconds to wait until a new Pod is ready. Default: 5 minutes.
	PodReadyTimeoutSeconds int `env:"SIMSVC_POD_READY_TIMEOUT_SECONDS" envDefault:"300"`
	// Timeout in seconds to wait until a new Node is ready. Default: 5 minutes.
	NodeReadyTimeoutSeconds int `env:"SIMSVC_NODE_READY_TIMEOUT_SECONDS" envDefault:"300"`
	// The number of live simulations a team can have running in parallel. Zero value means unlimited.
	MaxSimultaneousSimsPerOwner int `env:"SIMSVC_SIMULTANEOUS_SIMS_PER_TEAM" envDefault:"1"`
	// MaxDurationForSimulations is the maximum number of minutes a simulation can run in cloudsim.
	// This is a default value. Specific ApplicationTypes are expected to overwrite this.
	MaxDurationForSimulations int `env:"SIMSVC_SIM_MAX_DURATION_MINUTES" envDefault:"45"`
	// IsTest determines if a service is being used for a test
	IsTest bool
}

// ApplicationType represents an Application (eg. SubT). Applications are used
// to customize launched Simulations.
type ApplicationType interface {
	getApplicationName() string
	GetSchedulableTasks(ctx context.Context, s *Service, tx *gorm.DB) []SchedulableTask
	checkCanShutdownSimulation(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment, user *users.User) (bool, *ign.ErrMsg)
	checkValidNumberOfSimulations(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg
	customizeSimulationRequest(ctx context.Context, s *Service, r *http.Request, tx *gorm.DB,
		createSim *CreateSimulation, username string) *ign.ErrMsg
	// allow specific applications to create multiSimulations from a single CreateSimulation request.
	spawnChildSimulationDeployments(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) ([]*SimulationDeployment, *ign.ErrMsg)
	// invoked when a simulation is about to be restarted. The old simulation run should be invalidated.
	invalidateSimulation(ctx context.Context, tx *gorm.DB, simDep *SimulationDeployment) error
	deleteApplication(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg
	// allow specific applications to define the maximum allowed time for simulations. A value of 0 will
	// mean that the Cloudsim's default value should be used (defined by env var SIMSVC_SIM_MAX_DURATION_MINUTES).
	getMaxDurationForSimulation(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) time.Duration
	getGazeboWorldStatsTopicAndLimit(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (string, int, error)
	getGazeboWorldWarmupTopic(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (string, error)
	getSimulationWebsocketAddress(ctx context.Context, s *Service, tx *gorm.DB,
		dep *SimulationDeployment) (interface{}, *ign.ErrMsg)
	getSimulationWebsocketHost() string
	getSimulationWebsocketPath(groupID string) string
	getSimulationLogsForDownload(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment,
		robotName *string) (*string, *ign.ErrMsg)
	getSimulationLiveLogs(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment,
		robotName *string, lines int64) (interface{}, *ign.ErrMsg)
	launchApplication(ctx context.Context, s *Service, tx *gorm.DB, dep *SimulationDeployment,
		podNamePrefix string, baseLabels map[string]string) (interface{}, *ign.ErrMsg)
	updateMultiSimStatuses(ctx context.Context, tx *gorm.DB, userAccessor useracc.UserAccessor, simDep *SimulationDeployment) *ign.ErrMsg
	getCompetitionRobots() (interface{}, *ign.ErrMsg)
	ValidateSimulationLaunch(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg
	// TODO: Move simulationIsHeld to SubT implementation.
	simulationIsHeld(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) bool
}

// JobPool is a pool of Jobs that can accept jobs to be executed.
// For more details see project "github.com/panjf2000/ants".
type JobPool interface {
	Serve(args interface{}) error
}

// PoolFactory is a function responsible of initializing and returning a JobPool.
// Dev note: we created the PoolFactory abstraction to allow tests use
// synchronic pools.
type PoolFactory func(poolSize int, jobF func(interface{})) (JobPool, error)

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// NewSimulationsService creates a new simulations service
func NewSimulationsService(ctx context.Context, db *gorm.DB, nm NodeManager, kcli kubernetes.Interface,
	gloo gloo.Clientset, pf PoolFactory, ua useracc.UserAccessor, isTest bool) (SimService, error) {

	var err error
	s := Service{}
	s.AllowRequeuing = true
	s.DB = db
	s.baseCtx = ctx
	s.clientset = kcli
	s.glooClientset = gloo
	s.userAccessor = ua
	s.hostsSvc = nm
	s.runningSimulations = map[string]*RunningSimulation{}
	s.lockRunningSimulations = sync.RWMutex{}
	s.applications = map[string]ApplicationType{}
	s.scheduler = scheduler.GetInstance()

	// Read configuration from environment
	s.cfg = simServConfig{
		IsTest: isTest,
	}
	if err := env.Parse(&s.cfg); err != nil {
		return nil, err
	}

	// Configure the worker pools
	// Create the queues of pending user requests to process.
	// We use a buffered channel of a big size to avoid blocking callers (i.e. incoming http requests).
	s.launchHandlerQueue = queues.NewLaunchQueueRepository()
	s.terminationHandlerQueue = make(chan string, 1000)
	s.errorHandlerQueue = make(chan string, 1000)

	if pf == nil {
		pf = defaultPoolFactory
	}
	s.launcher, err = pf(s.cfg.PoolSizeLaunchSim, s.workerStartSimulation)
	if err != nil {
		return nil, err
	}

	s.terminator, err = pf(s.cfg.PoolSizeTerminateSim, s.workerTerminateSimulation)
	if err != nil {
		return nil, err
	}

	s.errorHandler, err = pf(s.cfg.PoolSizeErrorHandler, s.workerErrorHandler)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// PoolNotificationCallback type of the listeners
type PoolNotificationCallback func(poolEvent PoolEvent, groupID string, result interface{}, em *ign.ErrMsg)

// PoolEvent registers a single pool event listener that will receive
// notifications any time a pool worker "finishes" its job (either with result or error).
type PoolEvent int

// PoolEvent enum
const (
	PoolStartSimulation PoolEvent = iota
	PoolShutdownSimulation
	PoolRollbackFailedLaunch
	PoolCompleteFailedTermination
)

// SetPoolEventsListener sets a new PoolNotificationCallback to the poolNotificationCallback field.
func (s *Service) SetPoolEventsListener(cb PoolNotificationCallback) {
	s.poolNotificationCallback = cb
}

func (s *Service) notify(poolEvent PoolEvent, groupID string, result interface{}, em *ign.ErrMsg) {
	if s.poolNotificationCallback != nil {
		s.poolNotificationCallback(poolEvent, groupID, result, em)
	}
}

func (s *Service) queueLaunchRequest(groupID string) {
	s.launchHandlerQueue.Enqueue(groupID)
}

func (s *Service) queueShutdownRequest(groupID string) {
	s.terminationHandlerQueue <- groupID
}

func (s *Service) queueErrorHandlerRequest(groupID string) {
	s.errorHandlerQueue <- groupID
}

// defaultPoolFactory is the default implementation of the PoolFactory interface.
// It creates an ants.PoolWithFunc.
func defaultPoolFactory(poolSize int, jobF func(interface{})) (JobPool, error) {
	return ants.NewPoolWithFunc(poolSize, jobF)
}

// CustomizeSimRequest allows registered Applications to customize the incoming CreateSimulation request.
// Eg. reading specific SubT fields.
func (s *Service) CustomizeSimRequest(ctx context.Context, r *http.Request, tx *gorm.DB, createSim *CreateSimulation, username string) *ign.ErrMsg {
	return s.applications[createSim.Application].customizeSimulationRequest(ctx, s, r, tx, createSim, username)
}

// Start starts this simulation service. It needs to be invoked AFTER 'Applications'
// were registerd using 'RegisterApplication'.
func (s *Service) Start(ctx context.Context) error {
	// Start a routine that will move 'launch' requests from the Waiting Queue into
	// the WorkerPool. If all the Workers are busy then this goroutine will block.
	go func() {
		var groupID string
		var ok bool
		for {
			result, err := s.launchHandlerQueue.DequeueOrWait()
			groupID, ok = result.(string)
			if ok && err == nil {
				logger(ctx).Info("launchHandler queue - about to process launch task for groupID: " + groupID)
				// This call will block if all Workers are busy
				if err := s.launcher.Serve(groupID); err != nil {
					logMsg := fmt.Sprintf(
						"launchHandler queue - Error in launch task for groupID [%s]. Error:[%v]\n", groupID, err,
					)
					logger(ctx).Error(logMsg, err)
				}
			}
		}
	}()

	// Start a routine that will move 'shutdown' requests from the Waiting Queue into
	// the WorkerPool. If all the Workers are busy then this goroutine will block.
	go func() {
		for groupID := range s.terminationHandlerQueue {
			logger(ctx).Info("shutdownHandler queue - about to submit shutdown task for groupID: " + groupID)
			// This call will block if all Workers are busy
			if err := s.terminator.Serve(groupID); err != nil {
				logMsg := fmt.Sprintf(
					"shutdownHandler queue - Error in shutdown task for groupID [%s]. Error:[%v]\n", groupID, err,
				)
				logger(ctx).Error(logMsg, err)
			}
		}
	}()

	// Start a routine that will move 'error handler' requests from the Waiting Queue into
	// the WorkerPool. If all the Workers are busy then this goroutine will block.
	go func() {
		for groupID := range s.errorHandlerQueue {
			logger(ctx).Info("errorHandler queue - about to handle error for groupID: " + groupID)
			// This call will block if all Workers are busy
			if err := s.errorHandler.Serve(groupID); err != nil {
				logMsg := fmt.Sprintf(
					"errorHandler queue - Error while handling errors for groupID [%s]. Error:[%v]\n", groupID, err,
				)
				logger(ctx).Error(logMsg, err)
			}
		}
	}()

	// Initialize server state based on data from DB and and from kubernetes cluster Pods.
	// Important note: it is expected that the kubernetes cluster should be running already.
	if err := s.rebuildState(ctx, s.DB); err != nil {
		return err
	}
	s.StartExpiredSimulationsCleaner()
	s.StartMultiSimStatusUpdater()
	RegisterSchedulableTasks(s, ctx, s.DB)
	return nil
}

// Stop stops this Simulations Service
func (s *Service) Stop(ctx context.Context) error {
	s.StopExpiredSimulationsCleaner()
	s.StopMultiSimStatusUpdater()
	close(s.terminationHandlerQueue)
	close(s.errorHandlerQueue)
	return nil
}

// RegisterApplication registers a new application type.
func (s *Service) RegisterApplication(ctx context.Context, app ApplicationType) {
	logger(ctx).Info(fmt.Sprintf("Sim Service - Registered new Application [%s]", app.getApplicationName()))
	s.applications[app.getApplicationName()] = app
}

// GetNodeManager returns the NodeManager saved on hostsSvc
func (s *Service) GetNodeManager() NodeManager {
	return s.hostsSvc
}

// SetNodeManager sets a different to hostsSvc
func (s *Service) SetNodeManager(nm NodeManager) {
	s.hostsSvc = nm
}

// GetApplications returns a map of application name and application type.
func (s *Service) GetApplications() map[string]ApplicationType {
	return s.applications
}

// initializeRunningSimulationsFromCluster finds the existing Pods in the Kubernetes
// cluster and initializes the internal set of runningSimulations.
// Note: after a server restart there can be inconsistencies between DB data and
// live kubernetes. This function is not responsible for sanitizing such inconsistencies.
// TODO: There should be another call for SystemAdmins to list inconsistencies and allow them
// to act on those.
func (s *Service) initializeRunningSimulationsFromCluster(ctx context.Context, tx *gorm.DB) error {

	// Find all Pods associated to cloudsim
	podsInterface := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace)
	pods, err := podsInterface.List(metav1.ListOptions{LabelSelector: cloudsimTagLabel})
	if err != nil {
		logger(ctx).Error("Error getting initial list of Cloudsim Pods from cluster", err)
		return err
	}

	// First, filter the simulations that have all its Pods with status PodRunning.
	// Keep in mind that a simulation could have spawned multiple Pods.
	runningSims := make(map[string]bool)

	for _, p := range pods.Items {
		groupID := p.Labels[podLabelKeyGroupID]

		if p.ObjectMeta.DeletionTimestamp != nil {
			// DeletionTimestamp != nil means the system has requested a deletion of this Pod.
			// So, we won't consider this as a Running Pod.
			runningSims[groupID] = false
			continue
		}

		running, found := runningSims[groupID]
		if !found {
			// First pod processed for this simulation. Mark running with initial value to make the "&&"" work later
			running = true
		}
		// is the current pod running. Update the whole simulation running status based on that.
		running = running && (p.Status.Phase == corev1.PodRunning)
		runningSims[groupID] = running

	}

	// Now iterate the simulations marked as 'running' and create RunningSimulations for them.
	for groupID, running := range runningSims {
		if !running {
			continue
		}
		// Get the Simulation record from DB
		simDep, err := GetSimulationDeployment(tx, groupID)
		if err != nil {
			return err
		}

		// Only create a RunningSimulation if the whole simulation status was Running and the DB
		// deploymentStatus is Running as well.
		if simRunning.Eq(*simDep.DeploymentStatus) {
			// Register a new live RunningSimulation
			if err := s.createRunningSimulation(ctx, tx, simDep); err != nil {
				return err
			}

			logger(ctx).Info(fmt.Sprintf("Init - Added RunningSimulation for groupID: [%s]. Deployment Status in DB: [%d]", groupID, *simDep.DeploymentStatus))
		}
	}

	return nil
}

// DeployHeldCircuitSimulations launches the held simulation deployments for a given circuit
func (s *Service) DeployHeldCircuitSimulations(ctx context.Context, tx *gorm.DB, circuit string) error {
	deps, err := GetSimulationDeploymentsByCircuit(tx, circuit, simPending, simPending, boolptr(true))
	if err != nil {
		return err
	}
	for _, dep := range *deps {
		logger(ctx).Info(fmt.Sprintf("Deploying simulations -- Circuit: %s | Group ID: %s", circuit, *dep.GroupID))
		s.DeployHeldSimulation(ctx, tx, &dep)
	}
	return nil
}

// DeployHeldSimulation deploys a simulation that is being held by cloudsim
func (s *Service) DeployHeldSimulation(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {
	if err := dep.UpdateHeldStatus(tx, false); err != nil {
		return NewErrorMessageWithBase(ErrorLaunchHeldSimulation, err)
	}

	simsToLaunch, err := s.getLaunchableSimulations(ctx, tx, dep)
	if err != nil {
		return err
	}

	for _, sim := range simsToLaunch {
		if err := sim.UpdateHeldStatus(tx, false); err != nil {
			return NewErrorMessageWithBase(ErrorLaunchHeldSimulation, err)
		}

		logger(ctx).Info(fmt.Sprintf("DeployHeldSimulation about to submit launch task for groupID: %s", *sim.GroupID))
		if err := LaunchSimulation(s, ctx, tx, &sim); err != nil {
			logger(ctx).Error(fmt.Sprintf("DeployHeldSimulation -- Cannot launch simulation: %s", err.Msg))
		}
	}
	return nil
}

// rebuildState is called during this service startup to (re)build the queue of pending
// simulation requests, and also to mark with an Error status those simulations
// that were in the middle of a "launch" or "shutdown" operation when the server was
// previously stopped. Those simulations with error status will need to be reviewed by an admin.
func (s *Service) rebuildState(ctx context.Context, db *gorm.DB) error {

	// Initialize running simulation from the running kubernetes Pods.
	// Important note: it is expected that the kubernetes cluster should be running already.
	if err := s.initializeRunningSimulationsFromCluster(ctx, db); err != nil {
		return err
	}

	s.lockRunningSimulations.RLock()
	defer s.lockRunningSimulations.RUnlock()

	// Get all "single" or "child" simulations (ie. not Parent Sims) that were not fully terminated yet and without errors.
	// Those simulations could have been running during previous server run.
	var deps SimulationDeployments
	if err := db.Model(&SimulationDeployment{}).Where("error_status IS NULL").Where("multi_sim != ?", multiSimParent).
		Where("deployment_status BETWEEN ? AND ?", int(simPending), int(simTerminatingInstances)).Find(&deps).Error; err != nil {
		return err
	}

	for _, d := range deps {
		groupID := *d.GroupID

		if simPending.Eq(*d.DeploymentStatus) {
			// If still Pending then re-add it to the scheduler, by adding a 'launch simulation'
			// request to the Launcher Jobs-Pool
			logger(ctx).Info("rebuildState -- about to submit launch task for groupID: " + groupID)
			if err := LaunchSimulation(s, ctx, db, &d); err != nil {
				logger(ctx).Error(fmt.Sprintf("rebuildState -- Cannot launch simulation: %s", err.Msg))
			}
			continue
		}

		if simRunning.Eq(*d.DeploymentStatus) {
			_, podRunning := s.runningSimulations[groupID]
			if !podRunning {
				logger(ctx).Info(fmt.Sprintf("rebuildState -- GroupID [%s] expected to be Running "+
					"in DB but there is no matching Pod running. Marking with error", groupID))
				// if the SimulationDeployment DB record has 'running' status but there is no matching
				// running Pod in the cluster then we have an inconsistenty. Mark it as error.
				d.setErrorStatus(db, simErrorServerRestart)
			}
			continue
		}

		// For any other intermediate deployment status, we just mark the Simulation with an
		// Error, as we cannot confirm a successful completion of the ongoing operation
		// after a server restart.
		statusStr := DeploymentStatus(*d.DeploymentStatus).String()
		logger(ctx).Info(fmt.Sprintf("rebuildState -- GroupID [%s] found with intermediate "+
			"DeploymentStatus [%s]. Marking with error", groupID, statusStr))
		d.setErrorStatus(db, simErrorServerRestart)
	}

	return nil
}

// prepareSimulations prepares a Simulation Deployment to be launched and returns an array of simulations to deploy.
// If it is a multisim, it will return the child simulations.
// In any other cases, it will only return one simulation deployment.
func (s *Service) prepareSimulations(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) ([]*SimulationDeployment, *ign.ErrMsg) {
	simsToLaunch := []*SimulationDeployment{dep}
	childSims, em := s.applications[*dep.Application].spawnChildSimulationDeployments(ctx, tx, dep)
	if em != nil {
		return nil, em
	}

	// Is this a multiSimulation?
	if len(childSims) > 1 {
		if em := dep.MarkAsMultiSimParent(tx); em != nil {
			return nil, em
		}
		for _, child := range childSims {
			// Dev note: this call adds the child SimulationDeployment record to DB
			if em := child.MarkAsMultiSimChild(tx, dep); em != nil {
				return nil, em
			}
		}
		simsToLaunch = childSims
	}
	return simsToLaunch, nil
}

// getLaunchableSimulations returns an array of simulations that are ready to be launched
// If it is a multisim, it will return all the child simulations.
// In any other cases, will return a single simulation.
func (s *Service) getLaunchableSimulations(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) ([]SimulationDeployment, *ign.ErrMsg) {
	var deps []SimulationDeployment
	if dep.isMultiSimParent() {
		childsims, err := GetChildSimulationDeployments(tx, dep, simPending, simPending)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
		}
		deps = append(deps, *childsims...)
	} else {
		deps = append(deps, *dep)
	}
	return deps, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// StartMultiSimStatusUpdater initialize the routine that will update the
// DeploymentStatus and ErrorStatus of Parents in Multi-simulations. The statuses
// will be updated based on the status of their children.
func (s *Service) StartMultiSimStatusUpdater() {
	// bind a specific logger to the routing
	newLogger := logger(s.baseCtx).Clone("multisim-status-updater")
	ctx := ign.NewContextWithLogger(s.baseCtx, newLogger)

	s.multisimStatusUpdater = time.NewTicker(20 * time.Second)
	s.multisimStatusUpdaterDone = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-s.multisimStatusUpdaterDone:
				newLogger.Info("MultiSim Parent Status Updater is done.")
				return
			case <-s.multisimStatusUpdater.C:
				s.updateMultiSimStatuses(ctx, s.DB)
			}
		}
	}()
}

// StopMultiSimStatusUpdater stops the update of MultiSim Parents status process.
func (s *Service) StopMultiSimStatusUpdater() {
	s.multisimStatusUpdater.Stop()
	s.multisimStatusUpdaterDone <- true
}

func (s *Service) updateMultiSimStatuses(ctx context.Context, tx *gorm.DB) {
	logger(ctx).Debug("Updating the Statuses of MultiSim Parents...")
	parents, err := GetParentSimulationDeployments(tx, simPending, simTerminatingInstances,
		[]ErrorStatus{simErrorWhenInitializing, simErrorWhenTerminating})
	if err != nil {
		logger(ctx).Error("Error while trying to get Simulation Parents from DB", err)
		return
	}

	// Compute and set the status of each Parent based on its children
	for _, p := range *parents {
		if em := p.updateCompoundStatuses(tx); em != nil {
			logger(ctx).Error("Error computing and updating compound status for Parent: "+*p.GroupID, err)
		}
		s.applications[*p.Application].updateMultiSimStatuses(ctx, tx, s.userAccessor, &p)
	}
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// StartExpiredSimulationsCleaner initialize the routine that will check for expired
// simulations.
func (s *Service) StartExpiredSimulationsCleaner() {
	// bind a specific logger to the cleaner
	newLogger := logger(s.baseCtx).Clone("expired-simulations-cleaner")
	cleanerCtx := ign.NewContextWithLogger(s.baseCtx, newLogger)

	// We check for expired simulations each minute
	s.expiredSimulationsTicker = time.NewTicker(time.Minute)
	s.expiredSimulationsDone = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-s.expiredSimulationsDone:
				newLogger.Info("Expired Simulations Cleaner is done.")
				return
			case <-s.expiredSimulationsTicker.C:
				_ = s.checkForExpiredSimulations(cleanerCtx)
			}
		}
	}()
}

// StopExpiredSimulationsCleaner stops the cleaner process
func (s *Service) StopExpiredSimulationsCleaner() {
	s.expiredSimulationsTicker.Stop()
	s.expiredSimulationsDone <- true
}

// checkForExpiredSimulations is an internal helper that tests all the runningSimulations
// to check if they were alive more than expected, and in that case, schedules their termination.
func (s *Service) checkForExpiredSimulations(ctx context.Context) error {

	logger(ctx).Debug("Checking for expired simulations...")
	s.lockRunningSimulations.RLock()
	defer s.lockRunningSimulations.RUnlock()

	for groupID := range s.runningSimulations {
		rs := s.runningSimulations[groupID]

		if rs.IsExpired() || rs.Finished {
			dep, err := GetSimulationDeployment(s.DB, groupID)
			if err != nil {
				logger(ctx).Error("Error while trying to get Simulation from DB: "+groupID, err)
				continue
			}

			// Add a 'stop simulation' request to the Terminator Jobs-Pool.
			if err := s.scheduleTermination(ctx, s.DB, dep); err != nil {
				logger(ctx).Error("Error while trying to schedule automatic termination of Simulation: "+groupID, err)
			} else {
				reason := "expired"
				if rs.Finished {
					reason = "finished"
				}
				logger(ctx).Info(fmt.Sprintf("Scheduled automatic termination of %s simulation: %s", reason, groupID))
			}
		}
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// RegisterSchedulableTasks gets the tasks from each application and schedules them on the internal scheduler.
var RegisterSchedulableTasks = func(s *Service, ctx context.Context, tx *gorm.DB) {
	for app := range s.GetApplications() {
		for _, task := range s.applications[app].GetSchedulableTasks(ctx, s, tx) {
			s.scheduler.DoAt(task.Fn, task.Date)
		}
	}
}

// LaunchSimulation receives a simulation deployment as an argument and pushes it to the launch queue.
var LaunchSimulation = func(s *Service, ctx context.Context,
	tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {

	// Pre-hook
	if em := s.applications[*dep.Application].ValidateSimulationLaunch(ctx, tx, dep); em != nil {
		return em
	}

	// Process
	groupID := *dep.GroupID
	s.queueLaunchRequest(groupID)
	return nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// countPods is a test function connects to k8 master and returns the number of pods. It is a helper
// used to quickly check a valid connection to the cluster.
func (s *Service) countPods(ctx context.Context, user *users.User) (interface{}, *ign.ErrMsg) {

	// Only system admins
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	pods, err := s.clientset.CoreV1().Pods(s.cfg.KubernetesNamespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	a := fmt.Sprintf("There are %d pods in the cluster", len(pods.Items))
	logger(ctx).Debug(a)
	return &a, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// workerStartSimulation is a thread pool worker that invokes the startSimulation.
func (s *Service) workerStartSimulation(payload interface{}) {
	groupID, ok := payload.(string)
	if !ok {
		return
	}

	// bind a specific logger to the worker
	reqID := fmt.Sprintf("worker-start-sim-%s", groupID)
	newLogger := logger(s.baseCtx).Clone(reqID)
	workerCtx := ign.NewContextWithLogger(s.baseCtx, newLogger)

	newLogger.Info("Worker about to invoke StartSimulation for groupID: " + groupID)

	simDep, err := GetSimulationDeployment(s.DB, groupID)
	if err != nil {
		logger(workerCtx).Error(fmt.Sprintf("startSimulation - %v", err))
		return
	}

	res, em := s.startSimulation(workerCtx, s.DB, simDep)
	if res == launcherRelaunchNeeded {
		s.requeueSimulation(simDep)
	}
	s.notify(PoolStartSimulation, groupID, res, em)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// workerTerminateSimulation is a thread pool worker that invokes the shutdownSimulation.
func (s *Service) workerTerminateSimulation(payload interface{}) {
	groupID, ok := payload.(string)
	if !ok {
		return
	}
	// bind a specific logger to the worker-
	reqID := fmt.Sprintf("worker-finish-sim-%s", groupID)
	newLogger := logger(s.baseCtx).Clone(reqID)
	workerCtx := ign.NewContextWithLogger(s.baseCtx, newLogger)

	newLogger.Info("Worker about to invoke ShutdownSimulation for groupID: " + groupID)
	res, em := s.shutdownSimulation(workerCtx, s.DB, groupID)
	s.notify(PoolShutdownSimulation, groupID, res, em)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// workerErrorHandler is a pool worker responsible of handling errors.
func (s *Service) workerErrorHandler(payload interface{}) {
	groupID, ok := payload.(string)
	if !ok {
		return
	}
	// bind a specific logger to the worker
	reqID := fmt.Sprintf("worker-error-handler-%s", groupID)
	newLogger := logger(s.baseCtx).Clone(reqID)
	workerCtx := ign.NewContextWithLogger(s.baseCtx, newLogger)

	newLogger.Info("Worker about to (try to) handle error for groupID: " + groupID)
	dep, err := GetSimulationDeployment(s.DB, groupID)
	if err != nil {
		logMsg := fmt.Sprintf("workerErrorHandler - Error getting SimulationDeployment from DB for GroupID [%s]", groupID)
		newLogger.Error(logMsg, err)
		return
	}

	if *dep.ErrorStatus == string(simErrorWhenInitializing) {
		res, em := s.rollbackFailedLaunch(workerCtx, s.DB, dep)
		s.notify(PoolRollbackFailedLaunch, groupID, res, em)
	} else if *dep.ErrorStatus == string(simErrorWhenTerminating) {
		res, em := s.completeFailedTermination(workerCtx, s.DB, dep)
		s.notify(PoolCompleteFailedTermination, groupID, res, em)
	}
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// registerError fires a new request for an async error handling process. When
// this function is invoked, the involved simulationDeployment will be marked with
// an error status , and the error handling thread pool will be notified.
func (s *Service) registerError(ctx context.Context, tx *gorm.DB, simDep *SimulationDeployment, st ErrorStatus) *ign.ErrMsg {
	if em := simDep.setErrorStatus(tx, st); em != nil {
		return em
	}
	s.queueErrorHandlerRequest(*simDep.GroupID)

	return nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// 	Service API
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// StartSimulationAsync spawns a task to start a simulation
func (s *Service) StartSimulationAsync(ctx context.Context,
	tx *gorm.DB, createSim *CreateSimulation, user *users.User) (interface{}, *ign.ErrMsg) {

	// TODO: whether a user can read or write to an Organization is defined at ign-fuel's casbin db.
	// Not on cloudsim's casbin. Note: Casbin "caches" the data to avoid accesing the DB all the time.
	// So, a couple of options:
	// 1) allow local casbin to access fuel's casbin db (read-only)
	// and find a way to refresh local casbin keep it in sync with remote DB.
	// To refresh local cache, Casbin uses Watchers (eg. time based watchers, or etcd watcher, etc).
	// 2) Make Users a separate "server" with a REST api for these queries. The problem with this
	// option is that we will need to wait for an http call to return.
	// We are currently using Option 1, with a time-based refresh of permissions.
	// We should add a legend to UI that says "It can take up to 30 seconds to
	// populate new Team memberships and permissions to all servers"

	// Verify and set the owner
	owner := createSim.Owner
	if owner == "" {
		owner = *user.Username
	} else {
		// VerifyOwner checks to see if the 'owner' arg is an organization or a user. If the
		// 'owner' is an organization, it verifies that the given 'user' arg has the expected
		// permission in the organization. If the 'owner' is a user, it verifies that the
		// 'user' arg is the same as the owner.
		if ok, em := s.userAccessor.VerifyOwner(owner, *user.Username, per.Read); !ok {
			return nil, em
		}
	}

	private := true
	if createSim.Private != nil {
		private = *createSim.Private
	}

	isAdmin := s.userAccessor.IsSystemAdmin(*user.Username)

	stopOnEnd := false
	// Only system admins can request instances to stop on end
	if createSim.StopOnEnd != nil && isAdmin {
		stopOnEnd = *createSim.StopOnEnd
	}

	// Create and assign a new GroupID
	groupID := uuid.NewV4().String()

	// Create the SimulationDeployment record in DB. Set initial status.
	creator := *user.Username
	imageStr := SliceToStr(createSim.Image)
	simDep, err := NewSimulationDeployment()
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	simDep.Owner = &owner
	simDep.Name = &createSim.Name
	simDep.Creator = &creator
	simDep.Private = &private
	simDep.StopOnEnd = &stopOnEnd
	simDep.Platform = &createSim.Platform
	simDep.Application = &createSim.Application
	simDep.Image = &imageStr
	simDep.GroupID = &groupID
	simDep.DeploymentStatus = simPending.ToPtr()
	simDep.Extra = createSim.Extra
	simDep.ExtraSelector = createSim.ExtraSelector
	simDep.Robots = createSim.Robots
	simDep.Held = false

	// Set the maximum simulation expiration time.
	validFor := s.getMaxDurationForSimulation(ctx, tx, simDep)
	validForStr := validFor.String()
	simDep.ValidFor = &validForStr

	if err := tx.Create(simDep).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// Set held state if the user is not a sysadmin and the simulations needs to be held
	if !isAdmin && s.applications[*simDep.Application].simulationIsHeld(ctx, tx, simDep) {
		err := simDep.UpdateHeldStatus(tx, true)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
		}

		// Check if the simulation is a submission to a competition circuit.
		// If that's the case, the previous submission should be marked as superseded.
		if IsCompetitionCircuit(*simDep.ExtraSelector) {
			err = MarkPreviousSubmissionsSuperseded(tx, *simDep.GroupID, *simDep.Owner, *simDep.ExtraSelector)
			if err != nil {
				return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
			}
		}
	}



	// Set read and write permissions to owner (eg, the team) and to the Application
	// organizing team (eg. subt).
	if em := s.bulkAddPermissions(groupID, []per.Action{per.Read, per.Write}, owner, *simDep.Application); em != nil {
		return nil, em
	}

	// Sanity check: check for maximum number of allowed simultaneous simulations per Owner.
	// Also allow Applications to provide custom validations.
	// Dev note: in this case we check 'after' creating the record in the DB to make
	// sure that in case of a race condition then both records are added with pending state
	// and one of those (or both) can be rejected immediately.
	if em := s.checkValidNumberOfSimulations(ctx, tx, simDep); em != nil {
		// In case of error we delete the simulation request from DB and exit.
		tx.Model(simDep).Update(SimulationDeployment{
			DeploymentStatus: simRejected.ToPtr(),
			ErrorStatus:      simErrorRejected.ToStringPtr(),
		}).Delete(simDep)
		return nil, em
	}

	// By default, we launch a single simulation from a createSimulation request.
	// But we also allow specific ApplicationTypes (eg. SubT) to spawn multiple simulations
	// from a single request. When that happens, we call those "child simulations"
	// and they will be grouped by the same parent simulation's groupID.
	simsToLaunch, em := s.prepareSimulations(ctx, tx, simDep)
	if em != nil {
		return nil, em
	}

	// Add a 'launch simulation' request to the Launcher Jobs-Pool
	for _, sim := range simsToLaunch {
		groupID := *sim.GroupID
		logger(ctx).Info("StartSimulationAsync about to submit launch task for groupID: " + groupID)
		if err := LaunchSimulation(s, ctx, tx, sim); err != nil {
			logger(ctx).Error(fmt.Sprintf("StartSimulationAsync -- Cannot launch simulation: %s", err.Msg))
		}
	}

	return simDep, nil
}

// MarkPreviousSubmissionsSuperseded marks a set of submissions with the simSuperseded status.
func MarkPreviousSubmissionsSuperseded(tx *gorm.DB, groupID, owner, circuit string) error {
	return tx.Model(&SimulationDeployment{}).
		Where("group_id NOT LIKE ?", fmt.Sprintf("%s%%", groupID)).
		Where("owner = ?", owner).
		Where("extra_selector = ?", circuit).
		Where("held = true").
		Update("deployment_status", simSuperseded.ToInt()).Error
}

// LaunchSimulationAsync launches a simulation that is currently being held by cloudsim.
func (s *Service) LaunchSimulationAsync(ctx context.Context, tx *gorm.DB,
	groupID string, user *users.User) (interface{}, *ign.ErrMsg) {

	if !s.userAccessor.IsSystemAdmin(*user.Username) {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	dep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	if dep.IsRunning() {
		err := errors.New("Cannot launch a running simulation")
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	if err := s.DeployHeldSimulation(ctx, tx, dep); err != nil {
		return nil, err
	}

	return dep, nil
}

// RestartSimulationAsync re-launches a single (children) simulation that previosly
// finished with an error status.
func (s *Service) RestartSimulationAsync(ctx context.Context, tx *gorm.DB,
	groupID string, user *users.User) (interface{}, *ign.ErrMsg) {

	logger(ctx).Info("RestartSimulationAsync requested for groupID: " + groupID)

	mainDep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Is the user authorized to restart the simulation? Only application or system admins
	if ok, em := s.userAccessor.CanPerformWithRole(mainDep.Application, *user.Username, per.Admin); !ok {
		return nil, em
	}

	// Sanity checks
	if mainDep.isMultiSimParent() {
		err := errors.New("Cannot restart a MultiSim parent. Only children simulations")
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Check the simulation is not running already
	if mainDep.IsRunning() {
		err := errors.New("Cannot restart a running simulation")
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Create a clone of the original child simulation and mark it as Pending
	clone := mainDep.Clone()
	clone.DeploymentStatus = simPending.ToPtr()
	clone.ErrorStatus = nil
	clone.DeletedAt = nil
	clone.StoppedAt = nil
	// Update the max runtime limit in case the server configuration was updated
	clone.ValidFor = sptr(s.getMaxDurationForSimulation(ctx, tx, clone).String())
	// Reset the processed field to allow processing when simulations end
	clone.Processed = false

	// Find out if the old simulation was also a "retry" and get its retry number
	const retryStr = "-r-"
	retryNum := 1
	parts := strings.Split(*clone.GroupID, retryStr)
	baseGroupID := parts[0]
	// if the Split resulted in more than one slice then it was a retry
	if len(parts) > 1 {
		numStr := parts[1]
		if retryNum, err = strconv.Atoi(numStr); err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
		}
		retryNum++
	}
	clone.GroupID = sptr(fmt.Sprintf("%s%s%d", baseGroupID, retryStr, retryNum))

	// Save a new row with the clone/retry
	if err := tx.Create(clone).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// Set read and write permissions to owner (eg, the team) and to the Application
	// organizing team (eg. subt).
	if em := s.bulkAddPermissions(*clone.GroupID, []per.Action{per.Read, per.Write}, *clone.Owner, *clone.Application); em != nil {
		return nil, em
	}

	// Invalidate the old run (soft delete it)
	if err := tx.Delete(mainDep).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbDelete, err)
	}
	// Allow the specific Application to invalidate the individual old child run as well.
	// (eg. soft delete its score)
	s.applications[*mainDep.Application].invalidateSimulation(ctx, tx, mainDep)

	// If the restarted sim is a child simulation, then we need to mark the Parent
	// as 'Pending' again so we can compute its aggregated status.
	if mainDep.isMultiSimChild() {
		parentSim, err := GetParentSimulation(tx, mainDep)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
		}
		if err := tx.Model(&parentSim).Update(SimulationDeployment{
			DeploymentStatus: simPending.ToPtr(),
		}).Error; err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
		}
		// This is needed instead of 'nil' to ensure the Update will overrite with
		// NULL an existing value.
		// https://github.com/jinzhu/gorm/issues/1073
		if err := tx.Model(&parentSim).Update("error_status", gorm.Expr("NULL")).Error; err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
		}
	}

	// commit the DB transaction
	// Note: we commit the TX here on purpose, to be able to detect DB errors in
	// advance. And before sending the relaunch to the pending queue.
	if err := tx.Commit().Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	// Add a new 'launch simulation' request to the Launcher Jobs-Pool
	logger(ctx).Info("RestartSimulationAsync about to submit task to re-launch groupID: " + *clone.GroupID)
	if err := LaunchSimulation(s, ctx, tx, clone); err != nil {
		logger(ctx).Error(fmt.Sprintf("rebuildState -- Cannot launch simulation: %s", err.Msg))
	}

	return clone, nil
}

// bulkAddPermissions adds multiple permissions to multiple owners to access a resource.
func (s *Service) bulkAddPermissions(resID string, permissions []per.Action, owners ...string) *ign.ErrMsg {
	for _, o := range owners {
		for _, p := range permissions {
			if ok, em := s.userAccessor.AddResourcePermission(o, resID, p); !ok {
				return em
			}
		}
	}
	return nil
}

// checkValidNumberOfSimulations checks if the given owner hasn't gone beyond the
// maximum number of allowed concurrent simulations.
func (s *Service) checkValidNumberOfSimulations(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {

	// Default sanity check: check for maximum number of allowed "simultaneous" simulations per Owner.
	// Dev note: we check 'after' creating the record in the DB to make
	// sure that in case of a race condition then both records are added with pending state
	// and one of those (or both) can be rejected immediately.
	owner := *dep.Owner
	app := *dep.Application

	limit := s.cfg.MaxSimultaneousSimsPerOwner
	if limit != 0 {
		runningSims, err := s.getRunningSimulationDeploymentsByOwner(tx, owner)
		if err != nil {
			logger(ctx).Info("Failed to get running simulations by owner")
			return NewErrorMessageWithBase(
				ign.ErrorUnexpected,
				fmt.Errorf("failed to get running simulations by owner %w", err),
			)
		}
		if len(*runningSims) > limit {
			logger(ctx).Info(fmt.Sprintf(
				"Owner [%s] has reached the simultaneous simulations limit [%d]. Running simulations [%v]",
				owner, limit, *runningSims))

			return NewErrorMessageWithBase(
				ErrorOwnerSimulationsLimitReached,
				fmt.Errorf("cannot request new simulation, owner [%s] has reached the simultaneous simulations limit [%d]", owner, limit),
			)
		}
	}

	// Now allow registered Application to provide custom validations
	if em := s.applications[app].checkValidNumberOfSimulations(ctx, s, tx, dep); em != nil {
		return em
	}

	// All OK
	return nil
}

// GetRemainingSubmissions returns the number of remaining submissions for an
// owner in a circuit.
func (s *Service) GetRemainingSubmissions(ctx context.Context, tx *gorm.DB, user *users.User, circuit *string,
	owner *string) (interface{}, *ign.ErrMsg) {
	// TODO: whether a user can read or write to an Organization is defined
	//  at ign-fuel's casbin db. See StartSimulationAsync for more information.

	// Verify and set the owner
	if *owner == "" {
		owner = user.Username
	} else {
		// VerifyOwner checks to see if the 'owner' arg is an organization or a user. If the
		// 'owner' is an organization, it verifies that the given 'user' arg has the expected
		// permission in the organization. If the 'owner' is a user, it verifies that the
		// 'user' arg is the same as the owner.
		if ok, em := s.userAccessor.VerifyOwner(*owner, *user.Username, per.Read); !ok {
			return nil, em
		}
	}

	remaining, err := getRemainingSubmissions(tx, *circuit, *owner)
	if err != nil {
		return 0, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	if remaining != nil {
		// Ensure no negative values are returned
		val := Max(0, *remaining)
		remaining = &val
	}

	return remaining, nil
}

// CustomRuleList returns a paginated list of circuit custom rules.
// This operation can only be performed by a system administrator and team administrators.
func (s *Service) CustomRuleList(ctx context.Context, p *ign.PaginationRequest, tx *gorm.DB, user *users.User,
	application *string, circuit *string, owner *string, ruleType *CustomRuleType) (*CircuitCustomRules,
	*ign.PaginationResult, *ign.ErrMsg) {
	// Restrict access to application and system admins
	if ok, _ := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Admin); !ok {
		return nil, nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	// Create the DB query
	var rules CircuitCustomRules
	q := tx.Model(&CircuitCustomRule{})

	if circuit != nil {
		q = q.Where("circuit = ?", *circuit)
	}
	if owner != nil {
		q = q.Where("owner = ?", *owner)
	}
	if ruleType != nil {
		q = q.Where("rule_type = ?", *ruleType)
	}

	pagination, err := ign.PaginateQuery(q, &rules, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &rules, pagination, nil
}

// SetCustomRule creates or updates a custom rule for an owner in a circuit.
// This operation can only be performed by a system administrator and team administrators.
// TODO System and application admins are able to create rules for invalid owners because admin privileges override
//  invalid owner errors.
func (s *Service) SetCustomRule(ctx context.Context, tx *gorm.DB, user *users.User, application *string,
	circuit *string, owner *string, ruleType *CustomRuleType, value *string) (*CircuitCustomRule, *ign.ErrMsg) {
	// Restrict access to application and system admins
	if ok, _ := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Admin); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	// Validate the new/updated rule model instance
	newRule := CircuitCustomRule{
		Circuit:  circuit,
		Owner:    owner,
		RuleType: *ruleType,
		Value:    *value,
	}
	if em := ValidateStruct(newRule); em != nil {
		return nil, em
	}

	// Create or update the rule
	var rule CircuitCustomRule
	tx.Where(&CircuitCustomRule{
		Circuit:  circuit,
		Owner:    owner,
		RuleType: *ruleType,
	}).
		Assign(&newRule).
		FirstOrCreate(&rule)

	return &rule, nil
}

// DeleteCustomRule deletes a custom rule for an owner in a circuit.
// This operation can only be performed by a system administrator and team administrators.
// TODO System and team admins are able to delete rules for invalid owners because admin privileges override
//  invalid owner errors.
func (s *Service) DeleteCustomRule(ctx context.Context, tx *gorm.DB, user *users.User, application *string,
	circuit *string, owner *string, ruleType *CustomRuleType) (interface{}, *ign.ErrMsg) {
	// Restrict access to application and system admins
	if ok, _ := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Admin); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	rule := &CircuitCustomRule{}

	if err := tx.Where(&CircuitCustomRule{
		Circuit:  circuit,
		Owner:    owner,
		RuleType: *ruleType,
	}).
		First(rule).
		Delete(CircuitCustomRule{}).
		Error; err != nil {
		errMsg := fmt.Sprintf("Attempted to delete nonexistent rule [%s] entry for Owner [%s].", string(*ruleType), *owner)
		logger(ctx).Debug(errMsg)
		return nil, NewErrorMessageWithBase(ErrorRuleForOwnerNotFound, errors.Errorf(errMsg))
	}

	return rule, nil
}

// getMaxDurationForSimulation returns the max duration for a simulation, based
// on the chosen Application (eg. SubT).
func (s *Service) getMaxDurationForSimulation(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) time.Duration {

	// Find the max duration for simulations based on the chosen Application.
	maxDuration := s.applications[*dep.Application].getMaxDurationForSimulation(ctx, tx, dep)
	if maxDuration == 0 {
		// Set the max duration if not specified by the Application
		maxDuration = time.Duration(s.cfg.MaxDurationForSimulations) * time.Minute
	}
	return maxDuration
}

// StartSimulation is the main func to launch a new simulation.
// IMPORTANT: This function is invoked in a separate thread, from a Launcher Worker thread.
// @return: it can return a (launcherRelaunchNeeded) "relaunch" string value as result, which means the
// pool worker will send the simulation again to the Pending queue.
func (s *Service) startSimulation(ctx context.Context, tx *gorm.DB,
	simDep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	groupID := *simDep.GroupID
	logger(ctx).Info("startSimulation running for groupID: " + groupID)

	// Sanity checks

	// Check the simulation has the correct status
	if em := simDep.assertSimDepStatus(simPending); em != nil {
		logger(ctx).Warning(fmt.Sprintf("startSimulation - Invalid simulation status: %d", *simDep.DeploymentStatus))
		return nil, em
	}

	// Cannot launch a gz simulation from a MultiSim Parent. Only for child simulations
	if simDep.isMultiSimParent() {
		err := errors.New("Cannot launch a gz simulation from a MultiSim Parent. Only for child simulations")
		logger(ctx).Error(fmt.Sprintf("startSimulation - %v", err))
		tx.Model(simDep).Update(SimulationDeployment{
			DeploymentStatus: simRejected.ToPtr(),
			ErrorStatus:      simErrorRejected.ToStringPtr(),
		})
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
	}

	// Sanity check: Check if the Parent simulation doesn't have an error status already.
	// In that case we stop launching the child simulation.
	if simDep.isMultiSimChild() {
		parentSim, err := GetParentSimulation(tx, simDep)
		if err != nil {
			logger(ctx).Error(fmt.Sprintf("startSimulation - %v", err))
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
		if parentSim.ErrorStatus != nil {
			err := errors.New("Cannot launch a children simulation when its parent has ErrorStatus already")
			logger(ctx).Error(fmt.Sprintf("startSimulation - %v", err))
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
	}

	// Everything OK. Log the launch details (to Rollbar)
	str, err := simDep.toJSON()
	if err != nil {
		logger(ctx).Error(fmt.Sprintf("startSimulation - %v", err))
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	logger(ctx).Warning(fmt.Sprintf("startSimulation - SimulationDeployment to launch: [%s]", *str))

	// Move to 'launching nodes' status
	if em := simDep.updateSimDepStatus(tx, simLaunchingNodes); em != nil {
		logger(ctx).Error(fmt.Sprintf("startSimulation - %v", err))
		return nil, em
	}

	tstart := time.Now()

	// Run the following as a block, and in case of an error set the simDep's ErrorStatus field.
	_, em := func() (interface{}, *ign.ErrMsg) {

		// NOTE: This call will block until nodes are created.
		nodeSelectorGroupID, em := s.hostsSvc.launchNodes(ctx, tx, simDep)
		if em != nil {
			return nil, em
		}
		// Wait until Nodes are ready before updating the status.
		timeout := time.Duration(s.cfg.NodeReadyTimeoutSeconds) * time.Second
		if err := WaitForNodesReady(ctx, s.clientset, s.cfg.KubernetesNamespace, *nodeSelectorGroupID, timeout); err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}
		timeTrack(ctx, tstart, "startSimulation - launchNodes")

		// switch to next status
		if em := simDep.updateSimDepStatus(tx, simLaunchingPods); em != nil {
			return nil, em
		}

		logger(ctx).Info("startSimulation - about to launch pods for groupID: " + groupID)
		// After getting the nodes ready we can schedule the Pods.
		if _, em := s.launchGazeboServerInGroup(ctx, tx, groupID, simDep); em != nil {
			return nil, em
		}

		// Wait until Pods are actually running and ready before continuing.
		// TODO: wait until Gazebo server is actually running and ready to receive data.
		// Idea, use kubernetes's readinessProbes for that.
		groupIDLabel := getPodLabelSelectorForSearches(groupID)
		timeout = time.Duration(s.cfg.PodReadyTimeoutSeconds) * time.Second
		if err := WaitForPodsReady(ctx, s.clientset, s.cfg.KubernetesNamespace, groupIDLabel, timeout); err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Create, err)
		}

		// Register a new RunningSimulation
		if err := s.createRunningSimulation(ctx, tx, simDep); err != nil {
			return nil, NewErrorMessageWithBase(ErrorCreatingRunningSimulationNode, err)
		}

		timeTrack(ctx, tstart, "startSimulation - launchGazeboServerInGroup")

		// Finally, change the simulation status to Running
		if em := simDep.updateSimDepStatus(tx, simRunning); em != nil {
			return nil, em
		}
		return simDep, nil
	}()
	if em != nil {
		if em.ErrCode == ign.ErrorLaunchingCloudInstanceNotEnoughResources && s.AllowRequeuing {
			// If the EC2 instances could not be started due to insufficient
			// instances available then requeue this simulation
			return launcherRelaunchNeeded, em
		}
		// Otherwise mark the simulation as failed
		logMsg := fmt.Sprintf("startSimulation - error in startSimulation for groupid [%s]. Error: %v", groupID, em)
		logger(ctx).Error(logMsg, em)
		timeTrack(ctx, tstart, "startSimulation - time tracker until error")
		s.registerError(ctx, tx, simDep, simErrorWhenInitializing)
		return nil, em
	}

	logger(ctx).Info("startSimulation - successfully launched groupID: " + groupID)
	return simDep, nil
}

// createRunningSimulation is a helper func used to create and register a new RunningSimulation.
func (s *Service) createRunningSimulation(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) error {
	worldStatsTopic, maxSimSeconds, err := s.getGazeboWorldStatsTopicAndLimit(ctx, tx, dep)
	if err != nil {
		return err
	}

	// TODO: warmup topic is not a generic concept as it is specific of SubT. Need to move to a SubT custom code.
	// TODO: Consider allowing Applications to configure the RunningSimulation instance.
	worldWarmupTopic, err := s.getGazeboWorldWarmupTopic(ctx, tx, dep)
	if err != nil {
		return err
	}

	t, err := s.setupRunningSimulationTransportLayer(dep)
	if err != nil {
		return err
	}

	rs, err := NewRunningSimulation(ctx, dep, t, worldStatsTopic, worldWarmupTopic, maxSimSeconds)
	if err != nil {
		return err
	}
	s.addRunningSimulation(rs)
	return nil
}

// setupRunningSimulationTransportLayer initializes a new transport layer for the given simulation deployment.
func (s *Service) setupRunningSimulationTransportLayer(dep *SimulationDeployment) (ignws.PubSubWebsocketTransporter, error) {
	host := s.applications[*dep.Application].getSimulationWebsocketHost()
	path := s.applications[*dep.Application].getSimulationWebsocketPath(*dep.GroupID)

	var t ignws.PubSubWebsocketTransporter
	var err error
	for i := 1; i <= 10; i++ {
		t, err = newTransporter(host, path, *dep.AuthorizationToken, s.cfg.IsTest)
		if err == nil {
			break
		}
		// i * 10s
		Sleep(time.Duration(i*10) * time.Second)
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

// newTransporter returns a new ign websocket transport.
// If isTest is set to true, it will return the default transport test mock.
func newTransporter(host, path, token string, isTest bool) (ignws.PubSubWebsocketTransporter, error) {
	if isTest {
		if globals.TransportTestMock == nil {
			return nil, errors.New("mock for testing transport not initialized")
		}
		return globals.TransportTestMock, nil
	}
	return ignws.NewIgnWebsocketTransporter(host, path, transport.WebsocketSecureScheme, token)
}

func (s *Service) requeueSimulation(simDep *SimulationDeployment) *ign.ErrMsg {
	// Revert the simulation deployment status to Pending
	if em := simDep.updateSimDepStatus(s.DB, simPending); em != nil {
		return em
	}
	// Wait a little time and requeue the simulation
	Sleep(time.Minute)
	s.queueLaunchRequest(*simDep.GroupID)

	return nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// ShutdownSimulationAsync spawns a task to finish a simulation.
func (s *Service) ShutdownSimulationAsync(ctx context.Context, tx *gorm.DB,
	groupID string, user *users.User) (interface{}, *ign.ErrMsg) {

	logger(ctx).Info("ShutdownSimulationAsync requested for groupID: " + groupID)

	// Is the user authorized to shutdown the simulation? First we check generic
	// permissions. Then we allow specific Applications to reject requests as well.
	if ok, em := s.userAccessor.IsAuthorizedForResource(*user.Username, groupID, per.Read); !ok {
		return nil, em
	}

	mainDep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Allow specific Application to reject the permissions too.
	if ok, em := s.applications[*mainDep.Application].checkCanShutdownSimulation(ctx, s, tx, mainDep, user); !ok {
		return nil, em
	}

	// Check the simulation has the correct status
	if em := mainDep.assertSimDepStatus(simRunning); em != nil {
		return nil, em
	}

	var depsToTerminate *SimulationDeployments

	// Is this a multiSim?
	if mainDep.isMultiSimParent() {
		// Get all child simulations that have status simRunning.
		depsToTerminate, err = GetChildSimulationDeployments(tx, mainDep, simRunning, simRunning)
		if err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
		}
	} else {
		depsToTerminate = &SimulationDeployments{*mainDep}
	}

	// Schedule the terminations
	for _, d := range *depsToTerminate {
		if em := s.scheduleTermination(ctx, tx, &d); em != nil {
			return nil, em
		}
	}
	return mainDep, nil
}

// scheduleTermination marks a simulation as "to be terminated" and queues it
// into the Termination Pool.
func (s *Service) scheduleTermination(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {
	// Do not continue if the simulation has already started termination
	if *dep.DeploymentStatus >= int(simTerminateRequested) {
		depStatus := DeploymentStatus(*dep.DeploymentStatus)
		logger(ctx).Warning(fmt.Sprintf(
			"Attempted to terminate simulation [%s] with status %s.", *dep.GroupID, depStatus.String(),
		))
		return nil
	}

	if em := dep.updateSimDepStatus(tx, simTerminateRequested); em != nil {
		return em
	}
	// Add a 'stop simulation' request to the Terminator Jobs-Pool.
	s.queueShutdownRequest(*dep.GroupID)

	return nil
}

// internalShutdownSimulation is an internal helper function used to free
// resources. It is invoked by the normal shutdown simulation and by the
// error handlers.
func (s *Service) internalShutdownSimulation(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment, logPrefix string) (interface{}, *ign.ErrMsg) {

	if dep.isMultiSimParent() {
		// Parents aren't real simulations. They just spawn child simulations.
		// So they cannot be shutdown. If a Parent Sim is here then that's an error.
		err := errors.New("Cannot shutdown a MultiSim Parent. Only child simulations")
		return nil, ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
	}

	groupID := *dep.GroupID

	rs := s.removeRunningSimulation(groupID)
	if rs != nil {
		rs.Free(ctx)
	}

	// Mark the simulation as stopped
	if em := dep.recordStop(tx); em != nil {
		return nil, em
	}

	if *dep.DeploymentStatus == int(simDeletingPods) {
		// Delete the pods, services , etc
		logger(ctx).Info(fmt.Sprintf("%s - about to delete pods for groupID: %s", logPrefix, groupID))
		// It is expected that deleteGazeboServerInGroup will block until Pods and Services
		// were cleanly terminated (including preStop scripts, etc).
		// Blocking is needed because after this call, the host/node will be killed (or reused).
		if em := s.deleteGazeboServerInGroup(ctx, tx, dep); em != nil {
			return nil, em
		}
		// Switch to next state to allow the next block to run
		if em := dep.updateSimDepStatus(tx, simDeletingNodes); em != nil {
			return nil, em
		}
	}

	if *dep.DeploymentStatus == int(simDeletingNodes) {
		// Delete the kubernetes nodes
		_, em := s.hostsSvc.deleteK8Nodes(ctx, tx, *dep.GroupID)
		// We allow for cases where were not used and thus are not found now
		if em != nil && em.ErrCode != int(ErrorLabeledNodeNotFound) {
			return nil, em
		}
		// Switch to next state to allow the next block to run
		if em := dep.updateSimDepStatus(tx, simTerminatingInstances); em != nil {
			return nil, em
		}
	}

	if *dep.DeploymentStatus == int(simTerminatingInstances) {
		// ask the Cloud instances manager to terminate the instances
		_, em := s.hostsSvc.deleteHosts(ctx, tx, dep)
		if em != nil {
			return nil, em
		}
		// Switch to next state, marking is as terminated
		if em := dep.updateSimDepStatus(tx, simTerminated); em != nil {
			return nil, em
		}
	}

	return dep, nil
}

// ShutdownSimulation finishes all resources associated to a cloudsim simulation.
// (eg. Nodes, Hosts, Pods)
// IMPORTANT: this function is invoked in a separate thread, from a Terminator Worker thread.
func (s *Service) shutdownSimulation(ctx context.Context, tx *gorm.DB, groupID string) (interface{}, *ign.ErrMsg) {

	dep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Check the simulation has the correct status
	if em := dep.assertSimDepStatus(simTerminateRequested); em != nil {
		return nil, em
	}

	if em := dep.updateSimDepStatus(tx, simDeletingPods); em != nil {
		return nil, em
	}

	tstart := time.Now()

	// Run the following as a block, and in case of error set the simDep error status field.
	_, em := func() (interface{}, *ign.ErrMsg) {
		return s.internalShutdownSimulation(ctx, tx, dep, "shutdownSimulation")
	}()
	if em != nil {
		// If the simulation failed to upload logs, the EC2 instances will be kept alive to allow
		// admins to manually extract logs.
		if em.ErrCode == int(ErrorFailedToUploadLogs) {
			// Set the simulation error status
			if em := dep.setErrorStatus(tx, simErrorFailedToUploadLogs); em != nil {
				logger(ctx).Error(fmt.Sprintf("Could not update error status to %s.", simErrorFailedToUploadLogs))
			}
			logMsg := "shutdownSimulation - Errors occurred while uploading log files. Resources will not be terminated."
			logger(ctx).Error(logMsg, em)
			timeTrack(ctx, tstart, "shutdownSimulation - time tracker until error")
			// Return without calling the error handler to avoid it from terminating this simulation's resources
			return nil, em
		}
		logMsg := fmt.Sprintf("shutdownSimulation - error in shutdownSimulation for groupid [%s]. Error: %v", *dep.GroupID, em)
		logger(ctx).Error(logMsg, em)
		timeTrack(ctx, tstart, "shutdownSimulation - time tracker until error")
		return nil, em
	}

	logger(ctx).Info("shutdownSimulation - successfully removed groupID: " + groupID)
	timeTrack(ctx, tstart, "shutdownSimulation - Success")
	return dep, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// rollbackFailedLaunch tries to undo a failed launch, to release resources.
// It is invoked by the ErrorHandler worker.
func (s *Service) rollbackFailedLaunch(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	// Need to release the resources. Based on the status reached at launch, set the
	// equivalent status from the shutdown process and start the shutdown.
	var newSt DeploymentStatus
	switch st := *dep.DeploymentStatus; st {
	case int(simLaunchingPods):
		newSt = simDeletingPods
	case int(simLaunchingNodes):
		newSt = simDeletingNodes
	}

	tstart := time.Now()

	// Run the following as a block, and in case of error set the simDep error status field.
	_, em := func() (interface{}, *ign.ErrMsg) {
		if em := dep.updateSimDepStatus(tx, newSt); em != nil {
			return nil, em
		}
		return s.internalShutdownSimulation(ctx, tx, dep, "rollbackFailedLaunch")
	}()
	if em != nil {
		logMsg := fmt.Sprintf("rollbackFailedLaunch - error while doing rollback in groupid [%s]. Marking for Admin review. Error: %v", *dep.GroupID, em)
		logger(ctx).Error(logMsg, em)
		timeTrack(ctx, tstart, "rollbackFailedLaunch - time tracker until error")
		// There was an error during error handling. Marking for Admin Review
		dep.setErrorStatus(tx, simErrorAdminReview)
		return nil, em
	}
	timeTrack(ctx, tstart, "rollbackFailedLaunch - Success")
	return dep, nil
}

// completeFailedTermination tries to finish a failed termination, to release resources.
// It is invoked by the ErrorHandler worker.
func (s *Service) completeFailedTermination(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	tstart := time.Now()

	// Run the following as a block, and in case of error set the simDep error status field.
	_, em := func() (interface{}, *ign.ErrMsg) {
		return s.internalShutdownSimulation(ctx, tx, dep, "completeFailedTermination")
	}()
	if em != nil {
		logMsg := fmt.Sprintf("completeFailedTermination - error while completing failed termination for groupid [%s]. Marking for Admin review. Error: %v", *dep.GroupID, em)
		logger(ctx).Error(logMsg, em)
		timeTrack(ctx, tstart, "completeFailedTermination - time tracker until error")
		// There was an error during error handling. Marking for Admin Review
		dep.setErrorStatus(tx, simErrorAdminReview)
		return nil, em
	}
	timeTrack(ctx, tstart, "completeFailedTermination - Success")
	return dep, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// DeleteNodesAndHostsForGroup starts the shutdown of all the kubernetes nodes
// and associated Hosts (instances) of a given Cloudsim Simulation.
func (s *Service) DeleteNodesAndHostsForGroup(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment, user *users.User) (interface{}, *ign.ErrMsg) {

	// make sure the requesting user has the correct permissions
	// Only admins of the Application team (eg. subt Org) can do that at the moment.
	if ok, em := s.userAccessor.CanPerformWithRole(dep.Application, *user.Username, per.Admin); !ok {
		return nil, em
	}

	if em := dep.updateSimDepStatus(tx, simDeletingNodes); em != nil {
		return nil, em
	}

	return s.internalShutdownSimulation(ctx, tx, dep, "DeleteNodesAndHostsForGroup")
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// getSimulationPodNamePrefix returns the pod name prefix for a simulation
func getSimulationPodNamePrefix(groupID string) string {
	return fmt.Sprintf("sim-%s", groupID)
}

// launchGazeboServerInGroup launches a set of gzserver pods and associated services in the
// given group.
func (s *Service) launchGazeboServerInGroup(ctx context.Context, tx *gorm.DB, groupID string,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	// It is quite important that application's specific launchers do add the following
	// labels to the created Pods / Services.
	podName := getSimulationPodNamePrefix(groupID)
	labels := map[string]string{
		cloudsimTagLabelKey: "true",
		podLabelPodGroup:    podName,
		podLabelKeyGroupID:  groupID,
	}

	// Add the parent's groupID to the labels as well.
	if dep.isMultiSimChild() {
		labels["parent-group-id"] = *dep.ParentGroupID
	}

	// Find the specific Application handler and ask it to launch the app, using
	// the given base labels and groupID.
	return s.applications[*dep.Application].launchApplication(ctx, s, tx, dep, podName, labels)
}

// getPodLabelSelectorForSearches is a helper function to return the full groupID label
// used for searching Pods associated to a groupID.
func getPodLabelSelectorForSearches(groupID string) string {
	return podLabelKeyGroupID + "=" + groupID
}

// deleteGazeboServerInGroup removes an existing gzserver pod and its services from a group.
func (s *Service) deleteGazeboServerInGroup(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) *ign.ErrMsg {
	// Find the specific Application handler and ask it to delete the app, using
	// the given groupID to find the involved pods/services.
	return s.applications[*dep.Application].deleteApplication(ctx, s, tx, dep)
}

// getGazeboWorldStatsTopicAndLimit returns the topic to subscribe to get notifications about the simulation
// state (eg. /world/default/stats) and time, as well as the Maximum allowed
// Sim time seconds (before marking the simulation as expired).
// This request is delegated to the specific application being launched.
func (s *Service) getGazeboWorldStatsTopicAndLimit(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (string, int, error) {
	return s.applications[*dep.Application].getGazeboWorldStatsTopicAndLimit(ctx, tx, dep)
}

// getGazeboWorldWarmupTopic returns the topic to subscribe to get notifications about the simulation.
// finishing the Warmup period (ie. being ready to start).
// This request is delegated to the specific application being launched.
func (s *Service) getGazeboWorldWarmupTopic(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (string, error) {
	return s.applications[*dep.Application].getGazeboWorldWarmupTopic(ctx, tx, dep)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// SimulationDeploymentList returns a paginated list with cloudsim simulations.
// Members of the submitting team can see the list of simulations they submitted.
// Members of the organizing application (eg. members of SubT Organization) can see all
// simulations for that application.
func (s *Service) SimulationDeploymentList(ctx context.Context, p *ign.PaginationRequest,
	tx *gorm.DB, byStatus *DeploymentStatus, invertStatus bool,
	byErrStatus *ErrorStatus, invertErrStatus bool, byCircuit *string, user *users.User,
	application *string, includeChildren bool) (*SimulationDeployments, *ign.PaginationResult, *ign.ErrMsg) {

	// Create the DB query
	var sims SimulationDeployments
	q := tx.Model(&SimulationDeployment{})
	// Return the newest simulations first
	q = q.Order("created_at desc, id", true)

	if application == nil {
		// The user is requesting ALL simulations from all applications. Only system admins can do that.
		if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
			return nil, nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
	} else {
		// Only simulations from the given application (eg. subt).
		q = q.Where("application = ?", *application)
	}

	// Restrict including children to application and system admins
	if ok, _ := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Admin); !ok {
		// Regardless of the value passed as argument, we set it to False if the requestor
		// is neither an application or system admin.
		includeChildren = false
	}

	if !includeChildren {
		// Only Top Level simulations (ie. not child sims from MultiSims)
		q = q.Where("multi_sim != ?", multiSimChild)
	}

	// Filter by status
	if byStatus != nil {
		if invertStatus {
			q = q.Where("deployment_status != ?", int(*byStatus))
		} else {
			q = q.Where("deployment_status = ?", int(*byStatus))
		}
	}

	// Filter by error status
	if byErrStatus != nil {
		if invertErrStatus {
			q = q.Where("error_status != ?", string(*byErrStatus))
		} else {
			q = q.Where("error_status = ?", string(*byErrStatus))
		}
	}

	// Filter by circuit
	// TODO: This is SubT specific and should be moved
	if byCircuit != nil {
		q = q.Where("extra_selector = ?", *byCircuit)
	}

	// If user belongs to the application's main Org, then he can see all simulations.
	// Otherwise, only those simulations created by the user's team.
	if ok, _ := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Member); !ok {
		// filter resources based on privacy setting
		q = s.userAccessor.QueryForResourceVisibility(q, nil, user)
	}

	pagination, err := ign.PaginateQuery(q, &sims, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &sims, pagination, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// GetSimulationDeployment returns a single simulation deployment based on its groupID
func (s *Service) GetSimulationDeployment(ctx context.Context, tx *gorm.DB,
	groupID string, user *users.User) (interface{}, *ign.ErrMsg) {

	// make sure the user has the correct permissions
	if ok, em := s.userAccessor.IsAuthorizedForResource(*user.Username, groupID, per.Read); !ok {
		return nil, em
	}

	var dep *SimulationDeployment
	var err error

	dep, err = GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	var extra *ExtraInfoSubT
	extra, err = ReadExtraInfoSubT(dep)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// If the user is not a system admin, remove the RunIndex and WorldIndex fields.
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		extra.RunIndex = nil
		extra.WorldIndex = nil
	}

	dep.Extra, err = extra.ToJSON()
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	return dep, nil
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// GetSimulationWebsocketAddress returns a live simulation's websocket server address and authorization token.
// If the simulation is not running, an error is returned.
func (s *Service) GetSimulationWebsocketAddress(ctx context.Context, tx *gorm.DB, user *users.User,
	groupID string) (interface{}, *ign.ErrMsg) {

	dep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// Check that the requesting user has the correct permissions
	username := *user.Username
	// Parent simulations are not valid as they do not run simulations directly
	if dep.isMultiSimParent() {
		return nil, ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus)
	}
	// Multisim child simulations can only be accessed by admins
	if dep.isMultiSimChild() && !s.userAccessor.IsSystemAdmin(username) {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	// Check access permissions
	if ok, em := s.userAccessor.IsAuthorizedForResource(username, groupID, per.Read); !ok {
		return nil, em
	}

	// Find the specific Application handler and ask for the live logs
	return s.applications[*dep.Application].getSimulationWebsocketAddress(ctx, s, tx, dep)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// GetSimulationLogsForDownload returns the generated logs from a simulation.
func (s *Service) GetSimulationLogsForDownload(ctx context.Context, tx *gorm.DB,
	user *users.User, groupID string, robotName *string) (*string, *ign.ErrMsg) {

	dep, err := GetSimulationDeployment(tx, groupID)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// make sure the requesting user has the correct permissions
	if ok, em := s.userAccessor.IsAuthorizedForResource(*user.Username, groupID, per.Read); !ok {
		return nil, em
	}

	// Find the specific Application handler and ask it to generate the link to download logs.
	return s.applications[*dep.Application].getSimulationLogsForDownload(ctx, tx, dep, robotName)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// GetSimulationLiveLogs returns the live logs from a simulation.
func (s *Service) GetSimulationLiveLogs(ctx context.Context, tx *gorm.DB, user *users.User, groupID string,
	robotName *string, lines *int64) (interface{}, *ign.ErrMsg) {

	dep, err := GetSimulationDeployment(tx, groupID)

	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	// make sure the requesting user has the correct permissions
	if ok, em := s.userAccessor.IsAuthorizedForResource(*user.Username, groupID, per.Read); !ok {
		return nil, em
	}

	// Find the specific Application handler and ask for the live logs
	return s.applications[*dep.Application].getSimulationLiveLogs(ctx, s, tx, dep, robotName, *lines)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// GetCloudMachineInstances returns a paginated list with all cloud instances.
func (s *Service) GetCloudMachineInstances(ctx context.Context, p *ign.PaginationRequest,
	tx *gorm.DB, byStatus *MachineStatus, invertStatus bool, groupID *string, user *users.User,
	application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg) {

	// check if the requesting user has permission to access machines in the given
	// application. Only members of the Application team (ie. Org) can do that.

	// Dev Note: we assume that each "application" has a corresponding Organization
	// with the same name. Members of that Organization will be considered Admins
	// of the application.
	if ok, em := s.userAccessor.CanPerformWithRole(application, *user.Username, per.Member); !ok {
		return nil, nil, em
	}

	return s.hostsSvc.CloudMachinesList(ctx, p, tx, byStatus, invertStatus, groupID, application)
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// removeRunningSimulation deletes and return a RunningSimulation
func (s *Service) removeRunningSimulation(groupID string) *RunningSimulation {
	s.lockRunningSimulations.Lock()
	defer s.lockRunningSimulations.Unlock()
	rs := s.runningSimulations[groupID]
	delete(s.runningSimulations, groupID)
	return rs
}

// addRunningSimulation registers a new RunningSimulation
func (s *Service) addRunningSimulation(rs *RunningSimulation) {
	s.lockRunningSimulations.Lock()
	defer s.lockRunningSimulations.Unlock()
	s.runningSimulations[rs.GroupID] = rs
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// getRunningSimulationDeploymentsByOwner returns all the top level SimulationDeployments
// requests (multiSimSingle or multiSimParent -- not childs) that can be considered 'running',
// filtered by a given owner. It is used to count the number of active pending / running
// requests a user has made simulteneously.
func (s *Service) getRunningSimulationDeploymentsByOwner(tx *gorm.DB, owner string) (*SimulationDeployments, error) {
	deps, err := GetSimulationDeploymentsByOwner(tx, owner, simPending, simRunning)
	if err != nil {
		return nil, err
	}
	running := make(SimulationDeployments, 0)
	for _, d := range *deps {
		if !d.isMultiSimChild() && d.IsRunning() {
			running = append(running, d)
		}
	}
	return &running, nil
}

// GetCompetitionRobots returns the list of available robot configurations for a competition.
func (s *Service) GetCompetitionRobots(applicationName string) (interface{}, *ign.ErrMsg) {
	return s.applications[applicationName].getCompetitionRobots()
}

// ///////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////

// QueueGetElements returns a paginated list of elements from the launch queue.
// If no page or perPage arguments are passed, it sets those value to 0 and 10 respectively.
func (s *Service) QueueGetElements(ctx context.Context, user *users.User, page, perPage *int) ([]interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	if page == nil {
		page = intptr(0)
	}
	if perPage == nil {
		perPage = intptr(10)
	}
	offset := *page * *perPage
	limit := *perPage
	return s.launchHandlerQueue.Get(&offset, &limit)
}

// QueueCount returns the element count from the launch queue.
func (s *Service) QueueCount(ctx context.Context, user *users.User) (interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.launchHandlerQueue.Count(), nil
}

// QueueMoveElementToFront moves an element by the given groupID to the front of the queue.
func (s *Service) QueueMoveElementToFront(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.launchHandlerQueue.MoveToFront(groupID)
}

// QueueMoveElementToBack moves an element by the given groupID to the back of the queue.
func (s *Service) QueueMoveElementToBack(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.launchHandlerQueue.MoveToBack(groupID)
}

// QueueSwapElements swaps positions of groupIDs A and B.
func (s *Service) QueueSwapElements(ctx context.Context, user *users.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.launchHandlerQueue.Swap(groupIDA, groupIDB)
}

// QueueRemoveElement removes an element by the given groupID from the queue
func (s *Service) QueueRemoveElement(ctx context.Context, user *users.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.launchHandlerQueue.Remove(groupID)
}
