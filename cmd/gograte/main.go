package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Shin-Thant/gograte"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var usageMessage = `Usage:

// Create
gograte create [migration_name]

// Migrate
gograte migrate [db_driver] [db_url] [migrate_action]`

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println(usageMessage)
		os.Exit(1)
	}

	action := args[0]

	switch action {
	case "create":
		fmt.Println("Creating migration file...")
	case "migrate":
		gograte.Migrate(args)
	default:
		fmt.Println(usageMessage)
	}
}
