package main

import (
	"log"
	"net/http"
	"testing-server/handlers"
	"testing-server/dbInteractions"
	_ "github.com/mattn/go-sqlite3"
)


func main () {
	dbInteractions.InitDBPathFromArgs()
	handlers.InitHandlers()

	log.Println("Serving on port 80")
	err := http.ListenAndServe(":80", nil)
	log.Println("Failed to serve on port 80")

	if err != nil {
		log.Println("Falling back to 8080")
		http.ListenAndServe(":8080", nil)
	}
}
