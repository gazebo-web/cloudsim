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
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/migrations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	secrets "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	legacy "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"net/http/httptest"
	"testing"
)

func TestStartSimulationAction(t *testing.T) {
	// Set up context
	ctx := context.Background()

	// Connect to the database
	db, err := gorm.GetTestDBFromEnvVars()
	require.NoError(t, err)

	// Clean and migrate database
	migrations.DBDropModels(context.Background(), db)
	migrations.DBMigrate(context.Background(), db)

	// Migrate actions
	err = actions.CleanAndMigrateDB(db)
	require.NoError(t, err)

	// Initialize logger
	logger := ign.NewLoggerNoRollbar("Cloudsim", ign.VerbosityDebug)

	// Initialize mock for EC2
	ec2api := mock.NewEC2()
	ec2Machines, err := ec2.NewMachines(&ec2.NewInput{
		API: ec2api,
		Logger: logger,
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
	secretsManager := secrets.NewKubernetesSecrets(kubernetesClientset.CoreV1())

	// Initialize platform components
	c := platform.Components{
		Machines: ec2Machines,
		Storage:  s3.NewStorage(storageAPI, logger),
		Cluster:  cluster,
		Store:    configStore,
		Secrets:  secretsManager,
	}

	// Initialize platform
	p, _ := platform.NewPlatform("test", c)

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

	// Initializing permissions
	perm := permissions.Permissions{}
	err = perm.Init(db, "sysadmin")
	require.NoError(t, err)

	// Initialize user service
	userService, err := users.NewService(ctx, &perm, db, "sysadmin")
	require.NoError(t, err)

	baseApp := application.NewServices(simService, userService)

	// Initialize track repository.
	trackRepository := tracks.NewRepository(db, logger)

	// Initialize validator
	v := validator.New()

	// Initialize track services.
	trackService := tracks.NewService(trackRepository, v, logger)

	// Create testing track
	_, err = trackService.Create(tracks.CreateTrackInput{
		Name:          "Urban Circuit 1",
		Image:         "test.org/image",
		BridgeImage:   "test.org/bridge-image",
		StatsTopic:    "/stats",
		WarmupTopic:   "/warmup",
		MaxSimSeconds: 720,
		Public:        true,
		World:         "cloudsim_sim.ign;worldName:=simple_urban_01;circuit:=urban",
	})
	require.NoError(t, err)

	// Initialize summary service.
	summaryService := summaries.NewService(db)

	// Initialize subt application.
	app := subtapp.NewServices(baseApp, trackService, summaryService)

	t.Run("First phase", func(t *testing.T) {
		actionService := actions.NewService(logger)

		kClient := kfake.NewSimpleClientset(
			&apiv1.Node{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:   subtapp.GetPodNameGazeboServer(sim.GetGroupID()),
					Labels: subtapp.GetNodeLabelsGazeboServer(sim.GetGroupID()).Map(),
				},
				Spec: apiv1.NodeSpec{},
				Status: apiv1.NodeStatus{
					Conditions: []apiv1.NodeCondition{
						{
							Type:   "Ready",
							Status: "True",
						},
					},
				},
			},
			&apiv1.Node{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:   subtapp.GetPodNameFieldComputer(sim.GetGroupID(), subtapp.GetRobotID(0)),
					Labels: subtapp.GetNodeLabelsFieldComputer(sim.GetGroupID(), &robot).Map(),
				},
				Spec: apiv1.NodeSpec{},
				Status: apiv1.NodeStatus{
					Conditions: []apiv1.NodeCondition{
						{
							Type:   "Ready",
							Status: "True",
						},
					},
				},
			},
		)

		CopyDeepKubernetesClientset(kubernetesClientset, kClient)

		s := simulator.NewSimulator(simulator.Config{
			DB:                    db,
			ApplicationServices:   app,
			ActionService:         actionService,
			DisableDefaultActions: true,
		})

		startActions, err := actions.NewAction(actions.Jobs{
			jobs.CheckSimulationPendingStatus,
			jobs.CheckStartSimulationIsNotParent,
			jobs.CheckSimulationNoErrors,
			jobs.SetSimulationStatusToLaunchInstances,
			jobs.LaunchInstances,
			jobs.SetSimulationStatusToWaitInstances,
			jobs.WaitForInstances,
			jobs.SetSimulationStatusToWaitNodes,
			jobs.WaitForNodes,
			jobs.SetSimulationStatusToLaunchPods,
			jobs.CreateNetworkPolicyGazeboServer,
			jobs.LaunchGazeboServerPod,
		})
		require.NoError(t, err)

		appName := simulator.ApplicationName
		actionService.RegisterAction(&appName, simulator.ActionNameStartSimulation, startActions)

		// Start the simulation.
		err = s.Start(ctx, p, sim.GetGroupID())
		assert.NoError(t, err)
	})

	t.Run("Second phase", func(t *testing.T) {
		actionService := actions.NewService(logger)

		s := simulator.NewSimulator(simulator.Config{
			DB:                    db,
			ApplicationServices:   app,
			ActionService:         actionService,
			DisableDefaultActions: true,
		})

		kClient := kfake.NewSimpleClientset(
			&apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      subtapp.GetPodNameGazeboServer(sim.GetGroupID()),
					Namespace: p.Store().Orchestrator().Namespace(),
					Labels:    subtapp.GetPodLabelsGazeboServer(sim.GetGroupID(), nil).Map(),
				},
				Spec: apiv1.PodSpec{},
				Status: apiv1.PodStatus{
					Conditions: []apiv1.PodCondition{
						{
							Type:   apiv1.PodInitialized,
							Status: apiv1.ConditionTrue,
						},
						{
							Type:   apiv1.PodScheduled,
							Status: apiv1.ConditionTrue,
						},
						{
							Type:   apiv1.PodReady,
							Status: apiv1.ConditionTrue,
						},
					},
					PodIP: "1.1.1.1",
				},
			},
		)

		CopyDeepKubernetesClientset(kubernetesClientset, kClient)

		startActions, err := actions.NewAction(actions.Jobs{
			jobs.WaitForGazeboServerPod,
			jobs.CreateNetworkPolicyCommsBridges,
			jobs.CreateNetworkPolicyFieldComputers,
			jobs.LaunchCommsBridgePods,
			jobs.LaunchCommsBridgeCopyPods,
		})
		require.NoError(t, err)

		appName := simulator.ApplicationName
		actionService.RegisterAction(&appName, simulator.ActionNameStartSimulation, startActions)

		// Start the simulation.
		err = s.Start(ctx, p, sim.GetGroupID())
		assert.NoError(t, err)
	})

	t.Run("Third phase", func(t *testing.T) {
		actionService := actions.NewService(logger)

		s := simulator.NewSimulator(simulator.Config{
			DB:                    db,
			ApplicationServices:   app,
			ActionService:         actionService,
			DisableDefaultActions: true,
		})

		kClient := kfake.NewSimpleClientset(&apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridge(sim.GetGroupID(), subtapp.GetRobotID(0)),
				Namespace: p.Store().Orchestrator().Namespace(),
				Labels:    subtapp.GetPodLabelsCommsBridge(sim.GetGroupID(), nil, &robot).Map(),
			},
			Spec: apiv1.PodSpec{},
			Status: apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type:   apiv1.PodInitialized,
						Status: apiv1.ConditionTrue,
					},
					{
						Type:   apiv1.PodScheduled,
						Status: apiv1.ConditionTrue,
					},
					{
						Type:   apiv1.PodReady,
						Status: apiv1.ConditionTrue,
					},
				},
				PodIP: "1.1.1.1",
			},
		})

		CopyDeepKubernetesClientset(kubernetesClientset, kClient)

		startActions, err := actions.NewAction(actions.Jobs{
			jobs.WaitForCommsBridgePodIPs,
			jobs.GetCommsBridgePodIP,
			jobs.LaunchFieldComputerPods,
			jobs.SetSimulationStatusToWaitPods,
			// NOTE: This job have been commented out because we don't have a mechanism to return
			// a list of pods based in a group of labels, as the WaitSimulationPods job does.
			// jobs.WaitSimulationPods,

			// NOTE: This jobs should be commented out once we have a mechanism to better mock the Websocket connection.
			// jobs.SetWebsocketConnection,
			// jobs.AddRunningSimulation,
			jobs.SetSimulationStatusToRunning,
		})
		require.NoError(t, err)

		appName := simulator.ApplicationName
		actionService.RegisterAction(&appName, simulator.ActionNameStartSimulation, startActions)

		// Start the simulation.
		err = s.Start(ctx, p, sim.GetGroupID())
		assert.NoError(t, err)
	})
}

// CopyDeepKubernetesClientset performs a deep copy of the content from the second argument into the first argument.
func CopyDeepKubernetesClientset(to *kfake.Clientset, from *kfake.Clientset) {
	to.Resources = from.Resources
	to.ProxyReactionChain = from.ProxyReactionChain
	to.ReactionChain = from.ReactionChain
	to.WatchReactionChain = from.WatchReactionChain
}
