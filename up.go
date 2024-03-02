package gograte

import (
	"database/sql"
	"log"
	"sort"
	"strconv"
)

func upAllMigrate(db *sql.DB) {
	upMigrate(db, false)
}

func upOneMigrate(db *sql.DB) {
	upMigrate(db, true)
}

func upMigrate(db *sql.DB, upOne bool) {
	records, err := queryAllMigrations(db)
	if err != nil {
		log.Fatalf("Error querying applied migrations: %v\n", err)
	}

	matches, err := findMigrationFiles()
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	migrationFiles := validateMigrationFilePaths(matches)
	sort.SliceStable(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Timestamp < migrationFiles[j].Timestamp
	})

	var filteredMigrationFiles []migrationFile
	for _, m := range migrationFiles {
		if upOne && len(filteredMigrationFiles) == 1 {
			break
		}
		isInRecord := false

		for _, record := range records {
			if m.Timestamp == record.VersionID {
				isInRecord = true
				if !record.IsApplied {
					filteredMigrationFiles = append(filteredMigrationFiles, m)
				}
				break
			}
		}

		if !isInRecord {
			m.IsNewFile = true
			filteredMigrationFiles = append(filteredMigrationFiles, m)
		}
	}

	if len(filteredMigrationFiles) == 0 {
		log.Fatalln("No migration to run.")
	}

	startMigrate(filteredMigrationFiles, db, "up")
}

func upToMigrate(versionStr string, db *sql.DB) {
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		log.Fatalln("Invalid version number.")
	}

	records, err := queryNonAppliedMigrations(db)
	if err != nil {
		log.Fatalf("Error querying migrations: %v\n", err)
	}
	if len(records) == 0 {
		log.Fatalln("No migration found.")
	}

	recordFound := checkRecordExists(records, version)
	if !recordFound {
		log.Fatalln("Couldn't find version to migrate.")
	}

	filteredRecords := make([]migrationRecord, 0)
	for _, record := range records {
		if record.VersionID == version {
			filteredRecords = append(filteredRecords, record)
			break
		}
		filteredRecords = append(filteredRecords, record)
	}

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
		for _, record := range filteredRecords {
			if m.Timestamp == record.VersionID {
				filteredMigrationFiles = append(filteredMigrationFiles, m)
				break
			}
		}
	}

	if len(filteredMigrationFiles) == 0 {
		log.Fatalln("No migration to run.")
	}

	startMigrate(filteredMigrationFiles, db, "up")
}

func checkRecordExists(records []migrationRecord, version int) bool {
	for _, record := range records {
		if record.VersionID == version {
			return true
		}
	}
	return false
}

func onUpMigrate(m *migrationFile, tx *sql.Tx) error {
	if m.IsNewFile {
		_, err := insertVersionRecord(m, tx)
		if err != nil {
			log.Printf("Error inserting version record: %v\n", err)

			err = tx.Rollback()
			if err != nil {
				log.Printf("Error rolling-back transaction: %v\n", err)
			}
			return err
		}
		return err
	}

	_, err := updateVersionRecord(m, tx, true)
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
