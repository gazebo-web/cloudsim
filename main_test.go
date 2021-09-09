package main

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/migrations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"log"
	"os"
	"testing"
	"time"
)

// This function applies to ALL tests in the application.
// It will run the test and then clean the database.
func TestMain(m *testing.M) {
	code := m.Run()
	packageTearDown(nil)
	log.Println("Cleaned database tables after all tests")
	os.Exit(code)
}

// setup is the Test Setup function. It should be invoked by all tests on their
// first line.
// You can use the alternative function `setupWithCustomInitalizer`
func setup() {
	setupWithCustomInitializer(nil)
}

type customInitializer func(ctx context.Context)

// setup helper function
func setupWithCustomInitializer(customFn customInitializer) {
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityDebug)
	logCtx := ign.NewContextWithLogger(context.Background(), logger)

	worldStatsTopic := "/world/default/stat"
	worldWarmupTopic := "/subt/start"

	// Make sure we don't have data from other tests.
	// For this we drop db tables and recreate them.
	packageTearDown(logCtx)
	createDBTablesAndData(logCtx, worldStatsTopic, worldWarmupTopic)

	// Mocking
	//  Mock the Sleep function so that it returns instantly
	sim.Sleep = func(d time.Duration) {
		return
	}

	if customFn != nil {
		customFn(logCtx)
	}

	// Check for auth0 environment variables.
	if os.Getenv("IGN_TEST_JWT") == "" {
		log.Printf("Missing IGN_TEST_JWT env variable." +
			"Authentication will not work.")
	}

	// Create the router, and indicate that we are testing
	igntest.SetupTest(globals.Server.Router)

	if globals.TransportTestMock == nil {
		globals.TransportTestMock = ignws.NewPubSubTransporterMock()
	}

	expectedFn := mock.AnythingOfType("ign.Callback")
	globals.TransportTestMock.On("Subscribe", worldStatsTopic, expectedFn).Return(nil)
	globals.TransportTestMock.On("Subscribe", worldWarmupTopic, expectedFn).Return(nil)
	globals.TransportTestMock.On("Disconnect").Return(error(nil))
	globals.TransportTestMock.On("IsConnected").Return(true)
}

// Clean up our mess
func packageTearDown(ctx context.Context) {
	if ctx == nil {
		ctx = ign.NewContextWithLogger(context.Background(), ign.NewLoggerNoRollbar("test", ign.VerbosityDebug))
	}

	cleanDBTables(ctx)
}

func cleanDBTables(ctx context.Context) {
	migrations.DBDropModels(ctx, globals.Server.Db)
}

