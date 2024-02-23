package gograte

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

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
	fmt.Println("Driver:", driver)
	fmt.Println("Database URL:", validatedURL)
	fmt.Println("Action:", action)

	db, err := sql.Open(driver, dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalf("Error database connection: %v\n", err)
	}
	defer conn.Close()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}
	matches, err := filepath.Glob(path.Join(currentDir, "*.sql"))
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	result, err := db.Exec(`CREATE TABLE IF NOT EXISTS _gograte_db_versions (
		id VARCHAR(255) PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Error creating migration table: %v\n", err)
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())

	// rows, err := db.Query("SELECT * FROM _gograte_db_versions")
	// if err != nil {
	// 	log.Fatalf("Error querying database versions: %v\n", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	fmt.Println(rows.Columns())
	// }
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
