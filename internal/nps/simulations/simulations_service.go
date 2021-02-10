package simulations

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Service interface {
	simulations.Service
	Start(ctx context.Context, request StartRequest) (*StartResponse, error)
	Stop(ctx context.Context, request StopRequest) (*StopResponse, error)
}

type service struct {
	repository domain.Repository
	startQueue *ign.Queue
	stopQueue  *ign.Queue
	logger     ign.Logger
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

func NewService(repository domain.Repository, start *ign.Queue, stop *ign.Queue, logger ign.Logger) Service {
	return &service{
		repository: repository,
		startQueue: start,
		stopQueue:  stop,
		logger:     logger,
	}
}
