package migrations

import (
	"context"
	"github.com/jinzhu/gorm"
)

// DropModels drops all tables from db. Used by tests.
func DropModels(ctx context.Context, db *gorm.DB) {

	if db != nil {
		// IMPORTANT NOTE: DROP TABLE order is important, due to FKs
		db.DropTableIfExists(
		//&sim.SimulationDeploymentsSubTValue{},
		//&sim.Simulation{},
		//&sim.MachineInstance{},
		//&sim.SubTCircuitRules{},
		//&sim.CircuitCustomRule{},
		//&sim.SubTQualifiedParticipant{},
		)

		// Now also remove many_to_many tables, because they are not automatically removed.
		// db.DropTableIfExists("model_tags", "world_tags")
	}
}
