package simulator

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
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

	// Initializer logger
	logger := ign.NewLoggerNoRollbar("Cloudsim", ign.VerbosityDebug)

	// Initializer platform components
	c := platform.Components{
		Machines: nil,
		Storage:  nil,
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
