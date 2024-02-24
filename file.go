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

func validateMigrationFilePaths(paths []string) []migrationFile {
	migrations := make([]migrationFile, len(paths))

	for index, path := range paths {
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

	return migrations
}
