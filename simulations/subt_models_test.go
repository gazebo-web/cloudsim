package simulations

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestSubTCreateSimulation_robotImageBelongsToECROwner(t *testing.T) {
	test := func(valid bool, owner string, robotImages []string) {
		cs := &SubTCreateSimulation{
			CreateSimulation: CreateSimulation{
				Owner: owner,
			},
			RobotImage: robotImages,
		}
		assert.Equal(t, valid, cs.robotImagesBelongToECROwner())
	}

	subTRedTeam := "OSRF SubT RedTeam"

	// If all images are stored in the team's ECR repo this should succeed
	test(true, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:subt_seed",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"200670743174.dkr.ecr.us-east-2.amazonaws.com/osrf_subt_redteam:virtual_stix_arl",
	})

	// A non-ECR repo should always succeed
	test(true, subTRedTeam, []string{
		"https://hub.docker.com/r/osrf/subt-virtual-testbed",
	})

	// A mix between images stored in the team's ECR repo and in non-ECR repos should succeed
	test(true, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:subt_seed",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"https://hub.docker.com/r/osrf/subt-virtual-testbed",
	})

	// If an image is stored in an ECR repo from another team then it should fail
	test(false, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/subt_sim:cloudsim_sim_latest",
		"200670743174.dkr.ecr.us-east-2.amazonaws.com/osrf_subt_redteam:virtual_stix_arl",
	})
}


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
