package db

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"log"
)

// Migrate auto migrates database tables
func Migrate(ctx context.Context, db *gorm.DB) {
	// Note about Migration from GORM doc: http://jinzhu.me/gorm/database.html#migration
	//
	// WARNING: AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON'T change existing column's type or delete unused columns to protect your data.
	//

	if db != nil {
		db.AutoMigrate(
		//&sim.Simulation{},
		//&sim.SimulationDeploymentsSubTValue{},
		//&sim.MachineInstance{},
		//&sim.SubTCircuitRules{},
		//&sim.CircuitCustomRule{},
		//&sim.SubTQualifiedParticipant{},
		)

		migrateMultiSimRoles(ctx, db)
		migrateCircuitRules(ctx, db)
		migrateStatusConstantValues(ctx, db)
		migrateSimulationRobots(ctx, db)
	}
}

func migrateMultiSimRoles(ctx context.Context, db *gorm.DB) {
	log.Println("[MIGRATION] initializing MultiSim field in old SimulationDeployments")
	tx := db.Begin()

	if err := tx.Exec("UPDATE simulation_deployments SET multi_sim = 0 WHERE multi_sim IS NULL;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error while running 'UPDATE simulation_deployments SET multi_sim = 0 WHERE multi_sim IS NULL;'", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX to initialize MultiSim field in old SimulationDeployments", err)
	}
}

func migrateCircuitRules(ctx context.Context, db *gorm.DB) {
	log.Println("[MIGRATION] Migrating Worlds and Max Sim seconds in Circuit Rules")
	tx := db.Begin()

	if err := tx.Exec("UPDATE sub_t_circuit_rules SET world_stats_topics='/world/default/stats',world_max_sim_seconds='0' WHERE world_stats_topics IS NULL;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error while Migrating Worlds and Max Sim seconds in Circuit Rules", err)
	}
	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX in Migrating Worlds and Max Sim seconds in Circuit Rules", err)
	}
}

func migrateStatusConstantValues(ctx context.Context, db *gorm.DB) {

	log.Println("[MIGRATION] updating Status constant values. From 0..9 to 0..100")

	// First check if the migration is needed. Otherwise return
	var needToUpdate bool
	rows, err := db.Raw("SELECT * FROM simulation_deployments where deployment_status > 9;").Rows()
	defer rows.Close()
	if err != nil {
		needToUpdate = false
	} else {
		needToUpdate = !rows.Next()
	}

	if !needToUpdate {
		log.Println("[MIGRATION] updating Status constant values. NO needToUpdate")
		return
	}

	tx := db.Begin()

	// We updated all the contant values to their value multiplied by 10. But we
	// also added a new constant in the middle called 'simRunningWithErrors' with
	// value 40. So any old rows with value "4" should now be "50". The same applies
	// to rows with higher values ( >= 4).
	if err := tx.Exec("UPDATE simulation_deployments SET deployment_status = (deployment_status+1)*10 WHERE deployment_status BETWEEN 4 AND 9;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error updating Status constant values. From 0..9 to 0..100", err)
	}
	if err := tx.Exec("UPDATE simulation_deployments SET deployment_status = (deployment_status)*10 WHERE deployment_status <= 3;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error updating Status constant values. From 0..9 to 0..100", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX to update Status constant values", err)
	}
}

func migrateSimulationRobots(ctx context.Context, db *gorm.DB) {

	log.Println("[MIGRATION] Updating Simulation NULL robots values.")

	// First check if the migration is needed, otherwise return.
	var count int
	if err := db.Model(&simulations.Simulation{}).
		Where("robots IS NULL").
		Count(&count).Error; err != nil {
		log.Fatal("[MIGRATION] Migrating Status robots values: could not get number of entries for migration")
	} else if count == 0 {
		log.Println("[MIGRATION] Migrating Status robots values: migration is not required")
		return
	}

	tx := db.Begin()

	// Field `simulation_deployments.robots` should contain a comma-separated list of robot names based on the robots
	// defined in `simulation_deployments.extra`. This implementation relies on the limitations imposed by name
	// validators. If validators change in the future, then this might migrate incorrectly.
	if err := tx.Exec(`
		UPDATE simulation_deployments sd
		JOIN (SELECT id,  
					 REPLACE(
					   REPLACE(
						 SUBSTRING(sd.robots, 2, LENGTH(sd.robots)-2),
					   '"', ''), 
					 ', ', ',') robots
			  FROM (SELECT id, JSON_EXTRACT(extra, '$.robots[*].Name') robots 
					FROM simulation_deployments) sd) r
		  ON sd.id = r.id
		SET sd.robots = r.robots`).Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error updating Status robots values", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX to update Status robots values", err)
	}
}
