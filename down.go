package gograte

import (
	"database/sql"
	"log"
	"sort"
)

func downAllMigrate(db *sql.DB) {
	records, err := queryAppliedMigrations(db)
	if err != nil {
		log.Fatalf("Error querying applied migrations: %v\n", err)
	}
	if len(records) == 0 {
		log.Fatalln("No migration files found.")
	}

	downMigrate(records, db, false)
}

func downOneMigrate(db *sql.DB) {
	records, err := queryLatestAppliedMigration(db)
	if err != nil {
		log.Fatalf("Error querying applied migrations: %v\n", err)
	}
	if len(records) == 0 {
		log.Fatalln("No migration files found.")
	}

	downMigrate(records, db, true)
}

func downMigrate(records []migrationRecord, db *sql.DB, downOne bool) {
	matches, err := findMigrationFiles()
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	migrationFiles := validateMigrationFilePaths(matches)
	sort.SliceStable(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Timestamp > migrationFiles[j].Timestamp
	})

	var filteredMigrationFiles []migrationFile
	for _, m := range migrationFiles {
		if downOne && len(filteredMigrationFiles) == 1 {
			break
		}

		for _, record := range records {
			if m.Timestamp == record.VersionID {
				filteredMigrationFiles = append(filteredMigrationFiles, m)
				break
			}
		}
	}

	if len(filteredMigrationFiles) == 0 {
		log.Fatalln("No migration to run.")
	}

	startMigrate(filteredMigrationFiles, db, "down")
}

func onDownMigrate(m *migrationFile, tx *sql.Tx) error {
	_, err := updateVersionRecord(m, tx, false)
	if err != nil {
		log.Printf("Error updating version record: %v\n", err)

		err = tx.Rollback()
		if err != nil {
			log.Printf("Error rolling-back transaction: %v\n", err)
		}
		return err
	}
	return err
}
