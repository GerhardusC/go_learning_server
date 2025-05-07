package dbInteractions

import (
	"fmt"
	"database/sql"
	"log"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)


func WriteMeasurementToDB (measurement MeasurementInterface) error {
	err := measurement.performSingleWrite("MEASUREMENTS")
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong while writing measurement to the database")
	}
	return nil
}

func WriteLogEntryToDB (measurement MeasurementInterface) error {
	err := measurement.performSingleWrite("LOGS")
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong while writing log entry to the database")
	}
	return nil
}

type MeasurementInterface interface {
	performSingleWrite(tableName string) error
}

func (measurement DBRowMeasurement) performSingleWrite (tableName string) error {
	db, err := sql.Open("sqlite3", cliargs.DbPath)		
	defer db.Close()

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to open database.")
	}

	_, err = db.Exec(
		fmt.Sprintf(`INSERT INTO %s (timestamp, topic, value)
		VALUES (?, ?, ?)
		`, tableName),
		measurement.Timestamp,
		measurement.Topic,
		measurement.Value,
	)

	if err != nil {
		return err
	}
	return nil
}

