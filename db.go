package gograte

import (
	"database/sql"
	"fmt"

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
		version_id BIGINT UNIQUE NOT NULL,
		is_applied BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)
}

func runMigration(tx *sql.Tx, statement string) (sql.Result, error) {
	return tx.Exec(statement)
}

func insertVersionRecord(m *migrationFile, tx *sql.Tx) (sql.Result, error) {
	id := uuid.New()
	query := fmt.Sprintf(`INSERT INTO _gograte_db_versions (id, version_id) VALUES ('%s', %d);`, id.String(), m.Timestamp)
	return tx.Exec(query)
}

func updateVersionRecord(m *migrationFile, tx *sql.Tx, isApplied bool) (sql.Result, error) {
	query := fmt.Sprintf(`UPDATE _gograte_db_versions SET is_applied = %t WHERE version_id = %d;`, isApplied, m.Timestamp)
	return tx.Exec(query)
}

func queryMigrationRecord(db *sql.DB) ([]migrationRecord, error) {
	rows, err := db.Query(`
	SELECT * FROM _gograte_db_versions
	ORDER BY created_at ASC;
	`)
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
