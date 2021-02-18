package nps

// This file implement the cloudsim/pkg/simulations service for this application.

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
  gormrepo "gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories/gorm"
  "github.com/jinzhu/gorm"
)

type Service interface {
	simulations.Service
	Start(ctx context.Context, request StartRequest) (*StartResponse, error)
	Stop(ctx context.Context, request StopRequest) (*StopResponse, error)

  StartSimulation(ctx context.Context, groupID simulations.GroupID) error
  StopSimulation(ctx context.Context, groupID simulations.GroupID) error

  GetStartQueue() *ign.Queue
  GetStopQueue() *ign.Queue
}

type service struct {
	repository domain.Repository
	startQueue *ign.Queue
	stopQueue  *ign.Queue
	logger     ign.Logger
}

// NewService creates a new simulation service instance.
// `repository` Database repository that hold simulation data.
// `start` The start queue.
// `stop` The stop queue.
// `
func NewService(db *gorm.DB, logger ign.Logger) Service {
  s := &service{
    // Create a new repository to hold simulation instance data.
		repository: gormrepo.NewRepository(db, logger, &Simulation{}),
    // Create the start simulation queue
		startQueue: ign.NewQueue(),
    // Create the stop simulation queue
		stopQueue:  ign.NewQueue(),
    // Store the logger
		logger:     logger,
	}

  // Create a queue to handle start requests.
	go queueHandler(s.startQueue, s.StartSimulation, s.logger)

  // Create a queue to handle stop requests.
	go queueHandler(s.stopQueue, s.StopSimulation, s.logger)

  return s
}

// queueHandler is in charge of getting the next element from the queue and passing it to the do function.
// TODO: Move somewhere else
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

func (s *service) GetStartQueue() *ign.Queue {
  return s.startQueue
}

func (s *service) GetStopQueue() *ign.Queue {
  return s.stopQueue
}

func (s *service) StartSimulation(ctx context.Context, groupID simulations.GroupID) error {

	panic("todo: implement me")
}

func (s *service) StopSimulation(ctx context.Context, groupID simulations.GroupID) error {

	panic("todo: implement me")
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

func (s *service) Start(ctx context.Context, request StartRequest) (*StartResponse, error) {
	// Business logic

	// Validate request

	// Create simulation if needed (using repository)

	// Send the simulation's group id to the queue
	gid := simulations.GroupID("test")

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
