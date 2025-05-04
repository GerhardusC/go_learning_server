package main

import (
	"fmt"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)


func main () {
	http.HandleFunc("/person", peopleHandler)
	http.HandleFunc("/measurements", allMeasurementsHandler)
	http.HandleFunc("/measurements_since", sinceMeasurementsHandler)

	http.ListenAndServe(":8080", nil)
	fmt.Printf("Listening on port 8080")
}
