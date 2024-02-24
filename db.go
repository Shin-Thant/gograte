package gograte

import (
	"database/sql"
	"log"
)

func GetSQLDriver(driver string) string {
	switch driver {
	case "mssql":
		return "sqlserver"
	case "sqlite3":
		return "sqlite"
	case "postgres", "redshift":
		return "pgx"
	}
	return ""
}

func CreateMigrationTableIfNotExist(db *sql.DB) (sql.Result, error) {
	return db.Exec(`CREATE TABLE IF NOT EXISTS _gograte_db_versions (
		id VARCHAR(255) PRIMARY KEY,
		version_id BIGINT NOT NULL,
		is_applied BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)
}

func runMigration(db *sql.DB, statement string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return err
	}
	_, err = db.Exec(statement)
	if err != nil {
		log.Printf("Error executing migration: %v\n", err)

		err = tx.Rollback()
		if err != nil {
			log.Printf("Error rolling-back transaction: %v\n", err)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return err
	}

	return nil
}

// func CheckTableExists(db *sql.DB) bool {
// 	query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '_gograte_db_versions')"
// 	var exists bool
// 	db.QueryRow(query).Scan(&exists)
// 	return exists
// }
