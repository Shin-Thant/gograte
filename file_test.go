package gograte_test

import (
	"testing"

	"github.com/Shin-Thant/gograte"
)

func TestValidateMigrationFilePaths(t *testing.T) {
	mockPaths := []string{
		"migrations/0001_create_table.sql",
		"migrations/0002_create_table.sql",
		"migrations/0003_create_table.sql",
	}

	result := gograte.ValidateMigrationFilePaths(mockPaths)
	if len(mockPaths) != len(result) {
		t.Errorf("Expected %d, got %d", len(mockPaths), len(result))
	}
	for _, m := range result {
		if m.Timestamp == 0 {
			t.Errorf("Expected timestamp to be non-zero")
		}
		if m.Name == "" {
			t.Errorf("Expected name to be non-empty")
		}
		if m.Path == "" {
			t.Errorf("Expected path to be non-empty")
		}
		if m.IsNewFile {
			t.Errorf("Expected IsNewFile to be false")
		}
	}

	mockPaths = []string{
		"migrations/0001_create_table.sql",
		"migrations/0002_create_table.sql",
		"migrations/create_table.sql",
	}
	result = gograte.ValidateMigrationFilePaths(mockPaths)
	if len(mockPaths)-1 != len(result) {
		t.Errorf("Expected %d, got %d", len(mockPaths)-1, len(result))
	}
	for _, m := range result {
		if m.Timestamp == 0 {
			t.Errorf("Expected timestamp to be non-zero")
		}
		if m.Name == "" {
			t.Errorf("Expected name to be non-empty")
		}
		if m.Path == "" {
			t.Errorf("Expected path to be non-empty")
		}
		if m.IsNewFile {
			t.Errorf("Expected IsNewFile to be false")
		}
	}
}
