package dbInteractions

import (
	"database/sql"
	"fmt"
	"log"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)

type Measurement struct {
	Timestamp int		`json:"timestamp"`
	Topic string            `json:"topic"`
	Value float64           `json:"value"`
}

func scanMeasurementsIntoArray (measurements *[]Measurement, res *sql.Rows) error {
	for res.Next() {
		var measurement DBRowMeasurement[float64]
		scanErr := res.Scan(&measurement.Timestamp, &measurement.Topic, &measurement.Value)

		if scanErr != nil {
			return scanErr
		}
		
		*measurements = append(
			*measurements, 
			Measurement{
				Timestamp: measurement.Timestamp,
				Topic: measurement.Topic,
				Value: measurement.Value,
			},
		)
	}
	return nil
}

func ReadAllMeasurementsFromDB () ([]Measurement, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)		
	defer db.Close()
	
	if(err != nil){
		log.Println(err)
		return nil, fmt.Errorf("Failed to open database.")
	}
	
	res, err := db.Query("SELECT * FROM MEASUREMENTS;")

	if(err != nil){
		return nil, err
	}

	measurements := make([]Measurement, 0, 500)

	scanMeasurementsIntoArray(&measurements, res)	

	return measurements, nil
}

func ReadBetweenMeasurementsFromDB (start int, stop int) ([]Measurement, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		return nil, err
	}
	
	measurements := make([]Measurement, 0, 500)

	res, err := db.Query("SELECT * FROM MEASUREMENTS WHERE timestamp > ? AND timestamp < ?;", start, stop)

	if err != nil {
		return nil, err
	}

	scanMeasurementsIntoArray(&measurements, res)

	return measurements, nil
}

func ReadSinceMeasurementsFromDB (timestamp int) ([]Measurement, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()
	
	if err != nil {
		return nil, err
	}

	measurements := make([]Measurement, 0, 500)

	res, err := db.Query("SELECT * FROM MEASUREMENTS WHERE timestamp > ?;", timestamp)

	if err != nil {
		return nil, err
	}

	scanMeasurementsIntoArray(&measurements, res)

	return measurements, nil
}
