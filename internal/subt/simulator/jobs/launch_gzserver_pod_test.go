package jobs

import (
	"github.com/stretchr/testify/mock"
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
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
	fakeSecrets "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestLaunchGazeboServerPod(t *testing.T) {
	db, err := gorm.GetDBFromEnvVars()
	defer db.Close()
	require.NoError(t, err)

	err = actions.CleanAndMigrateDB(db)
	require.NoError(t, err)

	// Set up logger
	logger := ign.NewLoggerNoRollbar("TestLaunchGazeboServerPod", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock ignition store methods for this test
	storeIgnition.On("GazeboServerLogsPath").Return("/tmp/test")
	storeIgnition.On("IP").Return("127.0.0.1")
	storeIgnition.On("Verbosity").Return("0")
	storeIgnition.On("LogsCopyEnabled").Return(true)
	storeIgnition.On("SecretsName").Return("aws-secrets")
	storeIgnition.On("Region").Return("us-west-1")
	storeIgnition.On("AccessKeyLabel").Return("aws-access-key-id")
	storeIgnition.On("SecretAccessKeyLabel").Return("aws-secret-access-key")
	storeIgnition.On("GazeboBucket").Return("gz-logs")

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")
	storeOrchestrator.On("TerminationGracePeriod").Return(time.Second)
	storeOrchestrator.On("Nameservers").Return([]string{"8.8.8.8", "8.8.4.4"})

	// Set up SPDY initializer with fake implementation
	spdyInit := spdy.NewSPDYFakeInitializer()

	secretsManager := fakeSecrets.NewFakeSecrets()
	ctx := mock.AnythingOfType("*context.emptyCtx")

	secretsManager.On("Get", ctx, "aws-secrets", "default").Return(&secrets.Secret{Data: map[string][]byte{
		"aws-access-key-id":     []byte("12345678910"),
		"aws-secret-access-key": []byte("secret"),
	}}, error(nil))

	// Set up kubernetes component
	client := kfake.NewSimpleClientset()
	po := kubernetesPods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Nodes:           nil,
		Pods:            po,
		Ingresses:       nil,
		IngressRules:    nil,
		Services:        nil,
		NetworkPolicies: nil,
	})

	// Set up platform using fake store and fake kubernetes component
	p := platform.NewPlatform(platform.Components{
		Cluster: ks,
		Store:   fakeStore,
		Secrets: secretsManager,
	})

	// Initialize generic simulation service
	simservice := simfake.NewService()

	// Create a GroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Define track name
	trackName := "Cave Circuit World 1"

	// Create a simulation for the given track
	sim := fake.NewSimulation(fake.SimulationConfig{
		GroupID: gid,
		Status:  simulations.StatusRunning,
		Kind:    simulations.SimSingle,
		Error:   nil,
		Image:   "test.org",
		Track:   trackName,
	})

	// Make the get method return the fake simulation when using
	simservice.On("Get", gid).Return(sim, error(nil))

	// Initialize tracks service
	trackService := tfake.NewService()

	// Mock Get method from tracks service
	trackService.On("Get", trackName, 0, 0).Return(&tracks.Track{
		Name:          trackName,
		Image:         "world-image.org/image",
		BridgeImage:   "bridge-image.org/image",
		StatsTopic:    "test",
		WarmupTopic:   "test",
		MaxSimSeconds: 500,
		Public:        true,
	}, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice, nil), trackService, nil)

	// Create new state: Start simulation state.
	s := state.NewStartSimulation(p, app, gid)

	// Set up action store
	store := actions.NewStore(s)

	// Run job
	_, err = LaunchGazeboServerPod.Run(store, db, &actions.Deployment{CurrentJob: "test"}, s)

	// Check if there are any errors.
	require.NoError(t, err)
}
