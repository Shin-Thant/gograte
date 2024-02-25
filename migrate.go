package gograte

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
)

type migrationFile struct {
	Timestamp int
	Path      string
}

type ValidData []string

func (v *ValidData) String() string {
	output := ""
	length := len(*v)
	for index, item := range *v {
		if index == length-1 {
			output += item
		} else {
			output += item + ", "
		}
	}
	return output
}

// argument indexes
var DRIVER = 1
var DB_URL = 2
var ACTION = 3

func Migrate(args []string) {
	driver := args[DRIVER]
	dbURL := args[DB_URL]
	action := args[ACTION]

	if !validateDbDriver(driver) {
		log.Fatalf("Invalid database driver. Supported databases are: %s\n", DB_DRIVERS.String())
	}
	validatedURL, err := validateDbURL(dbURL)
	if err != nil {
		log.Fatalln("Invalid database URL. Please provide a valid URL.")
	}
	if !validateAction(action) {
		log.Fatalf("Invalid action. Supported actions are: %s\n", ACTIONS.String())
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

	matches, err := findMigrationFiles()
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	_, err = CreateMigrationTableIfNotExist(db)
	if err != nil {
		log.Fatalf("Error creating migration table: %v\n", err)
	}

	migrations := validateMigrationFilePaths(matches)
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].Timestamp < migrations[j].Timestamp
	})

	records, err := queryMigrationRecord(db)
	if err != nil {
		log.Fatalf("Error querying migration records: %v\n", err)
	}

	var filteredMigrationFiles []migrationFile
	for _, m := range migrations {
		if len(records) == 0 && action == "up" {
			filteredMigrationFiles = append(filteredMigrationFiles, m)
			continue
		}

		isInRecord := false

	RecordLoop:
		for _, record := range records {
			fmt.Println(record.VersionID, m.Timestamp, record.VersionID == m.Timestamp)
			switch action {
			case "up":
				if record.VersionID == m.Timestamp {
					isInRecord = true
					if !record.IsApplied {
						filteredMigrationFiles = append(filteredMigrationFiles, m)
					}
					break RecordLoop
				}
			case "down":
				if record.VersionID == m.Timestamp {
					if record.IsApplied {
						filteredMigrationFiles = append(filteredMigrationFiles, m)
					}
					break RecordLoop
				}
			}
		}

		if !isInRecord && action == "up" {
			filteredMigrationFiles = append(filteredMigrationFiles, m)
		}
	}

	if len(filteredMigrationFiles) == 0 {
		log.Println("No migration to run.")
		return
	}

	for _, m := range filteredMigrationFiles {
		file, err := os.Open(m.Path)
		if err != nil {
			log.Fatalf("Error migration opening file: %v\n", err)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil || info.IsDir() {
			continue
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var statement string
		isDownMigrateStarted := false

		for scanner.Scan() {
			line := scanner.Text()
			isDownMigrate := strings.HasPrefix(line, downMigrateComment)

			if action == "up" {
				if isDownMigrate {
					break
				}
				isComment := strings.HasPrefix(line, commentSyntax)
				if isComment {
					continue
				}

				statement += line
			}

			if action == "down" {
				if !isDownMigrate && !isDownMigrateStarted {
					continue
				}
				if !isDownMigrateStarted {
					isDownMigrateStarted = true
					continue
				}
				isComment := strings.HasPrefix(line, "--")
				if isComment {
					continue
				}

				statement += line
			}
		}

		fmt.Println("Running migration: ", m.Path)
		// err = runMigration(db, statement)
		// if err != nil {
		// 	return
		// }

		// _, err = insertVersionRecord(&m, db)
		// if err != nil {
		// 	log.Fatalf("Error inserting migration version: %v\n", err)
		// }
	}

}

var DB_DRIVERS ValidData = []string{"mysql", "postgres", "sqlite3", "mssql"}

func validateDbDriver(inputDriver string) bool {
	for _, name := range DB_DRIVERS {
		if name == inputDriver {
			return true
		}
	}
	return false
}

func validateDbURL(inputDbURL string) (*url.URL, error) {
	u, err := url.ParseRequestURI(inputDbURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

var ACTIONS ValidData = []string{"up", "down"}

func validateAction(inputAction string) bool {
	for _, action := range ACTIONS {
		if action == inputAction {
			return true
		}
	}
	return false
}
