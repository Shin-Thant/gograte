package gograte

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var initUsage = `No migration directory found.
Please run "gograte init" to create a migrations directory.`

var initContent = `
-- +gograte Up
-- SQL in section 'Up' is executed when this migration is applied

-- +gograte Down
-- SQL section 'Down' is executed when this migration is rolled back
`

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

	id := uuid.New()
	fmt.Println(id)

	// check directory exists
	info, err := os.Stat("./migrations")
	if err != nil {
		log.Fatalf(`Error checking for migration directory: %s

%s`, err, initUsage)

	}
	if !info.IsDir() {
		log.Fatalln(initUsage)
	}

	file, err := os.Create(filepath.Join("./migrations", id.String()+".sql"))
	if err != nil {
		log.Fatalln("Error creating migration file:", err)
	}
	defer file.Close()
	file.Write([]byte(initContent))
}
