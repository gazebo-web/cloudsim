package nps

// This file implement the cloudsim/pkg/simulations service for this application.

import (
	"context"
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
  "github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	// "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/env"
  ignapp "gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"

)

// StartSimulation is the state of the action that starts a simulation.
// WTF is this??
type StartSimulationData struct {
	state.PlatformGetter
	state.ServicesGetter
  platform             platform.Platform
  //services             subtapp.Services
  GroupID              simulations.GroupID
  // CreateMachinesInput  []cloud.CreateMachinesInput
  // CreateMachinesOutput []cloud.CreateMachinesOutput
}

// Services returns the underlying application services.
// \todo I really hate this pattern. 
//     1. StartSimulationData *must* have a state.ServicesGetter.
//     2. You then *must* define this function. But there is no compile-time
//        error if you don't implement this function.
//     3. If you fail to implement this function, then a segfaul will occur 
//        due to invalid memory address in a route handler.
/*func (s *StartSimulationData) Services() application.Services {
  return s.services
}*/

// Platform returns the underlying platform.
func (s *StartSimulationData) Platform() platform.Platform {
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
  applicationName string
	repository domain.Repository
	startQueue *ign.Queue
	stopQueue  *ign.Queue
	logger     ign.Logger
  db         *gorm.DB
  platform   platform.Platform
  services   ignapp.Services
  actions    actions.Servicer
}

// NewService creates a new simulation service instance.
func NewService(db *gorm.DB, logger ign.Logger) Service {
  cluster, _ := kubernetes.InitializeKubernetes(logger)

  // \todo the region string is very error prone. Can the `aws` interface provide a list of regions to choose from?
  storage, machines, _ := aws.InitializeAWS("us-east-1", logger)
  store := env.NewStore()

  // base := application.NewServices(simulationService, userService)
  // services := npsapp.NewServices(base)

	s := &service {
    applicationName: applicationName,
		// Create a new repository to hold simulation instance data.
		repository: gormrepo.NewRepository(db, logger, &Simulation{}),
		// Create the start simulation queue
		startQueue: ign.NewQueue(),
		// Create the stop simulation queue
		stopQueue: ign.NewQueue(),
		// Store the logger
		logger: logger,
    db: db,
    actions: actions.NewService(),
    // \todo: What is this, and how do I define each part?
    platform: platform.NewPlatform(platform.Components{
      // \todo How do you create a machine?
      Machines: machines,
      // \todo How do you create a storage?
      Storage: storage,
      // \todo: This is actually the orchestrator, accessed by the Orchestrator() function. Why is this named Cluster here?
      Cluster: cluster,
      // \todo How do you create a store, and what is the different from Storage above?
      Store: store,
      // \todo How do you create a secretes, and what are secrets?
      Secrets: nil,
    }),
	}

  registerActions(applicationName, s.actions)

	// Create a queue to handle start requests.
	go queueHandler(s.startQueue, s.StartSimulation, s.logger)

	// Create a queue to handle stop requests.
	go queueHandler(s.stopQueue, s.StopSimulation, s.logger)

	return s
}

// registerActions register a set of actions into the given service with the given application's name.
// It panics whenever an action could not be registered.
// \todo: This seems like a useful utility function.
func registerActions(name string, service actions.Servicer) {

	actions := map[string]actions.Jobs{
		actionNameStartSimulation: JobsStartSimulation,
		// actionNameStopSimulation:  JobsStopSimulation,
	}
	for actionName, jobs := range actions {
    fmt.Printf("Name[%s] ActionName[%s]\n", name, actionName)
    for _, j := range jobs {
      fmt.Printf("Job Name[%s]\n", j.Name)
    }
		err := registerAction(name, service, actionName, jobs)
		if err != nil {
			panic(err)
		}
	}
}

// registerAction registers the given jobs as a new action called actionName.
// The action gets registered into the given service for the given application name.
// \todo: This seems like a useful utility function.
func registerAction(applicationName string, service actions.Servicer, actionName string, jobs actions.Jobs) error {
	action, err := actions.NewAction(jobs)
	if err != nil {
		return err
	}
	err = service.RegisterAction(&applicationName, actionName, action)
	if err != nil {
		return err
	}
	return nil
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
    // service: s.services,
    // Copy the group id.
    GroupID: groupID,
  }
  store := actions.NewStore(state)

	execInput := &actions.ExecuteInput{
		ApplicationName: &s.applicationName,
		ActionName:      actionNameStartSimulation,
	}
	err := s.actions.Execute(store, s.db, execInput, state)
	if err != nil {
		return err
	}

  // \todo: What is this, why do I need it, and how do I create it?
  /* OLD
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
  */

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

func (s *service) GetWebsocketToken(groupID simulations.GroupID) (string, error) {
	panic("implement me")
}

// Start is called from the Start function in controller.go.
//
// Flow: user --> POST /start --> controller.Start() --> service.Start()
func (s *service) Start(ctx context.Context, request StartRequest) (*StartResponse, error) {
	// Business logic

	// Validate request

	// Create a simulation
  sim := Simulation{
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),

    // Name of the simulation
    Name: "nps-test-sim",

    // Create a group id
    GroupID: uuid.NewV4().String(),

    Image: request.Image,
    Args: request.Args,
  }

  if err := s.db.Create(&sim).Error; err != nil {
    return nil, err
  }

  gid := simulations.GroupID(sim.GroupID)

  // This will cause `StartSimulation` to be called because a groupId has been
  // push into the `startQueue` which is processed by the `queueHandler`.
	s.startQueue.Enqueue(gid)

	return &StartResponse{
    URI: "http://localhost:3030",
  }, nil
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
