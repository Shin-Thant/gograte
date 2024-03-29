package gograte

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// argument indexes
var MIGRATION = 1

func Create(args []string) {
	if len(args) != 2 {
		fmt.Println(UsageMessage)
		os.Exit(1)
	}

	migrationName := args[MIGRATION]
	if migrationName == "" {
		fmt.Println(UsageMessage)
		os.Exit(1)
	}
	fmt.Println("Creating migration file:", migrationName)

	// check directory exists
	info, err := os.Stat("./migrations")
	if err != nil {
		log.Fatalf(`Error checking for migration directory: %s

%s`, err, initUsage)

	}
	if !info.IsDir() {
		log.Fatalln(initUsage)
	}

	filename := generateFilename(migrationName)
	file, err := os.Create(filepath.Join("./migrations", filename))
	if err != nil {
		log.Fatalln("Error creating migration file:", err)
	}
	defer file.Close()
	file.Write([]byte(initContent))

	log.Println("Migration file created:", filename)
}

func generateFilename(migrationName string) string {
	return fmt.Sprintf("%v_%s.sql", time.Now().UTC().Format(timestampFormat), migrationName)
}
