package gograte_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Shin-Thant/gograte"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func TestGetSQLDriver(t *testing.T) {
	drivers := make(map[string]string)
	drivers["sqlite3"] = "sqlite"
	drivers["postgres"] = "pgx"
	drivers["mysql"] = "mysql"

	for driver, expected := range drivers {
		mappedDriver := gograte.GetSQLDriver(driver)
		if mappedDriver != expected {
			t.Errorf("Expected %s but got %s.", expected, mappedDriver)
		}
	}
}

func TestCreateMigrationTableIfNotExist(t *testing.T) {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Error loading .env file")
		}
		DATABASE_URL = os.Getenv("DATABASE_URL")
		if DATABASE_URL == "" {
			log.Fatalln("DATABASE_URL is not set.")
		}
	}

	// setup database connection
	db, err := sql.Open("pgx", DATABASE_URL)
	if err != nil {
		log.Fatalf("DB test: Error opening database: %v\n", err)
	}
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalf("DB test: Error database connection: %v\n", err)
	}
	defer conn.Close()

	_, err = gograte.CreateMigrationTableIfNotExist(db)
	if err != nil {
		t.Errorf("Should not throw error when creating migration table: %v", err)
	}

	_, err = gograte.CreateMigrationTableIfNotExist(db)
	if err != nil {
		t.Errorf("Should not throw error when creating migration table: %v", err)
	}
}
