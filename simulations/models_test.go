package simulations

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strconv"

	// "strconv"
	"testing"
	"time"
)

func TestSimulationDeployment_Clone(t *testing.T) {
	simDep := &SimulationDeployment{
		ID:               123,
		CreatedAt:        time.Now().Add(-time.Hour * 24),
		UpdatedAt:        time.Now().Add(-time.Hour * 24),
		DeletedAt:        timeptr(time.Now().Add(-time.Hour * 12)),
		StoppedAt:        timeptr(time.Now().Add(-time.Hour * 12)),
		ValidFor:         sptr("6h0m0s"),
		Owner:            sptr("Test"),
		Creator:          sptr("TestUser"),
		Private:          boolptr(true),
		StopOnEnd:        boolptr(false),
		Name:             sptr("TestSimDep"),
		Image:            sptr("test"),
		GroupID:          sptr("11111111-1111-1111-1111-111111111111-c-1"),
		ParentGroupID:    sptr("11111111-1111-1111-1111-111111111111-c-1"),
		MultiSim:         2,
		DeploymentStatus: intptr(90),
		ErrorStatus:      sptr("InitializationFailed"),
		Platform:         sptr("subt"),
		Application:      sptr("subt"),
		Extra:            sptr("{}"),
		ExtraSelector:    sptr("Test Circuit"),
		Robots:           sptr("X1,X2"),
	}

	simDepClone := simDep.Clone()
	simDepClone.GroupID = sptr(fmt.Sprintf("%s-r-1", *simDep.GroupID))

	// Check that the model fields have been cleared
	assert.Equal(t, uint(0), simDepClone.ID)
	assert.Nil(t, simDepClone.DeletedAt)
	assert.Nil(t, simDepClone.StoppedAt)

	// Check that the references are copied and can be overwritten
	assert.NotEqual(t, *simDep.GroupID, *simDepClone.GroupID)
}

func TestGetRemainingSubmissions(t *testing.T) {
	// Get database config
	config, err := ign.NewDatabaseConfigFromEnvVars()
	require.NoError(t, err)

	// Initialize database
	db, err := ign.InitDbWithCfg(&config)
	require.NoError(t, err)

	db.DropTableIfExists(&SimulationDeployment{})
	db.DropTableIfExists(&CircuitCustomRule{})

	// Auto migrate models
	db.AutoMigrate(&SimulationDeployment{})
	db.AutoMigrate(&CircuitCustomRule{})

	// Define data
	owner := "Ignition Robotics"
	circuit := "Cave Circuit"

	require.NoError(t, db.Model(&CircuitCustomRule{}).Save(&CircuitCustomRule{Circuit: &circuit, Value: "1", RuleType: MaxSubmissions}).Error)

	result, err := getRemainingSubmissions(db, circuit, "")
	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, 1, *result)

	// Define group ID of the first submission
	firstGroupID := "aaaa-bbbb-cccc-dddd"

	// Create the first submission
	first := &SimulationDeployment{
		GroupID:          &firstGroupID,
		DeploymentStatus: simPending.ToPtr(),
		Owner:            &owner,
		ExtraSelector:    &circuit,
		Held:             true,
		MultiSim:         int(multiSimParent),
	}

	require.NoError(t, db.Model(&SimulationDeployment{}).Save(&first).Error)

	// Create child sims for the first submission
	createTestChildSims(t, db, first, 3)

	result, err = getRemainingSubmissions(db, circuit, owner)
	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, 0, *result)

	gid := "aaaa-bbbb-cccc-eeee"
	second := &SimulationDeployment{
		GroupID:          &gid,
		DeploymentStatus: simPending.ToPtr(),
		Owner:            &owner,
		ExtraSelector:    &circuit,
		Held:             true,
		MultiSim:         int(multiSimParent),
	}
	db.Model(&SimulationDeployment{}).Save(&second)

	// Create child sims for the second submission
	createTestChildSims(t, db, second, 3)

	require.NoError(t, MarkPreviousSubmissionsSuperseded(db, gid, owner, circuit))

	result, err = getRemainingSubmissions(db, circuit, owner)
	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, 0, *result)
}

func TestMachineInstanceTestSuite(t *testing.T) {
	suite.Run(t, &MachineInstanceTestSuite{})
}

type MachineInstanceTestSuite struct {
	suite.Suite
	ctx context.Context
	db  *gorm.DB
	// Application the machines are created for
	application string
	// Number of test machines to create by default
	machineCount int
	// List of machines created by default
	machines MachineInstances
}

