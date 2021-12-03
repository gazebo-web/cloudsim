package e2e

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	cloud "gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	email "gitlab.com/ignitionrobotics/web/cloudsim/pkg/email/implementations/ses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/migrations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/runsim"
	secrets "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/store"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	legacy "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	fuelusers "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"net/http/httptest"
	"testing"
)

func TestStopSimulationAction(t *testing.T) {
	// Set up context
	ctx := context.Background()

	// Connect to the database
	db, err := gorm.GetTestDBFromEnvVars()
	require.NoError(t, err)

	migrations.DBDropModels(context.Background(), db)
	migrations.DBMigrate(context.Background(), db)

	// Clean and migrate database
	err = gorm.CleanAndMigrateModels(
		db,
		&fuelusers.User{},
		&fuelusers.Organization{},
	)
	require.NoError(t, err)

	// Migrate actions
	err = actions.CleanAndMigrateDB(db)
	require.NoError(t, err)

	// Initialize logger
	logger := ign.NewLoggerNoRollbar("Cloudsim", ign.VerbosityDebug)

	// Initialize base application services
	simService := legacy.NewSubTSimulationServiceAdaptor(db)

	robot := legacy.SubTRobot{
		Name:    "X1",
		Type:    "X1",
		Image:   "test.org/image",
		Credits: 270,
	}

	extra := legacy.ExtraInfoSubT{
		Circuit: "Urban Circuit 1",
		Robots:  []legacy.SubTRobot{robot},
	}

	extraInfo, err := extra.ToJSON()
	require.NoError(t, err)

	sim, err := simService.Create(simulations.CreateSimulationInput{
		Name:      "sim-test",
		Owner:     nil,
		Creator:   "sysadmin",
		Image:     []string{"test.org/image"},
		Private:   false,
		StopOnEnd: false,
		Extra:     *extraInfo,
		Track:     "Urban Circuit 1",
		Robots:    "X1",
	})
	require.NoError(t, err)

	err = simService.UpdateStatus(sim.GetGroupID(), simulations.StatusTerminateRequested)
	require.NoError(t, err)

	// Initialize mock for EC2
	ec2api := mock.NewEC2(
		mock.NewEC2Instance("test-fc-1", subtapp.GetTagsInstanceSpecific("field-computer", sim.GetGroupID(), "sim", "cloudsim", "field-computer")),
		mock.NewEC2Instance("test-gz-1", subtapp.GetTagsInstanceSpecific("gzserver", sim.GetGroupID(), "sim", "cloudsim", "gzserver")),
	)
	ec2Machines, err := ec2.NewMachines(&ec2.NewInput{
		API:            ec2api,
		CostCalculator: cloud.NewCostCalculatorEC2(nil),
		Logger:         logger,
		Zones: []ec2.Zone{
			{
				Zone:     "test",
				SubnetID: "test",
			},
		},
	})
	require.NoError(t, err)

	// Initialize mock for S3
	storageBackend := s3mem.New()
	storageFake := gofakes3.New(storageBackend)
	storageServer := httptest.NewServer(storageFake.Server())

	storageSessionConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String(storageServer.URL),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	storageSession, err := session.NewSession(storageSessionConfig)
	require.NoError(t, err)

	storageAPI := s3.NewAPI(storageSession)

	// Initialize mock for Kubernetes orchestrator
	kubernetesClientset := kfake.NewSimpleClientset()

	fakeSPDY := spdy.NewSPDYFakeInitializer()

	cluster := kubernetes.NewDefaultKubernetes(kubernetesClientset, fakeSPDY, logger)

	// Initialize env vars
	configStore, err := store.NewStoreFromEnvVars()
	require.NoError(t, err)

	// Initialize secrets
	secrets := secrets.NewKubernetesSecrets(kubernetesClientset.CoreV1())

	runsimManager := runsim.NewManager()

	// Initialize email API using mocks.
	emailAPI := mock.NewSES()

	// Initialize platform components
	c := platform.Components{
		Machines:           ec2Machines,
		Storage:            s3.NewStorage(storageAPI, logger),
		Cluster:            cluster,
		Store:              configStore,
		Secrets:            secrets,
		RunningSimulations: runsimManager,
		EmailSender:        email.NewEmailSender(emailAPI, logger),
	}

	// Initialize platform
	p, _ := platform.NewPlatform("test", c)

	// Simulation should be set to terminate requested.
	err = simService.UpdateStatus(sim.GetGroupID(), simulations.StatusTerminateRequested)
	require.NoError(t, err)

	// Initializing permissions
	perm := permissions.Permissions{}
	err = perm.Init(db, "sysadmin")
	require.NoError(t, err)

	// Initialize user service
	userService, err := users.NewService(ctx, &perm, db, "sysadmin")
	require.NoError(t, err)

	baseApp := application.NewServices(simService, userService, nil)

	// Initialize track repository.
	trackRepository := tracks.NewRepository(db, logger)

	// Initialize validator
	v := validator.New()

	// Initialize track services.
	trackService := tracks.NewService(trackRepository, v, logger)

	maxSimSeconds := 720

	// Create testing track
	_, err = trackService.Create(tracks.CreateTrackInput{
		Name:          "Urban Circuit 1",
		Image:         "test.org/image",
		BridgeImage:   "test.org/bridge-image",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: maxSimSeconds,
		Public:        true,
		World:         "cloudsim_sim.ign;worldName:=simple_urban_01;circuit:=urban",
	})
	require.NoError(t, err)

	// Initialize summary service.
	summaryService := summaries.NewService(db)

	// Initialize subt application.
	app := subtapp.NewServices(baseApp, trackService, summaryService)

	rs := runsim.NewRunningSimulation(sim)
	ws := ignws.NewPubSubTransporterMock()

	ws.On("Disconnect").Return(error(nil))
	ws.On("IsConnected").Return(false)

	err = runsimManager.Add(sim.GetGroupID(), rs, ws)
	require.NoError(t, err)

	username := "sysadmin"
	email := "test@openrobotics.org"
	name := "Tester"
	require.NoError(t, db.Model(&fuelusers.User{}).Create(&fuelusers.User{
		Name:     &name,
		Username: &username,
		Email:    &email,
	}).Error)

	*kubernetesClientset = *kfake.NewSimpleClientset(
		&apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetServiceNameWebsocket(sim.GetGroupID()),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameGazeboServer(sim.GetGroupID()),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridge(sim.GetGroupID(), subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameFieldComputer(sim.GetGroupID(), subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameMappingServer(sim.GetGroupID()),
				Namespace: "default",
			},
		},
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameGazeboServer(sim.GetGroupID()),
				Namespace: "default",
			},
		},
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridge(sim.GetGroupID(), subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
		&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameFieldComputer(sim.GetGroupID(), subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
	)

	stopAction, err := actions.NewAction(actions.Jobs{
		jobs.CheckSimulationTerminateRequestedStatus,
		jobs.CheckStopSimulationIsNotParent,
		jobs.SetStoppedAt,
		jobs.DisconnectWebsocket,
		jobs.RemoveRunningSimulation,
		jobs.SetSimulationStatusToProcessingResults,
		jobs.UploadLogs,
		// TODO: These jobs have been postponed from being tested because we're not able to read files from pods right now.
		// jobs.ReadScore,
		// jobs.ReadStats,
		// jobs.ReadRunData,

		// TODO: Since no score has been read, this job fails.
		// jobs.SaveScore,
		jobs.SaveSummary,
		jobs.SendSummaryEmail,
		jobs.SetSimulationStatusToDeletingPods,

		// TODO: We haven't configure the gloo fake client, and therefore we can't test removing rules from the Gloo ingress.
		// jobs.RemoveIngressRulesGloo,
		jobs.RemoveWebsocketService,
		jobs.RemoveNetworkPolicies,
		jobs.RemovePods,
		jobs.SetSimulationStatusToDeletingNodes,
		jobs.RemoveInstances,
		jobs.SetSimulationStatusToTerminated,
	})
	require.NoError(t, err)

	actionService := actions.NewService(logger)

	// Initialize simulator
	s := simulator.NewSimulator(simulator.Config{
		DB:                    db,
		ApplicationServices:   app,
		ActionService:         actionService,
		DisableDefaultActions: true,
	})

	appName := simulator.ApplicationName
	actionService.RegisterAction(&appName, simulator.ActionNameStopSimulation, stopAction)

	// Stop the simulation.
	err = s.Stop(ctx, p, sim.GetGroupID())
	assert.NoError(t, err)
}
