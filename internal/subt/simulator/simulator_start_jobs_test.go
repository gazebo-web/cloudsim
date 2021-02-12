package simulator

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
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
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

	// Initialize platform components
	c := platform.Components{
		Machines: ec2.NewMachines(ec2api, logger),
		Storage:  s3.NewStorage(storageAPI, logger),
		Cluster:  nil,
		Store:    nil,
		Secrets:  nil,
	}

	// Initialize platform
	p := platform.NewPlatform(c)

	// Initialize base application services
	simService := fake.NewService()

	// Initialize user service
	userService, err := users.NewService(ctx, nil, db, "sysadmin")
	require.NoError(t, err)

	baseapp := application.NewServices(simService, userService)

	// Initialize track repository.
	trackRepository := tracks.NewRepository(db, logger)

	// Initialize validator
	v := validator.New()

	// Initialize track services.
	trackService := tracks.NewService(trackRepository, v, logger)

	// Initialize summary service.
	summaryService := summaries.NewService(db)

	// Initialize subt application.
	app := subtapp.NewServices(baseapp, trackService, summaryService)

	// Initialize simulator
	s := NewSimulator(Config{
		DB:                  db,
		Platform:            p,
		ApplicationServices: app,
		ActionService:       nil,
	})

	// Start the simulation.
	err = s.Start(ctx, gid)
	assert.NoError(t, err)
}
