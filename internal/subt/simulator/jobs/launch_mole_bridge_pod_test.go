package jobs

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	tfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	ufake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users/fake"
	gormutils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestLaunchMoleBridgePods(t *testing.T) {
	db, err := gormutils.GetTestDBFromEnvVars()
	defer db.Close()
	require.NoError(t, err)

	err = actions.CleanAndMigrateDB(db)
	require.NoError(t, err)

	// Set up logger
	logger := ign.NewLoggerNoRollbar("TestLaunchMoleBridgePods", ign.VerbosityDebug)

	// Set up store
	storeOrchestrator := sfake.NewFakeOrchestrator()
	moleOrchestrator := sfake.NewFakeMole()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil, moleOrchestrator)

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")
	storeOrchestrator.On("TerminationGracePeriod").Return(time.Second)
	storeOrchestrator.On("Nameservers").Return([]string{"8.8.8.8", "8.8.4.4"})

	// Mock mole store methods for this test
	moleOrchestrator.On("BridgePulsarAddress").Return("mole-pulsar-proxy")
	moleOrchestrator.On("BridgePulsarPort").Return(6650)
	moleOrchestrator.On("BridgePulsarHTTPPort").Return(8080)
	moleOrchestrator.On("BridgeTopicRegex").Return("^subt/")

	// Set up SPDY initializer with fake implementation
	spdyInit := spdy.NewSPDYFakeInitializer()

	// Set up kubernetes component
	client := kfake.NewSimpleClientset()
	po := kubernetesPods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Pods: po,
	})

	// Set up platform using fake store and fake kubernetes component
	p, _ := platform.NewPlatform("test", platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	// Initialize generic simulation service
	simService := simfake.NewService()

	// Create a GetGroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Define track name
	trackName := "Cave Circuit World 1"

	// Create a simulation for the given track
	owner := "test"
	sim := fake.NewSimulation(fake.SimulationConfig{
		GroupID: gid,
		Status:  simulations.StatusLaunchingPods,
		Kind:    simulations.SimSingle,
		Error:   nil,
		Image:   "test.org",
		Track:   trackName,
		Owner:   &owner,
		Robots: []simulations.Robot{
			simfake.NewRobot("testA", "X1"),
			simfake.NewRobot("testB", "X2"),
			simfake.NewRobot("testC", "X3"),
		},
	})

	// Make the get method return the fake simulation
	simService.On("Get", gid).Return(sim, error(nil))

	// Initialize user service
	userService := ufake.NewFakeService()

	org := &users.Organization{
		Model: gorm.Model{ID: 123},
	}
	userService.On("GetOrganization", owner).Return(org, (*ign.ErrMsg)(nil))

	// Initialize tracks service
	trackService := tfake.NewService()

	// Mock Get method from tracks service
	moleBridgeImage := "mole-bridge-image.org/image"
	trackService.On("Get", trackName, 0, 0).Return(&tracks.Track{
		Name:            trackName,
		Image:           "world-image.org/image",
		BridgeImage:     "bridge-image.org/image",
		MoleBridgeImage: &moleBridgeImage,
		StatsTopic:      "test",
		WarmupTopic:     "test",
		MaxSimSeconds:   500,
		Public:          true,
	}, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simService, userService), trackService, nil)

	// Create new state: Start simulation state.
	s := state.NewStartSimulation(p, app, gid)
	s.GazeboServerIP = "localhost"

	// Set up action store
	store := actions.NewStore(s)

	// Run job
	_, err = LaunchMoleBridgePod.Run(store, db, &actions.Deployment{CurrentJob: "test"}, s)

	// Check if there are any errors.
	require.NoError(t, err)
}
