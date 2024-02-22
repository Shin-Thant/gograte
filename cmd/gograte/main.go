package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Shin-Thant/gograte"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	flag.Parse()
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
	case "migrate":
		gograte.Migrate(args)
	default:
		fmt.Println(gograte.UsageMessage)
	}
}
