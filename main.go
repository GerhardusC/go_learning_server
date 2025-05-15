package main

import (
	"testing-server/handlers"
	"testing-server/dataCollection"
	"testing-server/cliArgs"
	"testing-server/dbInteractions"
	"testing-server/utils"
	_ "github.com/mattn/go-sqlite3"
)


func main () {
	cliargs.InitArgs()
	dbInteractions.InitDB()

	go utils.SendExampleEmail()

	go datacollection.CollectData()
	handlers.InitHandlers()
}
