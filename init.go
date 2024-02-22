package gograte

import (
	"fmt"
	"log"
	"os"
)

func Init() {
	info, err := os.Stat("./migrations")
	if err == nil && info.IsDir() {
		fmt.Println("Migration directory already exists.")
		return
	}

	if err := os.Mkdir("./migrations", 0755); err != nil {
		log.Fatalln("Error creating migrations directory:", err)
	}
}
