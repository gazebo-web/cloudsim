package nps

// This file implement the cloudsim/pkg/simulations service for this application.

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	ignapp "gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"time"
)

// Service implements the busniess logic behind the controller. A request
// comes into the controller, which then executes the appropriate function(s)
// in this service in order to handle the request.
type Service interface {
	simulations.Service

	// Start will run the StartSimulationAction to launch cloud machines
	// and a docker image.
	Start(tx *gorm.DB, ctx context.Context, request StartRequest) (*StartResponse, error)

	// Stop will run the StopSimulationAction to terminate clouds machines.
	Stop(ctx context.Context, request StopRequest) (*StopResponse, error)

	// StartQueueHandler processes entries in the startQueue.
	StartQueueHandler(ctx context.Context, groupID simulations.GroupID) error

	// StopQueueHandler processes entries in the startQueue.
	StopQueueHandler(ctx context.Context, groupID simulations.GroupID) error

	// GetStartQueue returns the start queue
	GetStartQueue() *ign.Queue

	// GetStopQueue returns the stop queue
	GetStopQueue() *ign.Queue
}

// service stores data necessary to implement Service functions.
type service struct {
	// applicationName is the name of the application.
	applicationName string
	repository      domain.Repository
	startQueue      *ign.Queue
	stopQueue       *ign.Queue
	logger          ign.Logger
	db              *gorm.DB
	platform        platform.Platform
	services        ignapp.Services
	actions         actions.Servicer
}

// NewService creates a new simulation service instance.
func NewService(db *gorm.DB, logger ign.Logger) Service {
	// This gets a reference to the Kubernetes cluster that the services can use.
	cluster, _ := kubernetes.InitializeKubernetes(logger)

	// \todo the region string is very error prone. Can the `aws` interface
	// provide a list of regions to choose from?
	storage, machines, _ := aws.InitializeAWS("us-east-1", logger)

	// \todo Why do I need to make a Store here?
	store := env.NewStore()

	s := &service{
		applicationName: applicationName,
		// Create a new repository to hold simulation instance data.
		repository: gormrepo.NewRepository(db, logger, &Simulation{}),
		// Create the start simulation queue. The start queue is used to process
		// simulation start requests.
		startQueue: ign.NewQueue(),
		// Create the stop simulation queue. The stop queue is used to process
		// simulation stop requests.
		stopQueue: ign.NewQueue(),
		// Store the logger
		logger: logger,
		// Store the database reference.
		db:      db,
		actions: actions.NewService(),
		// \todo: What is a "Platform"?
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
	go queueHandler(s.startQueue, s.StartQueueHandler, s.logger)

	// Create a queue to handle stop requests.
	go queueHandler(s.stopQueue, s.StopQueueHandler, s.logger)

	return s
}

// GetStartQueue returns the start queue
func (s *service) GetStartQueue() *ign.Queue {
	return s.startQueue
}

// GetStopQueue returns the stop queue
func (s *service) GetStopQueue() *ign.Queue {
	return s.stopQueue
}

// StartQueueHandler is called from service.Start(), and it should actually
// start the simulation running.
//
// Origin: user --> POST /start --> controller.Start() --> service.Start() --> service.StartQueueHandler
func (s *service) StartQueueHandler(ctx context.Context, groupID simulations.GroupID) error {

	// You must create a data structure to hold data that is then "stored" in a
	// NewStore on the following line. This store and the data contained in the
	// store is passed into the jobs, which perform the work of launching
	// K8 nodes (cloud machines) and K8 pods (docker containers).
	state := &StartSimulationData{
		// Copy the platform information.
		platform: s.platform,
		// Copy the group id.
		GroupID: groupID,
		logger:  s.logger,
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

	s.logger.Info("Starting simulation for groupID[%s]\n", groupID)
	return nil
}

func (s *service) StopQueueHandler(ctx context.Context, groupID simulations.GroupID) error {

	panic("todo: StopQueueHandler")
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

func (s *service) MarkStopped(groupID simulations.GroupID) (error) {
	panic("implement me")
}

// Start is called from the Start function in controller.go.
//
// Origin: user --> POST /start --> controller.Start() --> service.Start()
// Next: StartQueueHandler
func (s *service) Start(tx *gorm.DB, ctx context.Context, request StartRequest) (*StartResponse, error) {
	// Add business logic here to validate a request, update a database table,
	// etc.

	// This creates the database entry to keep track of simulation instances.
	sim := Simulation{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		// Name of the simulation
		Name: "nps-test-sim",

		// Create a group id
		GroupID: uuid.NewV4().String(),
		Status:  "starting",

		Image: request.Image,
		Args:  request.Args,
	}

	if err := tx.Create(&sim).Error; err != nil {
		return nil, err
	}

	gid := simulations.GroupID(sim.GroupID)

	// This will cause `StartQueueHandler` to be called because a groupId has been
	// push into the `startQueue` which is processed by the `queueHandler`.
	s.startQueue.Enqueue(gid)

	return &StartResponse{
		Message: "Simulation instance is starting. Use the URI to get status updates",
		Simulation: GetSimulationResponse{
			Name:    sim.Name,
			GroupID: sim.GroupID,
			Status:  sim.Status,
			Image:   sim.Image,
			Args:    sim.Args,
			URI:     "simulations/" + sim.GroupID,
		},
	}, nil
}

// Stop is called from the Stop function in controller.go.
//
// Origin: user --> POST /stop --> controller.Stop() --> service.Stop()
// Next: StopQueueHandler
func (s *service) Stop(ctx context.Context, request StopRequest) (*StopResponse, error) {
	// Add business logic here to validate a request, update a database table,
	// etc.
	// Send the group id to the queue
	gid := simulations.GroupID("test")

	s.stopQueue.Enqueue(gid)

	return &StopResponse{}, nil
}

// ListSimulations is called from the ListSimulations function in
// controller.go.
//
// Origin: user --> GET /list --> controller.ListSimulations() --> service.ListSimulations()
/*func (s *service) ListSimulations(ctx context.Context, request ListRequest) (*ListResponse, error) {

  var simulations Simulations
  s.db.Find(&simulations)

  var response ListResponse
  for _, sim := range simulations {
    response.Simulations = append(response.Simulations, ListResponseSimulation{
      Name: sim.Name,
      GroupID: sim.GroupID,
      Status: sim.Status,
      Image: sim.Image,
      Args: sim.Args,
      URI: sim.URI,
    })
  }

	// Send the group id to the queue
	return &response, nil
}*/

///////////////////////////////////////
// It would be nice to make the following function general purpose functions
// that live in the main Cloudsim codebase

// registerActions registers a set of actions into the given service with the given application's name.
// It panics whenever an action could not be registered.
// \todo: This seems like a useful utility function that should exist in the
// main cloudsim code. This was copied from the subt application.
func registerActions(name string, service actions.Servicer) {

	actions := map[string]actions.Jobs{
		actionNameStartSimulation: StartSimulationAction,
		// actionNameStopSimulation:  StopSimulationAction,
	}
	for actionName, jobs := range actions {
		err := registerAction(name, service, actionName, jobs)
		if err != nil {
			panic(err)
		}
	}
}

// registerAction registers the given jobs as a new action called actionName.
// The action gets registered into the given service for the given application name.
// \todo: This seems like a useful utility function that should exist in the
// main cloudsim code. This was copied from the subt application.
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
// \todo: This seems like a useful utility function that should exist in the
// main cloudsim code. This was copied from the subt application.
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

///////////////////////////////////////
