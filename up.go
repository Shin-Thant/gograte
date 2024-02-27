package gograte

import (
	"database/sql"
	"log"
)

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
