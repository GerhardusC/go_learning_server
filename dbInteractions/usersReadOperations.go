package dbInteractions

import (
	"database/sql"
	"errors"
	"log"
	"testing-server/cliArgs"
	"testing-server/utils"

	_ "github.com/mattn/go-sqlite3"
)

func (user UserPreAuth) GetFromDB() (User, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to retrieve user")
		return User{}, err
	}

	query := `
		SELECT id, created_at, email, username, permission_level, hashed_pwd FROM USERS
		WHERE email = ? AND username = ?;
	`
	row := db.QueryRow(query, user.Email, user.Username)

	var hashedPwd string
	var retrievedUser User

	err = row.Scan(
		&retrievedUser.ID,
		&retrievedUser.CreatedAt,
		&retrievedUser.Email,
		&retrievedUser.Username,
		&retrievedUser.PermissionLevel,
		&hashedPwd,
	)

	if err != nil {
		return User{}, err
	}

	log.Println("Pre-user: ", user, "\nRetrieved User", retrievedUser, "\nHashed Pwd", hashedPwd)

	if hashedPwd == "" || user.UnhashedPwd == "" {
		return User{}, errors.New("Password authentication not enabled.")
	}

	receivedHash := pwdauth.SaltAndHashPwd(user.Username, user.UnhashedPwd)

	log.Println("StoredHash: ", hashedPwd, "\nReceivedHash: ", receivedHash)

	if  receivedHash == hashedPwd {
		return retrievedUser, nil
	}

	return User{}, errors.New("User not authenticated.")
}
