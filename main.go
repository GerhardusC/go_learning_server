package main

import (
	"fmt"
	"log"
	"strconv"
	"database/sql"
	"net/http"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	Name string		`json:"name"`
	Surname string		`json:"surname"`
	Siblings []string	`json:"siblings"`
}

type Measurement struct {
	Timestamp int		`json:"timestamp"`
	Topic string            `json:"topic"`
	Value float64           `json:"value"`
}

type DBRowMeasurement struct {
	timestamp int
	topic string
	value float64
}

func readAllMeasurementsFromDB () ([]Measurement, error) {
	db, openErr := sql.Open("sqlite3", "./data.db")		
	defer db.Close()
	
	if(openErr != nil){
		log.Fatal(openErr)
	}
	
	res, queryErr := db.Query("SELECT * FROM READINGS;")

	if(queryErr != nil){
		return nil, queryErr
	}

	measurements := make([]Measurement, 0, 500)

	for res.Next() {
		var measurement DBRowMeasurement
		scanErr := res.Scan(&measurement.timestamp, &measurement.topic, &measurement.value)

		if scanErr != nil {
			return nil, scanErr
		}

		measurements = append(
			measurements,
			Measurement{
				Timestamp: measurement.timestamp,
				Topic: measurement.topic,
				Value: measurement.value,
			},
		)
	}
	return measurements, nil	
}

func allMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	measurements, err := readAllMeasurementsFromDB();

	if err != nil {
		http.Error(writer, err.Error(), 404)
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func readSinceMeasurementsFromDB (timestamp int) ([]Measurement, error) {
	db, openErr := sql.Open("sqlite3", "./data.db")		
	defer db.Close()
	
	if openErr != nil {
		return nil, openErr
	}

	measurements := make([]Measurement, 0, 500)

	res, queryErr := db.Query("SELECT * FROM READINGS WHERE timestamp < ?;", timestamp)

	if queryErr != nil {
		return nil, queryErr
	}

	m

	for res.Next() {
		var measurement DBRowMeasurement
		res.Scan()
		
	}

}

func sinceMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	sinceTimestampStr := request.URL.Query().Get("timestamp")
	sinceTimestamp, err := strconv.ParseInt(sinceTimestampStr, 0, 64)	

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}


	
}

func peopleHandler (writer http.ResponseWriter, request *http.Request) {
	// Set headers.
	writer.Header().Set("Content-Type", "application/json")

	names := []string{
		"Adam",
		"Sally",
		"Steven",
		"Sarah",
	}

	surnames := []string{
		"Smith",
		"Smith",
		"Phtevens",
		"Smith",
	}

	var length int

	surnameslen := len(surnames)
	nameslen := len(names)

	length = min(surnameslen, nameslen)

	for i:=range length {
		names[i] = names[i] + " " + surnames[i]
	}

	person := Person {
		Name: "John",
		Surname: "Smith",
		Siblings: names,
	}
	// Write output.
	json.NewEncoder(writer).Encode(person)
}

func main () {
	http.HandleFunc("/person", peopleHandler)
	http.HandleFunc("/measurements", allMeasurementsHandler)
	http.HandleFunc("/measurements_since", sinceMeasurementsHandler)

	http.ListenAndServe(":8080", nil)
	fmt.Printf("Listening on port 8080")
}
