package dbInteractions

import (
	"database/sql"
	"testing-server/types" 
	"testing-server/cliArgs" 
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)



func scanRowsIntoArray (measurements *[]types.Measurement, res *sql.Rows) error {
	for res.Next() {
		var measurement DBRowMeasurement
		scanErr := res.Scan(&measurement.Timestamp, &measurement.Topic, &measurement.Value)

		if scanErr != nil {
			return scanErr
		}
		
		*measurements = append(
			*measurements, 
			types.Measurement{
				Timestamp: measurement.Timestamp,
				Topic: measurement.Topic,
				Value: measurement.Value,
			},
		)
	}
	return nil
}

func ReadAllMeasurementsFromDB () ([]types.Measurement, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)		
	defer db.Close()
	
	if(err != nil){
		log.Println(err)
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
	db, err := sql.Open("sqlite3", cliargs.DbPath)
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
	db, err := sql.Open("sqlite3", cliargs.DbPath)
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
