package gograte

import (
	"database/sql"
	"log"
)

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
