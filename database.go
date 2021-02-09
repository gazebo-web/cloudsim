package main

// Import this file's dependencies
import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"log"
)

// DBMigrate auto migrates database tables
func DBMigrate(ctx context.Context, db *gorm.DB) {
	// Note about Migration from GORM doc: http://jinzhu.me/gorm/database.html#migration
	//
	// WARNING: AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON'T change existing column's type or delete unused columns to protect your data.
	//

	if db != nil {
		db.AutoMigrate(
			&sim.SimulationDeployment{},
			&sim.SimulationDeploymentsSubTValue{},
			&sim.MachineInstance{},
			&sim.SubTCircuitRules{},
			&sim.CircuitCustomRule{},
			&sim.SubTQualifiedParticipant{},
		)

		migrateMultiSimRoles(ctx, db)
		migrateCircuitRules(ctx, db)
		migrateDeploymentStatusConstantValues(ctx, db)
		migrateSimulationDeploymentRobots(ctx, db)
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

func migrateDeploymentStatusConstantValues(ctx context.Context, db *gorm.DB) {

	log.Println("[MIGRATION] updating DeploymentStatus constant values. From 0..9 to 0..100")

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
		log.Println("[MIGRATION] updating DeploymentStatus constant values. NO needToUpdate")
		return
	}

	tx := db.Begin()

	// We updated all the contant values to their value multiplied by 10. But we
	// also added a new constant in the middle called 'simRunningWithErrors' with
	// value 40. So any old rows with value "4" should now be "50". The same applies
	// to rows with higher values ( >= 4).
	if err := tx.Exec("UPDATE simulation_deployments SET deployment_status = (deployment_status+1)*10 WHERE deployment_status BETWEEN 4 AND 9;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error updating DeploymentStatus constant values. From 0..9 to 0..100", err)
	}
	if err := tx.Exec("UPDATE simulation_deployments SET deployment_status = (deployment_status)*10 WHERE deployment_status <= 3;").Error; err != nil {
		tx.Rollback()
		log.Fatal("[MIGRATION] Error updating DeploymentStatus constant values. From 0..9 to 0..100", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX to update DeploymentStatus constant values", err)
	}
}

func migrateSimulationDeploymentRobots(ctx context.Context, db *gorm.DB) {

	log.Println("[MIGRATION] Updating SimulationDeployment NULL robots values.")

	// First check if the migration is needed, otherwise return.
	var count int
	if err := db.Model(&sim.SimulationDeployment{}).
		Where("robots IS NULL").
		Count(&count).Error; err != nil {
		log.Fatal("[MIGRATION] Migrating DeploymentStatus robots values: could not get number of entries for migration")
	} else if count == 0 {
		log.Println("[MIGRATION] Migrating DeploymentStatus robots values: migration is not required")
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
		log.Fatal("[MIGRATION] Error updating DeploymentStatus robots values", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("[MIGRATION] Error while committing TX to update DeploymentStatus robots values", err)
	}
}

// DBDropModels drops all tables from DB. Used by tests.
func DBDropModels(ctx context.Context, db *gorm.DB) {

	if db != nil {
		// IMPORTANT NOTE: DROP TABLE order is important, due to FKs
		db.DropTableIfExists(
			&sim.SimulationDeploymentsSubTValue{},
			&sim.SimulationDeployment{},
			&sim.MachineInstance{},
			&sim.SubTCircuitRules{},
			&sim.CircuitCustomRule{},
			&sim.SubTQualifiedParticipant{},
		)

		// Now also remove many_to_many tables, because they are not automatically removed.
		// db.DropTableIfExists("model_tags", "world_tags")
	}
}

// DBAddDefaultData adds default data.
func DBAddDefaultData(ctx context.Context, db *gorm.DB) {

	// if db != nil {
	// }
}

// DBAddCustomIndexes allows application to add custom indexes that cannot be added automatically
// by GORM.
func DBAddCustomIndexes(ctx context.Context, db *gorm.DB) {
	// TIP: command to check for existing foreign keys in db:
	// SELECT TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'fuel';
	// db.Model(&users.User{}).AddForeignKey("username", "unique_owners(name)", "RESTRICT", "RESTRICT")

	// TIP: You can check created indexes by executing in mysql: `show index from models;`

	// Just an example:
	// First add indexes for Models
	// found, err := indexIsPresent(db, "models", "models_fultext")
	// if err != nil {
	// 	ign.LoggerFromContext(ctx).Critical("Error with DB while checking index", err)
	// 	log.Fatal("Error with DB while checking index", err)
	// 	return
	// }
	// if !found {
	// 	db.Exec("ALTER TABLE models ADD FULLTEXT models_fultext (name, description);")
	// 	db.Exec("ALTER TABLE tags ADD FULLTEXT tags_fultext (name);")
	// }

}

// indexIsPresent returns true if the index with name idxName already exists in the given table
func indexIsPresent(db *gorm.DB, table string, idxName string) (bool, error) {
	// Raw SQL
	rows, err := db.Raw("select * from information_schema.statistics where table_schema=database() and table_name=? and index_name=?;",
		table, idxName).Rows() //(*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}
