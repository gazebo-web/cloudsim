package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

type RoutesGetter interface {
	GetRoutes() ign.Routes
}

// Server is an HTTP server that exposes an API Rest when calling ListenAndServe in the given address.
type Server interface {
	// ListenAndServe serves the HTTP server in the given address.
	// This operation should be called after initializing all the components.
	ListenAndServe(addr string) error
}

// server is a Server implementation using gorilla/mux and gorm.
type server struct {
	router     *mux.Router
	logger     ign.Logger
	db         *gorm.DB
	simulator  simulator.Simulator
	startQueue *ign.Queue
	stopQueue  *ign.Queue
}

// ListenAndServe serves the HTTP server in the given address.
// This operation should be called after initializing all the components.
func (s *server) ListenAndServe(addr string) error {
	s.logger.Debug("server: Serving HTTP on the given address:", addr)
	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		s.logger.Debug("server: Error while serving HTTP, error:", err)
		return err
	}
	return nil
}

// Config is passed to NewServer to configure the API Rest with certain components.
type Config struct {
	Router     *mux.Router
	DB         *gorm.DB
	Logger     ign.Logger
	Simulator  simulator.Simulator
	StartQueue *ign.Queue
	StopQueue  *ign.Queue
}

// NewServer initializes a new Server implementation.
func NewServer(cfg Config) Server {
	s := server{
		router:     cfg.Router,
		db:         cfg.DB,
		logger:     cfg.Logger,
		simulator:  cfg.Simulator,
		startQueue: cfg.StartQueue,
		stopQueue:  cfg.StopQueue,
	}

	go queueHandler(s.startQueue, s.simulator.Start, s.logger)

	go queueHandler(s.stopQueue, s.simulator.Stop, s.logger)

	return &s
}

// This could be moved somewhere else
func queueHandler(queue *ign.Queue, do func(ctx context.Context, gid simulations.GroupID) error, logger ign.Logger) {
	for {
		element, em := queue.DequeueOrWaitForNextElement()
		if em != nil {
			logger.Error("queue: failed to dequeue next element, error:", em.BaseError)
			break
		}

		gid, ok := element.(simulations.GroupID)
		if !ok {
			logger.Error("queue: invalid input data")
			break
		}

		ctx := context.Background()

		err := do(ctx, gid)
		if err != nil {
			logger.Error("queue: failed perform operation on the next element, error:", err)
			logger.Debug("queue: pushing element into the queue:", gid)
			queue.Enqueue(gid)
			break
		}
	}
}
