package main

import (
	"os"
	"testing-server/cliArgs"
	"testing-server/dataCollection"
	"testing-server/dbInteractions"
	"testing-server/handlers"
	"testing-server/utils"

	_ "github.com/mattn/go-sqlite3"
)


func main () {
	cliargs.InitArgs()
	dbInteractions.InitDB()

	if os.Getenv("SEND_WELCOME_EMAIL") == "YES" {
		go utils.SendExampleEmail()
	}

	go datacollection.CollectData()
	handlers.InitHandlers()
}
