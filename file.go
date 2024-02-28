package gograte

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func findMigrationFiles() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return filepath.Glob(path.Join(currentDir, "migrations", "*.sql"))
}

var timestampIdx = 0
var nameIdx = 1

func validateMigrationFilePaths(paths []string) []migrationFile {
	migrations := make([]migrationFile, len(paths))

	for index, path := range paths {
		targetFile := filepath.Base(path)
		fileSlice := strings.SplitN(targetFile, "_", 2)
		if len(fileSlice) != 2 {
			continue
		}

		nameWithExt := fileSlice[nameIdx]
		name := strings.TrimSuffix(nameWithExt, ".sql")
		if !strings.HasSuffix(nameWithExt, ".sql") || name == "" {
			continue
		}

		numericPart := fileSlice[timestampIdx]
		result, err := strconv.Atoi(numericPart)
		if err != nil {
			continue
		}

		migrations[index] = migrationFile{
			Name:      name,
			Timestamp: result,
			Path:      path,
			IsNewFile: false,
		}
	}

	return migrations
}
