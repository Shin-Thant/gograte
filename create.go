package gograte

import (
	"fmt"
	"os"

	"github.com/google/uuid"
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

	id := uuid.New()
	fmt.Println(id)
}
