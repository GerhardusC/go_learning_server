package dbInteractions

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing-server/cliArgs"

	_ "github.com/mattn/go-sqlite3"
)

func (user *UserPreAuth) CheckUsernameAndEmailAvailability () error {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to retrieve user")
		return err
	}

	query := `
		SELECT email, username FROM USERS
		WHERE username = ? OR email = ?;
	`
	row := db.QueryRow(query, user.Username, user.Email)

	var retrievedEmail string
	var retrievedUsername string

	err = row.Scan(&retrievedEmail, &retrievedUsername)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		return err
	}

	if retrievedEmail != "" {
		return errors.New("Email not available")
	}

	if retrievedUsername != "" {
		return errors.New("Username not available")
	}

	return nil
}

/** 
** A shortcut function to get user. Only use if the user is already authenticated.*/
func GetUserByUsername (username string) (User, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to retrieve user")
		return User{}, err
	}

	query := `
		SELECT id, created_at, email, username, permission_level FROM USERS
		WHERE username = ?;
	`
	row := db.QueryRow(query, username)

	var retrievedUser User

	err = row.Scan(
		&retrievedUser.ID,
		&retrievedUser.CreatedAt,
		&retrievedUser.Email,
		&retrievedUser.Username,
		&retrievedUser.PermissionLevel,
	)

	if err != nil {
		return User{}, err
	}

	return retrievedUser, nil
}

func (user *UserWithHashedPwd) GetFromDB() (User, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to retrieve user")
		return User{}, err
	}

	query := `
		SELECT id, created_at, email, username, permission_level FROM USERS
		WHERE username = ?;
	`
	row := db.QueryRow(query, user.Username)

	var retrievedUser User

	err = row.Scan(
		&retrievedUser.ID,
		&retrievedUser.CreatedAt,
		&retrievedUser.Email,
		&retrievedUser.Username,
		&retrievedUser.PermissionLevel,
	)

	if err != nil {
		fmt.Println("Failed to get user from DB")
		return User{}, err
	}

	return retrievedUser, nil
}

func (user *UserPreAuth) GetFromDB() (User, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to retrieve user")
		return User{}, err
	}

	query := `
		SELECT id, created_at, email, username, permission_level, hashed_pwd FROM USERS
		WHERE username = ?;
	`
	row := db.QueryRow(query, user.Username)

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

	receivedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.UnhashedPwd)))

	log.Println("StoredHash: ", hashedPwd, "\nReceivedHash: ", receivedHash)

	if  subtle.ConstantTimeCompare([]byte(receivedHash), []byte(hashedPwd)) == 1 {
		return retrievedUser, nil
	}

	return User{}, errors.New("User not authenticated.")
}
