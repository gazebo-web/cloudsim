package nps

// This file implement the cloudsim/pkg/simulations service for this application.

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
  ignapp "gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"

)

// StartSimulation is the state of the action that starts a simulation.
// WTF is this??
type StartSimulationData struct {
	state.PlatformGetter
	state.ServicesGetter
  platform             platform.Platform
  GroupID              simulations.GroupID
}

// Platform returns the underlying platform.
func (s *StartSimulationData) Platform() platform.Platform {
  fmt.Printf("Getting platform\n")
  fmt.Println(s.platform)
  return s.platform
}

// Service implements the busniess logic behind the controller. A request
// comes into the controller, which then executes the appropriate function(s)
// in this service in order to handle the request.
type Service interface {
	simulations.Service
	Start(ctx context.Context, request StartRequest) (*StartResponse, error)
	Stop(ctx context.Context, request StopRequest) (*StopResponse, error)

	StartSimulation(ctx context.Context, groupID simulations.GroupID) error
	StopSimulation(ctx context.Context, groupID simulations.GroupID) error

	GetStartQueue() *ign.Queue
	GetStopQueue() *ign.Queue
}

// service stores data necessary to implement Service functions.
type service struct {
	repository domain.Repository
	startQueue *ign.Queue
	stopQueue  *ign.Queue
	logger     ign.Logger
  db         *gorm.DB
  platform   platform.Platform
  services   ignapp.Services
}

// NewService creates a new simulation service instance.
func NewService(db *gorm.DB, logger ign.Logger) Service {
	s := &service{
		// Create a new repository to hold simulation instance data.
		repository: gormrepo.NewRepository(db, logger, &Simulation{}),
		// Create the start simulation queue
		startQueue: ign.NewQueue(),
		// Create the stop simulation queue
		stopQueue: ign.NewQueue(),
		// Store the logger
		logger: logger,
    db: db,
    // \todo: What is this, and how do I define each part?
    platform: platform.NewPlatform(platform.Components{
      // \todo How do you create a machine?
      Machines: nil,
      // \todo How do you create a storage?
      Storage: nil,
      // \todo: This is actually the orchestrator, accessed by the Orchestrator() function. Why is this named Cluster here?
      Cluster: nil,
      // \todo How do you create a store, and what is the different from Storage above?
      Store: nil,
      // \todo How do you create a secretes, and what are secrets?
      Secrets: nil,
    }),
	}

	// Create a queue to handle start requests.
	go queueHandler(s.startQueue, s.StartSimulation, s.logger)

	// Create a queue to handle stop requests.
	go queueHandler(s.stopQueue, s.StopSimulation, s.logger)

	return s
}

// queueHandler is in charge of getting the next element from the queue and passing it to the do function.
func queueHandler(queue *ign.Queue, do func(ctx context.Context, gid simulations.GroupID) error, logger ign.Logger) {
	for {
		element, em := queue.DequeueOrWaitForNextElement()
		if em != nil {
			logger.Error("queue: failed to dequeue next element, error:", em.BaseError)
			continue
		}
		gid, ok := element.(simulations.GroupID)
		if !ok {
			logger.Error("queue: invalid input data")
			continue
		}
		ctx := context.Background()
		err := do(ctx, gid)
		if err != nil {
			logger.Error("queue: failed perform operation on the next element, error:", err)
			logger.Debug("queue: pushing element into the queue:", gid)
			queue.Enqueue(gid)
		}
	}
}

// GetStartQueue returns the start queue
func (s *service) GetStartQueue() *ign.Queue {
	return s.startQueue
}

// GetStopQueue returns the stop queue
func (s *service) GetStopQueue() *ign.Queue {
	return s.stopQueue
}

/////////////////////////////////////////////
var LaunchGazeboServerPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-gzserver-pod",
	PreHooks:        []actions.JobFunc{prepareGazeboCreatePodInput},
	// PostHooks:       []actions.JobFunc{},
	// RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&StartSimulationData{}),
	OutputType:      actions.GetJobDataType(&StartSimulationData{}),
})

func prepareGazeboCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\n\nPREHOOK!\n\n")

	s := store.State().(*StartSimulationData)
  fmt.Printf("State\n")
  fmt.Println(s)

	// What is this, and why is it needed???
	// namespace := s.Platform().Store().Orchestrator().Namespace()
  fmt.Printf("-------------------------\n")

	// TODO: Get ports from Ignition Store
	ports := []int32{11345, 11311}

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	volumes := []orchestrator.Volume{
		{
			Name:      "xauth",
			MountPath: "/tmp/.docker.xauth",
			HostPath:  "/tmp/.docker.xauth",
		},
		{
			Name:      "localtime",
			MountPath: "/etc/localtime",
			HostPath:  "/etc/localtime",
		},
		{
			Name:      "devinput",
			MountPath: "/dev/input",
			HostPath:  "/dev/input",
		},
		{
			Name:      "x11",
			MountPath: "/tmp/.X11-unix",
			HostPath:  "/tmp/.X11-unix",
		},
	}

	envVars := map[string]string{
		"DISPLAY":          ":0",
		"TERM":             "",
		"QT_X11_NO_MITSHM": "1",
		"XAUTHORITY":       "/tmp/.docker.xauth",
		"USE_XVFB":         "1",
	}

  // \todo: Are the regular nameservers? Are they manadatory?
	// nameservers := s.Platform().Store().Orchestrator().Nameservers()

	return jobs.LaunchPodsInput{
		{
      // Name is the name of the pod that will be created.
      // \todo: Should this be unique, and where is name used?
			Name:                          "MyTestSimulation",

      // Namespace is the namespace where the pod will live in.
      // \todo: What is a namespace?
			Namespace:                     "namespace",

      // Labels are the map of labels that will be applied to the pod.
      // \todo: What are the labels used for?
      Labels:                        map[string]string{"key":"value"},

      // RestartPolicy defines how the pod should react after an error.
      // \todo: What are the restart policies, and how do I choose one?
			RestartPolicy:                 orchestrator.RestartPolicyNever,

      // TerminationGracePeriodSeconds is the time duration in seconds the pod needs to terminate gracefully.
      // \todo: What does this do?
			TerminationGracePeriodSeconds: 0,

      // NodeSelector defines the node where the pod should run in.
      // \todo: What does this mean, and how do I know what value to put in???
			NodeSelector:                  orchestrator.NewSelector(map[string]string{
    "cloudsim_groupid": s.GroupID.String() }),

      // Containers is the list of containers that should be created inside the pod.
      // \todo: What is a container? 
			Containers: []orchestrator.Container{
        {
          // Name is the container's name.
					Name:                     "nps-novnc",
          // Image is the image running inside the container.
					Image:                    "osrf/ros:melodic-desktop-full",
          // Args passed to the Command. Cannot be updated.
					Args:                     []string{"gazebo"},
          // Privileged defines if the container should run in privileged mode.
					Privileged:               &privileged,
          // AllowPrivilegeEscalation is used to define if the container is allowed to scale its privileges.
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
          // Ports is the list of ports that should be opened.
					Ports:                    ports,
          // Volumes is the list of volumes that should be mounted in the container.
					Volumes:                  volumes,
          // EnvVars is the list of env vars that should be passed into the container.
					EnvVars:                  envVars,
				},
			},
			Volumes:     volumes,

      // \todo: Is this required?
			// Nameservers: nameservers,
		},
	}, nil
}

/////////////////////////////////////////////


// StartSimulation is called from service.Start(), and it should actually start
// the simulation running.
//
// Flow: user --> POST /start --> controller.Start() --> service.Start() --> service.StartSimulation
func (s *service) StartSimulation(ctx context.Context, groupID simulations.GroupID) error {

  // You must create a data structure to hold data that is then "stored" in a
  // NewStore on the following line. This store and the data contained in the 
  // store is passed into the jobs, which perform the work of launching 
  // K8 nodes (cloud machines) and K8 pods (docker containers). 
  state := &StartSimulationData{
    // Copy the platform information. 
    platform: s.platform,
    // Copy the group id.
    GroupID: groupID,
  }
  store := actions.NewStore(state)

  // \todo: What is this, why do I need it, and how do I create it?
  action := &actions.Deployment{}

  // \todo: What is this, why do I need it, and how do I create it?
  launchPodsInput := jobs.LaunchPodsInput{}

  // Run the job. This will launch the docker container, hooray!!
  _, err := LaunchGazeboServerPod.Run(store, s.db, action, launchPodsInput)

  // Check for errors, always a good thing to do.
  if err != nil {
    fmt.Printf("\n\nError launching pod\n\n")
    fmt.Println(err)
  }

	fmt.Printf("StartSimulation for groupID[%s]\n", groupID)
	return nil
}

func (s *service) StopSimulation(ctx context.Context, groupID simulations.GroupID) error {

	panic("todo: StopSimulation")
}

func (s *service) Get(groupID simulations.GroupID) (simulations.Simulation, error) {
	panic("implement me")
}

func (s *service) Reject(groupID simulations.GroupID) (simulations.Simulation, error) {
	panic("implement me")
}

func (s *service) GetParent(groupID simulations.GroupID) (simulations.Simulation, error) {
	panic("implement me")
}

func (s *service) UpdateStatus(groupID simulations.GroupID, status simulations.Status) error {
	panic("implement me")
}

func (s *service) Update(groupID simulations.GroupID, simulation simulations.Simulation) error {
	panic("implement me")
}

func (s *service) GetRobots(groupID simulations.GroupID) ([]simulations.Robot, error) {
	panic("implement me")
}

// Start is called from the Start function in controller.go.
//
// Flow: user --> POST /start --> controller.Start() --> service.Start()
func (s *service) Start(ctx context.Context, request StartRequest) (*StartResponse, error) {
	// Business logic

	// Validate request

	// Create simulation if needed (using repository)

	// Send the simulation's group id to the queue
	gid := simulations.GroupID("test")

  // This will cause `StartSimulation` to be called because a groupId has been
  // push into the `startQueue` which is processed by the `queueHandler`.
	s.startQueue.Enqueue(gid)

	return &StartResponse{}, nil
}

func (s *service) Stop(ctx context.Context, request StopRequest) (*StopResponse, error) {
	// Business logic

	// Validate request

	// Mark simulation as stopped

	// Send the group id to the queue
	gid := simulations.GroupID("test")

	s.stopQueue.Enqueue(gid)

	return &StopResponse{}, nil
}
