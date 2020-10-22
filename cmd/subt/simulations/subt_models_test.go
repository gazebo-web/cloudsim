package simulations

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestCountSimulationsByCircuitReturnsZero(t *testing.T) {
	// Get database config
	config, err := ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		t.FailNow()
	}

	// Initialize database
	db, err := ign.InitDbWithCfg(&config)
	if err != nil {
		t.FailNow()
	}

	db.DropTableIfExists(&SimulationDeployment{})

	// Auto migrate simulation deployments
	db.AutoMigrate(&SimulationDeployment{})

	// Define data
	owner := "Ignition Robotics"
	circuit := "Cave Circuit"

	count, err := countSimulationsByCircuit(db, owner, circuit)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 0, *count)
}

func TestCountSimulationByCircuitReturnsCountWhenCircuitIsSubmitted(t *testing.T) {
	// Get database config
	config, err := ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		t.FailNow()
	}

	// Initialize database
	db, err := ign.InitDbWithCfg(&config)
	if err != nil {
		t.FailNow()
	}

	db.DropTableIfExists(&SimulationDeployment{})

	// Auto migrate simulation deployments
	db.AutoMigrate(&SimulationDeployment{})

	// Define data
	owner := "Ignition Robotics"
	circuit := "Cave Circuit"

	gid := "aaaa-bbbb-cccc-eeee"

	// Create the first submission
	first := &SimulationDeployment{
		GroupID:          &gid,
		DeploymentStatus: simPending.ToPtr(),
		Owner:            &owner,
		ExtraSelector:    &circuit,
		Held:             true,
		MultiSim:         int(multiSimParent),
	}
	db.Model(&SimulationDeployment{}).Save(&first)

	// Create child sims for the submission
	createTestChildSims(t, db, first, 3)

	count, err := countSimulationsByCircuit(db, owner, circuit)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 1, *count)
}

func TestCountSimulationByCircuitReturnsCountWhenCircuitGetsSuperseded(t *testing.T) {
	// Get database config
	config, err := ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		t.FailNow()
	}

	// Initialize database
	db, err := ign.InitDbWithCfg(&config)
	if err != nil {
		t.FailNow()
	}

	db.DropTableIfExists(&SimulationDeployment{})

	// Auto migrate simulation deployments
	db.AutoMigrate(&SimulationDeployment{})

	// Define data
	owner := "Ignition Robotics"
	circuit := "Cave Circuit"

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
	db.Model(&SimulationDeployment{}).Save(&first)

	// Create child sims for the submission
	createTestChildSims(t, db, first, 3)

	secondGroupID := "aaaa-bbbb-cccc-eeee"

	// Create the second submission
	second := &SimulationDeployment{
		GroupID:          &secondGroupID,
		DeploymentStatus: simPending.ToPtr(),
		Owner:            &owner,
		ExtraSelector:    &circuit,
		Held:             true,
		MultiSim:         int(multiSimParent),
	}
	db.Model(&SimulationDeployment{}).Save(&second)

	// Create child sims for the submission
	createTestChildSims(t, db, second, 3)

	assert.NoError(t, MarkPreviousSubmissionsSuperseded(db, secondGroupID, owner, circuit))

	count, err := countSimulationsByCircuit(db, owner, circuit)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 1, *count)
}
