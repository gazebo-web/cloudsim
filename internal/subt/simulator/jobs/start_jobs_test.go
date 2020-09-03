package jobs

import (
	"context"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	simctx "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestStartJobs(t *testing.T) {
	suite.Run(t, new(startJobsTestSuite))
}

type startJobsTestSuite struct {
	suite.Suite
	actionService     actions.Servicer
	platform          platform.Platform
	appServices       application.Services
	db                *gorm.DB
	simulationService *fake.Service
	store             store.Store
	logger            ign.Logger
	ec2               cloud.Machines
	s3                cloud.Storage
	k8s               orchestrator.Cluster
	spdyInit          *spdy.Fake
}

func (s *startJobsTestSuite) SetupTest() {
	var err error
	s.actionService = actions.NewService()
	s.simulationService = fake.NewService()
	s.store = env.NewStore()
	s.logger = ign.NewLoggerNoRollbar("SimulatorJobsTest", ign.VerbosityDebug)

	s.ec2 = ec2.NewMachines(&ec2TestMock{}, s.logger)
	s.s3 = s3.NewStorage(&s3TestMock{}, s.logger)
	s.k8s = kubernetes.NewFakeKubernetes(s.logger)

	s.platform = platform.NewPlatform(s.ec2, s.s3, s.k8s, s.store)
	s.appServices = application.NewServices(s.simulationService)

	s.db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		s.FailNow(err.Error())
	}

	err = actions.MigrateDB(s.db, true)
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *startJobsTestSuite) TestCheckPendingStatusSuccessWhenStatusIsPending() {
	ctx := s.prepareActions(CheckPendingStatus)

	// Prepare simulation
	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimSingle, nil, "test")

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

func (s *startJobsTestSuite) TestCheckPendingStatusFailsWhenStatusIsNotPending() {
	ctx := s.prepareActions(CheckPendingStatus)

	// Prepare simulation
	sim := fake.NewSimulation("test-group-id", simulations.StatusRunning, simulations.SimSingle, nil, "test")

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

func (s *startJobsTestSuite) prepareActions(jobs ...*actions.Job) actions.Context {
	// Set up action with jobs
	action, err := actions.NewAction(jobs)
	s.NoError(err)

	// Register action
	err = s.actionService.RegisterAction(nil, "test", action)
	s.NoError(err)

	// Prepare context
	ctx := context.Background()
	ctx = context.WithValue(ctx, simctx.CtxPlatform, s.platform)
	ctx = context.WithValue(ctx, simctx.CtxServices, s.appServices)
	actCtx := actions.NewContext(ctx, s.logger)
	return actCtx
}

type ec2TestMock struct {
	ec2iface.EC2API
}

type s3TestMock struct {
	s3iface.S3API
}
