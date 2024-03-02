package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Shin-Thant/gograte"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	help := flag.Bool("help", false, "Show usage message")
	flag.Parse()

	if *help {
		fmt.Println(gograte.UsageMessage)
		return
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println(gograte.UsageMessage)
		os.Exit(1)
	}

	action := args[0]

	switch action {
	case "init":
		gograte.Init()
	case "create":
		gograte.Create(args)
	case "status":
		gograte.Status(args)
	case "migrate":
		gograte.Migrate(args)
	default:
		fmt.Println(gograte.UsageMessage)
	}
}