func createDBTablesAndData(ctx context.Context, worldStatsTopic, worldWarmupTopic string) {
	migrations.DBMigrate(ctx, globals.Server.Db)
	// After removing tables we can ask casbin to re initialize
	if err := globals.Permissions.Reload(sysAdminForTest); err != nil {
		log.Fatal("Error reloading casbin policies", err)
	}
	// Apply custom indexes. Eg: fulltext indexes
	migrations.DBAddCustomIndexes(ctx, globals.Server.Db)

	// Insert SubT Circuit Rules
	circuits := []*sim.SubTCircuitRules{
		{
			Circuit:            sptr("Virtual Stix"),
			Worlds:             sptr("worldName:=testworld"),
			Times:              sptr("1"),
			Image:              sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:        sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:   sptr(worldStatsTopic),
			WorldWarmupTopics:  sptr(worldWarmupTopic),
			WorldMaxSimSeconds: sptr("0"),
			MaxCredits:         intptr(1000),
		},
		{
			Circuit:               sptr("Tunnel Circuit"),
			Worlds:                sptr("worldName:=testworld"),
			Times:                 sptr("1"),
			Image:                 sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:           sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:      sptr(worldStatsTopic),
			WorldWarmupTopics:     sptr(worldWarmupTopic),
			WorldMaxSimSeconds:    sptr("0"),
			Seeds:                 sptr("10"),
			MaxCredits:            nil,
			CompetitionDate:       nil,
			RequiresQualification: boolptr(false),
		},
		{
			Circuit:            sptr("Tunnel Practice 1"),
			Worlds:             sptr("worldName:=testworld"),
			Times:              sptr("1"),
			Image:              sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:        sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:   sptr(worldStatsTopic),
			WorldWarmupTopics:  sptr(worldWarmupTopic),
			WorldMaxSimSeconds: sptr("0"),
			MaxCredits:         intptr(10),
		},
		{
			Circuit:            sptr("Urban Practice 1"),
			Worlds:             sptr("worldName:=urban_circuit,worldName:=urban_circuit"),
			Times:              sptr("3,3"),
			Image:              sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:        sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:   sptr(worldStatsTopic),
			WorldWarmupTopics:  sptr(worldWarmupTopic),
			WorldMaxSimSeconds: sptr("0"),
			MaxCredits:         nil,
			CompetitionDate:    timeptr(time.Now().Add(time.Hour * 24)),
		},
		{
			Circuit:               sptr("Urban Practice 2"),
			Worlds:                sptr("worldName:=urban_circuit"),
			Times:                 sptr("1"),
			Image:                 sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:           sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:      sptr(worldStatsTopic),
			WorldWarmupTopics:     sptr(worldWarmupTopic),
			WorldMaxSimSeconds:    sptr("0"),
			MaxCredits:            nil,
			RequiresQualification: boolptr(false),
			CompetitionDate:       timeptr(time.Now().Add(time.Second * 2)),
		},
		{
			Circuit:            sptr("Urban Practice 3"),
			Worlds:             sptr("worldName:=urban_circuit, worldName:=urban_circuit"),
			Times:              sptr("3, 3"),
			Image:              sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:        sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:   sptr(worldStatsTopic),
			WorldWarmupTopics:  sptr(worldWarmupTopic),
			WorldMaxSimSeconds: sptr("0"),
			MaxCredits:         nil,
			CompetitionDate:    timeptr(time.Now().Add(-time.Hour * 24)),
		},
		{
			Circuit:               sptr("Urban Circuit"),
			Worlds:                sptr("worldName:=urban_circuit"),
			Times:                 sptr("1"),
			Image:                 sptr("infrastructureascode/aws-cli:latest"),
			BridgeImage:           sptr("infrastructureascode/aws-cli:latest"),
			WorldStatsTopics:      sptr(worldStatsTopic),
			WorldWarmupTopics:     sptr(worldWarmupTopic),
			WorldMaxSimSeconds:    sptr("0"),
			Seeds:                 sptr("10"),
			MaxCredits:            nil,
			CompetitionDate:       timeptr(time.Now().Add(time.Hour * 24)),
			RequiresQualification: boolptr(true),
		},
	}

	var availableCircuits []string
	for _, circuit := range circuits {
		globals.Server.Db.Create(circuit)
		availableCircuits = append(availableCircuits, *circuit.Circuit)
	}

	// TODO: This code is specific to the simulations package
	// Temporarily add previous competition circuits to the list of available circuits
	sim.SubTCompetitionCircuits = []string{
		sim.CircuitUrbanCircuit,
		sim.CircuitCaveCircuit,
		sim.CircuitVirtualStixCircuit,
		sim.CircuitVirtualStixCircuit2,
		sim.CircuitFinalsPreliminaryRound,
	}

	for _, circuit := range availableCircuits {
		// We need to check that the circuit is not already in the list to avoid
		// multiple tests from adding the same list of circuits
		if !sim.StrSliceContains(circuit, sim.SubTCircuits) {
			// Circuit is prepended to help StrSliceContains find circuits faster
			sim.SubTCircuits = append([]string{circuit}, sim.SubTCircuits...)
		}
	}

	// Insert qualified teams
	qualifiedTeams := []*sim.SubTQualifiedParticipant{
		{
			Circuit: "Urban Circuit",
			Owner:   "TeamA",
		},
	}

	for _, qualifiedTeam := range qualifiedTeams {
		globals.Server.Db.Create(qualifiedTeam)
	}

	// Reinitialize UsersDB (ie. recreate test users and organization membership)
	// HACK application "subt" shouldn't be hardcoded here.
	mockUsers := users.NewUserAccessorDataMock(ctx, globals.UserAccessor, sysAdminIdentityForTest, "subt")
	mockUsers.ReloadEverything(ctx)
	migrations.DBAddDefaultData(ctx, globals.Server.Db)
}
