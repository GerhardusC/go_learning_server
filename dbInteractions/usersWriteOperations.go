package dbInteractions

import (
	"database/sql"
	"log"
	"testing-server/cliArgs"
	"testing-server/utils"

	_ "github.com/mattn/go-sqlite3"
)

func (user UserPreAuth) SaveToDb(permissionLevel int) (User, error) {
	db, err := sql.Open("sqlite3", cliargs.DbPath)
	defer db.Close()

	if err != nil {
		log.Println("Unable to open database to create user")
		return User{}, err
	}
	
	//TODO! Ensure user is not spammer, do some sort of 2 step email or something.
	
	hashedPwd := pwdauth.SaltAndHashPwd(user.Username, user.UnhashedPwd)

	query := `
		INSERT INTO USERS	(created_at, email,	username,	hashed_pwd,	permission_level)
		VALUES			(datetime(), ?,		?,		?,		?		);
	`

	_, err = db.Exec(query, user.Email, user.Username, hashedPwd, permissionLevel)

	if err != nil {
		return User{}, err
	}

	createdUser, err := user.GetFromDB()

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return createdUser, nil
}


