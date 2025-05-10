package dbInteractions

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/golang-jwt/jwt/v5"
)

func (user UserPreAuth) AuthenticateUsernamePwd () (string, error) {
	authorisedUser, err := user.GetFromDB()
	
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "test-secret"
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_details": authorisedUser,
		},
	)

	if err != nil {
		return "", err
	}
	
	return token.SignedString([]byte(sec))
}
