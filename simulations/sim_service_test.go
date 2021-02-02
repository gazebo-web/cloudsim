package simulations

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"testing"
	"time"
)

func TestMarkPreviousSubmissionsSuperseded(t *testing.T) {
	// Get database config
	db, err := gormUtils.GetTestDBFromEnvVars()

	db.DropTableIfExists(&SimulationDeployment{})

	// Auto migrate simulation deployments
	db.AutoMigrate(&SimulationDeployment{})

	// Define data
	owner := "Ignition Robotics"
	circuit := "Cave Circuit"

	// Define group ID of the previous submission
	previousGroupID := "aaaa-bbbb-cccc-dddd"

	// Create the first submission
	first := &SimulationDeployment{
		GroupID:          &previousGroupID,
		DeploymentStatus: simPending.ToPtr(),
		Owner:            &owner,
		ExtraSelector:    &circuit,
		Held:             true,
		MultiSim:         int(multiSimParent),
	}
	db.Model(&SimulationDeployment{}).Save(&first)

	// Create child sims for the first submission
	createTestChildSims(t, db, first, 3)

	// Create the second submission
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

	// Mark previous as superseded
	assert.NoError(t, MarkPreviousSubmissionsSuperseded(db, gid, owner, circuit))

	// Get the list of previous submissions
	var previousSubmissions SimulationDeployments
	err = db.Model(&SimulationDeployment{}).
		Where("group_id NOT LIKE ?", fmt.Sprintf("%s%%", gid)).
		Where("owner = ?", owner).
		Where("extra_selector = ?", circuit).
		Find(&previousSubmissions).Error

	assert.NoError(t, err)

	// Check that the previous submissions have the superseded status.
	for _, s := range previousSubmissions {
		assert.True(t, simSuperseded.Eq(*s.DeploymentStatus))
	}

	// Get the list of simulations from the latest submission
	var lastSubmission SimulationDeployments
	err = db.Model(&SimulationDeployment{}).
		Where("group_id LIKE ?", fmt.Sprintf("%s%%", gid)).
		Where("owner = ?", owner).
		Where("extra_selector = ?", circuit).
		Find(&lastSubmission).Error
	assert.NoError(t, err)

	// Check that all the simulations in the latest submission don't have the superseded status.
	for _, s := range lastSubmission {
		assert.False(t, simSuperseded.Eq(*s.DeploymentStatus))
	}
}

func createTestChildSims(t *testing.T, db *gorm.DB, sim *SimulationDeployment, amount int) {
	for i := 1; i <= amount; i++ {
		var child SimulationDeployment
		child = *sim
		groupID := fmt.Sprintf("%s-c-%d", *sim.GroupID, i)
		child.ID = 0
		child.CreatedAt = time.Now()
		child.UpdatedAt = time.Now()
		child.GroupID = &groupID
		child.ParentGroupID = sim.GroupID
		child.DeploymentStatus = simPending.ToPtr()
		child.MultiSim = int(multiSimChild)
		err := db.Model(&SimulationDeployment{}).Save(&child).Error
		if err != nil {
			t.FailNow()
		}
	}
}
