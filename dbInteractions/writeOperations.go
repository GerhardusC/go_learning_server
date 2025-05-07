package dbInteractions

import (
	"fmt"
	"database/sql"
	"log"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)

func (measurement DBRowMeasurement) WriteToTable (tableName string) error {
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

