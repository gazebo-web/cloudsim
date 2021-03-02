package e2e

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	legacy "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	kfake "k8s.io/client-go/kubernetes/fake"
	"net/http/httptest"
	"testing"
)

func TestStartSimulationAction(t *testing.T) {
	// Set up context
	ctx := context.Background()

	// Define simulation GroupID.
	gid := simulations.GroupID(uuid.NewV4().String())

	// Connect to the database
	db, err := gorm.GetTestDBFromEnvVars()
	require.NoError(t, err)

	// Clean and migrate database
	err = gorm.CleanAndMigrateModels(db, &legacy.SimulationDeployment{})
	require.NoError(t, err)

	// Migrate actions
	err = actions.CleanAndMigrateDB(db)
	require.NoError(t, err)

	// Initialize logger
	logger := ign.NewLoggerNoRollbar("Cloudsim", ign.VerbosityDebug)

	// Initialize mock for EC2
	ec2api := mock.NewEC2()

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
	configStore := env.NewStore()

	// Initialize secrets
	secrets := secrets.NewKubernetesSecrets(kubernetesClientset.CoreV1())

	// Initialize platform components
	c := platform.Components{
		Machines: ec2.NewMachines(ec2api, logger),
		Storage:  s3.NewStorage(storageAPI, logger),
		Cluster:  cluster,
		Store:    configStore,
		Secrets:  secrets,
	}

	// Initialize platform
	p := platform.NewPlatform(c)

	// Initialize base application services
	simService := legacy.NewSubTSimulationServiceAdaptor(db)

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

	// Initialize summary service.
	summaryService := summaries.NewService(db)

	// Initialize subt application.
	app := subtapp.NewServices(baseApp, trackService, summaryService)

	// Initialize simulator
	s := simulator.NewSimulator(simulator.Config{
		DB:                  db,
		Platform:            p,
		ApplicationServices: app,
		ActionService:       actions.NewService(),
	})

	// Start the simulation.
	err = s.Start(ctx, gid)
	assert.NoError(t, err)
}
