package gograte

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type migrationRecord struct {
	ID        string
	VersionID int
	IsApplied bool
	CreatedAt string
}

func NewMigrationRecord() *migrationRecord {
	return &migrationRecord{}
}

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

func insertVersionRecord(m *migrationFile, db *sql.DB) (sql.Result, error) {
	id := uuid.New()
	query := fmt.Sprintf(`INSERT INTO _gograte_db_versions (id, version_id) VALUES ('%s', %d);`, id.String(), m.Timestamp)
	return db.Exec(query)
}

func queryMigrationRecord(db *sql.DB) ([]migrationRecord, error) {
	rows, err := db.Query("SELECT * FROM _gograte_db_versions;")
	if err != nil {
		return nil, fmt.Errorf("error querying database versions: %v", err)
	}
	defer rows.Close()

	var records []migrationRecord
	for rows.Next() {
		record := NewMigrationRecord()
		err = rows.Scan(&record.ID, &record.VersionID, &record.IsApplied, &record.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning migration record: %v", err)
		}
		records = append(records, *record)
	}
	return records, nil
}
