package gograte

import (
	"bufio"
	"context"
	"database/sql"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
)

type migrationFile struct {
	Timestamp int
	FileName  string
	Path      string
	IsNewFile bool
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

	_, err = CreateMigrationTableIfNotExist(db)
	if err != nil {
		log.Fatalf("Error creating migration table: %v\n", err)
	}

	switch action {
	case "up":
		upAllMigrate(db)
	case "up-one":
		upOneMigrate(db)
	}
	return

	matches, err := findMigrationFiles()
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	migrations := validateMigrationFilePaths(matches)
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].Timestamp < migrations[j].Timestamp
	})

	records, err := queryAllMigrations(db)
	if err != nil {
		log.Fatalf("Error querying migration records: %v\n", err)
	}

	var filteredMigrationFiles []migrationFile
	for _, m := range migrations {
		if len(records) == 0 && action == "up" {
			m.IsNewFile = true
			filteredMigrationFiles = append(filteredMigrationFiles, m)
			continue
		}

		isInRecord := false

	RecordLoop:
		for _, record := range records {
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
			m.IsNewFile = true
			filteredMigrationFiles = append(filteredMigrationFiles, m)
		}
	}

	if len(filteredMigrationFiles) == 0 {
		log.Println("No migration to run.")
		return
	}

	startMigrate(filteredMigrationFiles, db, action)
}

func startMigrate(migrationFiles []migrationFile, db *sql.DB, action string) {
	for _, m := range migrationFiles {
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
		hasPassedDownMigrateCmt := false

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
				if !isDownMigrate && !hasPassedDownMigrateCmt {
					continue
				}
				if !hasPassedDownMigrateCmt {
					hasPassedDownMigrateCmt = true
					continue
				}
				isComment := strings.HasPrefix(line, commentSyntax)
				if isComment {
					continue
				}

				statement += line
			}
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Error beginning transaction: %v\n", err)
		}

		_, err = runMigration(tx, statement)
		if err != nil {
			log.Printf("Error executing migration: %v\n", err)

			err = tx.Rollback()
			if err != nil {
				log.Printf("Error rolling-back transaction: %v\n", err)
			}
			return
		}

		switch action {
		case "up":
			err = onUpMigrate(&m, tx)
		case "down":
			err = onDownMigrate(&m, tx)
		}
		if err != nil {
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Fatalf("Error committing transaction: %v\n", err)
			return
		}

		log.Println("OK    ", m.FileName)
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

var ACTIONS ValidData = []string{"up", "down", "up-one"}

func validateAction(inputAction string) bool {
	for _, action := range ACTIONS {
		if action == inputAction {
			return true
		}
	}
	return false
}
