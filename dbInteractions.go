package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

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
		return nil, fmt.Errorf("Failed to open database.")
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

func readSinceMeasurementsFromDB (timestamp int) ([]Measurement, error) {
	db, openErr := sql.Open("sqlite3", "./data.db")
	defer db.Close()
	
	if openErr != nil {
		return nil, openErr
	}

	measurements := make([]Measurement, 0, 500)

	res, queryErr := db.Query("SELECT * FROM READINGS WHERE timestamp > ?;", timestamp)

	if queryErr != nil {
		return nil, queryErr
	}

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
