package dbInteractions

import (
	"database/sql"
	"errors"
	"log"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB () error {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		return errors.New("Could not open initial connection to DB")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS MEASUREMENTS (
				timestamp int,
				topic varchar(255),
				value float
			)
		`)

	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS LOGS (
				timestamp int,
				topic varchar(255),
				value varchar(255)
			)
		`)

	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USERS (
				id INTEGER PRIMARY KEY ASC AUTOINCREMENT NOT NULL,
				created_at date NOT NULL,
				email TEXT NOT NULL UNIQUE,
				username TEXT NOT NULL UNIQUE,
				hashed_pwd TEXT,
				permission_level int
			)
		`)

	if err != nil {
		log.Println(err)
	}

	return nil
}
