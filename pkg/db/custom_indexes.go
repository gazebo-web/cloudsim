package db

import (
	"context"
	"github.com/jinzhu/gorm"
)

// AddCustomIndexes allows application to add custom indexes that cannot be added automatically
// by GORM.
func AddCustomIndexes(ctx context.Context, db *gorm.DB) {
	// TIP: command to check for existing foreign keys in db:
	// SELECT TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'fuel';
	// db.Model(&users.User{}).AddForeignKey("username", "unique_owners(name)", "RESTRICT", "RESTRICT")

	// TIP: You can check created indexes by executing in mysql: `show index from models;`

	// Just an example:
	// First add indexes for Models
	// found, err := indexIsPresent(db, "models", "models_fultext")
	// if err != nil {
	// 	ign.LoggerFromContext(ctx).Critical("Error with db while checking index", err)
	// 	log.Fatal("Error with db while checking index", err)
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
