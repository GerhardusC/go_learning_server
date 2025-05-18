package dbInteractions

import (
	"database/sql"
	"crypto/sha256"
	"fmt"
	"log"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)

func (user UserPreAuth) SaveToDb(permissionLevel int) error {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to create user")
		return err
	}
	
	hashedPwd := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.UnhashedPwd)))

	query := `
		INSERT INTO USERS	(created_at, email,	username,	hashed_pwd,	permission_level)
		VALUES			(datetime(), ?,		?,		?,		?		);
	`

	_, err = db.Exec(query, user.Email, user.Username, hashedPwd, permissionLevel)

	if err != nil {
		return err
	}
	return nil
}

func (user UserWithHashedPwd) SaveToDb(permissionLevel int) error {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to create user")
		return err
	}

	query := `
		INSERT INTO USERS	(created_at, email,	username,	hashed_pwd,	permission_level)
		VALUES			(datetime(), ?,		?,		?,		?		);
	`

	_, err = db.Exec(query, user.Email, user.Username, user.HashedPwd, permissionLevel)

	if err != nil {
		return err
	}
	return nil
	
}
