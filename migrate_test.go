package gograte_test

import (
	"testing"

	"github.com/Shin-Thant/gograte"
)

func TestValidateDbDriver(t *testing.T) {
	// validate supported drivers
	inputDrivers := []string{"sqlite3", "postgres", "mysql"}
	for _, driver := range inputDrivers {
		result := gograte.ValidateDbDriver(driver)
		if !result {
			t.Errorf("Expected true but got false for driver: %s", driver)
		}
	}

	// validate unsupported drivers
	unsupportedDrivers := []string{"sqlite", "pgx", "mysqlx", "mssql"}
	for _, driver := range unsupportedDrivers {
		result := gograte.ValidateDbDriver(driver)
		if result {
			t.Errorf("Expected false but got true for driver: %s", driver)
		}
	}
}
