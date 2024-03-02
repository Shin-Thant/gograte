package gograte

import (
	"context"
	"database/sql"
	"log"
)

func Status(args []string) {
	driver := args[DRIVER]
	dbURL := args[DB_URL]
	if len(args) > 3 {
		log.Fatalf("Invalid number of arguments.")
	}

	if !validateDbDriver(driver) {
		log.Fatalf("Invalid database driver. Supported databases are: %s\n", DB_DRIVERS.String())
	}
	validatedURL, err := validateDbURL(dbURL)
	if err != nil {
		log.Fatalln("Invalid database URL. Please provide a valid URL.")
	}

	driver = GetSQLDriver(driver)

	db, err := sql.Open(driver, validatedURL.String())
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalf("Error database connection: %v\n", err)
	}
	defer conn.Close()

	records, err := queryAllMigrations(db)
	if err != nil {
		log.Fatalf("Error querying applied migrations: %v\n", err)
	}
	log.Println("\tApplied \tMigration")
	log.Println("\t==============================")
	for _, record := range records {
		log.Printf("\t%v\t\t%d_%s.sql\n", record.IsApplied, record.VersionID, record.Name)
	}
}
