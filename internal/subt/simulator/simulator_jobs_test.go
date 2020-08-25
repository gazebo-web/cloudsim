package simulator

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestJobs(t *testing.T) {
	suite.Run(t, new(jobsTestSuite))
}

type jobsTestSuite struct {
	suite.Suite
	actionService     *actions.Service
	platform          platform.Platform
	appServices       application.Services
	db                *gorm.DB
	simulationService *fake.Service
}

func (s *jobsTestSuite) SetupTest() {
	var err error
	s.actionService = actions.NewService()
	s.simulationService = fake.NewService()
	s.appServices = application.NewServices(s.simulationService)
	s.db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		s.FailNow(err.Error())
	}

	s.db.DropTableIfExists(&actions.Deployment{})
	s.db.DropTableIfExists(&actions.DeploymentDataSet{})
	s.db.DropTableIfExists(&actions.DeploymentError{})

	s.db.AutoMigrate(&actions.Deployment{})
	s.db.AutoMigrate(&actions.DeploymentDataSet{})
	s.db.AutoMigrate(&actions.DeploymentError{})

}

func (s *jobsTestSuite) TestCheckPendingStatusSuccessWhenStatusIsPending() {
	ctx := s.prepareActions(JobCheckPendingStatus)

	// Prepare simulation
	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimSingle, nil)

	// Mock method
	s.simulationService.On("Get", sim.GroupID()).Return(sim, nil)

	// Execute action
	input := &actions.ExecuteInput{
		ApplicationName: nil,
		ActionName:      "test",
		Deployment:      nil,
	}
	err := s.actionService.Execute(ctx, s.db, input, sim.GroupID())
	s.NoError(err)
}

func (s *jobsTestSuite) TestCheckPendingStatusFailsWhenStatusIsNotPending() {
	ctx := s.prepareActions(JobCheckPendingStatus)

	// Prepare simulation
	sim := fake.NewSimulation("test-group-id", simulations.StatusRunning, simulations.SimSingle, nil)

	// Mock method
	s.simulationService.On("Get", sim.GroupID()).Return(sim, nil)

	// Execute action
	input := &actions.ExecuteInput{
		ApplicationName: nil,
		ActionName:      "test",
		Deployment:      nil,
	}
	err := s.actionService.Execute(ctx, s.db, input, sim.GroupID())
	s.Error(err)
	s.Equal(simulations.ErrIncorrectStatus, err)
}

func (s *jobsTestSuite) prepareActions(jobs ...*actions.Job) actions.Context {

	// Set up action with jobs
	action, err := actions.NewAction(jobs)
	s.NoError(err)

	// Register action
	err = s.actionService.RegisterAction(nil, "test", action)
	s.NoError(err)

	// Prepare context
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextPlatform, s.platform)
	ctx = context.WithValue(ctx, contextServices, s.appServices)
	actCtx := actions.NewContext(ctx)
	return actCtx
}
