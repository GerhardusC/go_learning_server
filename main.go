package main

import (
	"testing-server/handlers"
	"testing-server/dataCollection"
	"testing-server/cliArgs"
	"testing-server/dbInteractions"
	_ "github.com/mattn/go-sqlite3"
)


func main () {
	cliargs.InitArgs()
	dbInteractions.InitDB()

	go datacollection.CollectData()
	handlers.InitHandlers()
}
