package gograte

import (
	"database/sql"
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

func CreateTableIfNotExist(db *sql.DB) (sql.Result, error) {
	return db.Exec(`CREATE TABLE IF NOT EXISTS _gograte_db_versions (
		id VARCHAR(255) PRIMARY KEY,
		version_id VARCHAR(255) NOT NULL,
		is_applied BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)
}

// func CheckTableExists(db *sql.DB) bool {
// 	query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '_gograte_db_versions')"
// 	var exists bool
// 	db.QueryRow(query).Scan(&exists)
// 	return exists
// }
