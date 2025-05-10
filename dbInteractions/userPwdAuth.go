package dbInteractions

import (
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/golang-jwt/jwt/v5"
)

func (user UserPreAuth) AuthenticateUsernamePwd () (string, error) {
	authorisedUser, err := user.GetFromDB()
	
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "test-secret"
	}

	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_details": authorisedUser,
			"exp": now.Add(time.Duration(2*time.Hour)),
		},
	)

	if err != nil {
		return "", err
	}
	
	return token.SignedString([]byte(sec))
}
