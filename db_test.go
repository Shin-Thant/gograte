package gograte_test

import (
	"testing"

	"github.com/Shin-Thant/gograte"
)

func TestGetSQLDriver(t *testing.T) {
	inputDriver := "invalid one"
	mappedDriver := gograte.GetSQLDriver(inputDriver)
	if mappedDriver != "" {
		t.Error("Result should be empty string.")
	}
}
