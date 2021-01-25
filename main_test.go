package main

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"reflect"
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
	// Use the K8Mock as default K8 test mock
	// Individual tests can change this.
	kcli := NewK8Mock(logCtx)
	sim.AssertMockedClientset(globals.KClientset).SetImpl(kcli)
	//  Mock the Sleep function so that it returns instantly
	sim.Sleep = func(d time.Duration) {
		return
	}
	// Mock WaitForMatchPodsCondition so that the calls return instantly and no
	// wait is needed
	sim.WaitForMatchPodsCondition = func(ctx context.Context, c kubernetes.Interface,
		namespace string, opts metav1.ListOptions, condStr string, timeout time.Duration,
		condition sim.PodCondition) error {
		return nil
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
	globals.TransportTestMock.On("Disconnect").Return()
	globals.TransportTestMock.On("IsConnected").Return(true)
}

// useEC2NodeManager forces the application to use an Ec2Client object as the
// node manager, meant for testing. A function is returned to restore the node
// manager to its original value. The returned function should be deferred right
// after calling this function.
func useEC2NodeManager() func() {
	service := sim.SimServImpl.(*sim.Service)
	nm := service.GetNodeManager()

	// Only replace the node manager if it's not an Ec2Client
	nmType := reflect.TypeOf(nm).String()
	if nmType != "*simulations.Ec2Client" {
		logCtx := context.Background()
		ec2nm, err := sim.NewEC2Client(logCtx, globals.KClientset, globals.EC2Svc)
		if err != nil {
			log.Fatal("Could not create EC2 client. Error:", err)
		}
		for _, application := range service.GetApplications() {
			ec2nm.RegisterPlatform(logCtx, application.(sim.PlatformType))
		}
		service.SetNodeManager(ec2nm)
	}

	return func() {
		service.SetNodeManager(nm)
	}
}

// Clean up our mess
func packageTearDown(ctx context.Context) {
	if ctx == nil {
		ctx = ign.NewContextWithLogger(context.Background(), ign.NewLoggerNoRollbar("test", ign.VerbosityDebug))
	}

	// Restore settings set by tests
	// TODO This should be restored automatically depending on the test
	sim.MaxAWSRetries = 8
	sim.SimServImpl.(*sim.Service).AllowRequeuing = true

	cleanDBTables(ctx)
}

func cleanDBTables(ctx context.Context) {
	DBDropModels(ctx, globals.Server.Db)
}

func createDBTablesAndData(ctx context.Context, worldStatsTopic, worldWarmupTopic string) {
	DBMigrate(ctx, globals.Server.Db)
	// After removing tables we can ask casbin to re initialize
	if err := globals.Permissions.Reload(sysAdminForTest); err != nil {
		log.Fatal("Error reloading casbin policies", err)
	}
	// Apply custom indexes. Eg: fulltext indexes
	DBAddCustomIndexes(ctx, globals.Server.Db)

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
	DBAddDefaultData(ctx, globals.Server.Db)
}