func (suite *MachineInstanceTestSuite) SetupSuite() {
	var err error

	// Get the DB connection
	suite.db, err = gormUtils.GetTestDBFromEnvVars()
	suite.NoError(err)

	// Create machine instances
	suite.application = "test"
	suite.machineCount = 5
}

func (suite *MachineInstanceTestSuite) SetupTest() {
	// Clear and create MachineInstance models
	gormUtils.CleanAndMigrateModels(suite.db, &MachineInstance{})

	suite.machines = suite.createMachineInstances(suite.machineCount)
}

func (suite *MachineInstanceTestSuite) getMachineInstances() MachineInstances {
	machines := MachineInstances{}
	err := suite.db.Find(&machines).Error
	suite.NoError(err)

	return machines
}

func (suite *MachineInstanceTestSuite) createMachineInstances(n int) MachineInstances {
	machines := make(MachineInstances, n)
	for i := 0; i < n; i++ {
		groupID := strconv.Itoa(i + 1)
		instanceID := fmt.Sprintf("i-%d", i+1)

		instance := MachineInstance{
			InstanceID:  &instanceID,
			GroupID:     &groupID,
			Application: &suite.application,
		}

		err := suite.db.Create(&instance).Error
		suite.NoError(err)

		machines[i] = instance
	}

	return machines
}

func (suite *MachineInstanceTestSuite) TestMachineInstances_getInstanceIDs() {
	instanceIDs := []string{
		"i-1",
		"i-2",
		"i-3",
		"i-4",
		"i-5",
	}
	result := suite.machines.getInstanceIDs()

	suite.Equal(len(instanceIDs), len(result))
	for i, id := range instanceIDs {
		suite.Equal(id, *result[i])
	}
}

func (suite *MachineInstanceTestSuite) TestMachineInstances_updateMachinesStatus() {
	machines := suite.getMachineInstances()

	// Machines should have nil last known status
	suite.Equal(suite.machineCount, len(machines))
	for _, machine := range machines {
		suite.Nil(machine.LastKnownStatus)
	}

	// Update machine status
	suite.Nil(suite.machines.updateMachineStatuses(suite.ctx, suite.db, macRunning))

	// Validate new status
	machines = suite.getMachineInstances()
	suite.Equal(suite.machineCount, len(machines))
	for _, machine := range machines {
		suite.Equal(*machine.LastKnownStatus, *macRunning.ToStringPtr())
	}
}


func TestUpdateScore(t *testing.T) {
	// Get database config
	config, err := ign.NewDatabaseConfigFromEnvVars()
	require.NoError(t, err)

	// Initialize database
	db, err := ign.InitDbWithCfg(&config)
	require.NoError(t, err)

	require.NoError(t, db.DropTableIfExists(&SimulationDeployment{}).Error)

	// Auto migrate models
	require.NoError(t, db.AutoMigrate(&SimulationDeployment{}).Error)

	// Scores to use
	score := 1.123
	updatedScore := 99.99

	// Create a submission
	simDep := SimulationDeployment{
		ValidFor:         sptr("6h0m0s"),
		Owner:            sptr("Test"),
		Creator:          sptr("TestUser"),
		Private:          boolptr(true),
		StopOnEnd:        boolptr(false),
		Name:             sptr("TestSimDep"),
		Image:            sptr("test"),
		GroupID:          sptr("11111111-1111-1111-1111-111111111111-c-1"),
		ParentGroupID:    sptr("11111111-1111-1111-1111-111111111111-c-1"),
		MultiSim:         2,
		DeploymentStatus: intptr(90),
		ErrorStatus:      sptr("InitializationFailed"),
		Platform:         sptr("subt"),
		Application:      sptr("subt"),
		Extra:            sptr("{}"),
		ExtraSelector:    sptr("Test Circuit"),
		Robots:           sptr("X1,X2"),
		Score:            &score,
	}

	require.NoError(t, db.Model(&SimulationDeployment{}).Save(&simDep).Error)

	// Verify unmodified score.
	assert.Equal(t, simDep.Score, &score)

	// Update score and verify both variable and database row values.
	simDep.UpdateScore(db, &updatedScore)

	assert.Equal(t, simDep.Score, &updatedScore)

	var dbSimDep SimulationDeployment
	require.NoError(t, db.Model(&SimulationDeployment{}).Where("id = ?", simDep.ID).First(&dbSimDep).Error)
	assert.Equal(t, dbSimDep.Score, &updatedScore)
}
