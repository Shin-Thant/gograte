package gograte

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
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
	matches, err := filepath.Glob(path.Join(currentDir, "migrations", "*.sql"))
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}
	if len(matches) == 0 {
		log.Fatalln("No migration files found.")
	}

	_, err = CreateTableIfNotExist(db)
	if err != nil {
		log.Fatalf("Error creating migration table: %v\n", err)
	}

	migrations := make([]migrationFile, len(matches))
	for index, path := range matches {
		targetFile := filepath.Base(path)
		fileSlice := strings.Split(targetFile, "_")
		if len(fileSlice) != 2 {
			continue
		}
		numericPart := fileSlice[0]
		result, err := strconv.Atoi(numericPart)
		if err != nil {
			continue
		}
		migrations[index] = migrationFile{
			Timestamp: result,
			Path:      path,
		}
	}
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].Timestamp < migrations[j].Timestamp
	})

	for _, m := range migrations {
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
		downMigrateStart := false

		for scanner.Scan() {
			line := scanner.Text()
			isDownMigrate := strings.HasPrefix(line, "-- +gograte Down")

			if action == "up" {
				if isDownMigrate {
					break
				}
				isComment := strings.HasPrefix(line, "--")
				if isComment {
					continue
				}

				statement += line
			}

			if action == "down" {
				if !isDownMigrate && !downMigrateStart {
					continue
				}
				if !downMigrateStart {
					downMigrateStart = true
					continue
				}
				isComment := strings.HasPrefix(line, "--")
				if isComment {
					continue
				}

				statement += line
			}
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error beginning transaction: %v\n", err)
			return
		}
		_, err = db.Exec(statement)
		if err != nil {
			log.Printf("Error executing migration: %v\n", err)

			err = tx.Rollback()
			if err != nil {
				log.Printf("Error rolling back transaction: %v\n", err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			return
		}
	}

	// rows, err := db.Query("SELECT * FROM _gograte_db_versions;")
	// if err != nil {
	// 	log.Fatalf("Error querying database versions: %v\n", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	res, err := rows.Columns()
	// 	fmt.Println(res, err)
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
