package dbInteractions

import (
	"database/sql"
	"testing-server/types" 
	"flag"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type DBRowMeasurement struct {
	timestamp int
	topic string
	value float64
}

var dbPath string

func InitDBPathFromArgs() {
	flag.StringVar(&dbPath, "d", "./data.db", "Database path")
	flag.Parse()
	log.Println("DB Path: ", dbPath)
}

func scanRowsIntoArray (measurements *[]types.Measurement, res *sql.Rows) error {
	for res.Next() {
		var measurement DBRowMeasurement
		scanErr := res.Scan(&measurement.timestamp, &measurement.topic, &measurement.value)

		if scanErr != nil {
			return scanErr
		}
		
		*measurements = append(
			*measurements, 
			types.Measurement{
				Timestamp: measurement.timestamp,
				Topic: measurement.topic,
				Value: measurement.value,
			},
		)
	}
	return nil
}

func ReadAllMeasurementsFromDB () ([]types.Measurement, error) {
	db, err := sql.Open("sqlite3", dbPath)		
	defer db.Close()
	
	if(err != nil){
		log.Fatal(err)
		return nil, fmt.Errorf("Failed to open database.")
	}
	
	res, err := db.Query("SELECT * FROM READINGS;")

	if(err != nil){
		return nil, err
	}

	measurements := make([]types.Measurement, 0, 500)

	scanRowsIntoArray(&measurements, res)	

	return measurements, nil
}

func ReadBetweenMeasurementsFromDB (start int, stop int) ([]types.Measurement, error) {
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		return nil, err
	}
	
	measurements := make([]types.Measurement, 0, 500)

	res, err := db.Query("SELECT * FROM READINGS WHERE timestamp > ? AND timestamp < ?;", start, stop)

	if err != nil {
		return nil, err
	}

	scanRowsIntoArray(&measurements, res)

	return measurements, nil
}

func ReadSinceMeasurementsFromDB (timestamp int) ([]types.Measurement, error) {
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	
	if err != nil {
		return nil, err
	}

	measurements := make([]types.Measurement, 0, 500)

	res, err := db.Query("SELECT * FROM READINGS WHERE timestamp > ?;", timestamp)

	if err != nil {
		return nil, err
	}

	scanRowsIntoArray(&measurements, res)

	return measurements, nil
}
