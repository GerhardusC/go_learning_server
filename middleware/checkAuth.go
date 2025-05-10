package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"testing-server/dbInteractions"

	"github.com/golang-jwt/jwt/v5"
)

type userContextKey string

const AuthUserKey userContextKey = "user context key"

type userJwtClaims struct {
	UserDetails dbInteractions.User `json:"user_details"`
	jwt.RegisteredClaims
}

func CheckAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func (writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")

		splitAuthHeader := strings.Split(authHeader, " ")

		if len(splitAuthHeader) < 2 {
			http.Error(writer, errors.New("Authorization header malformed").Error(), http.StatusUnauthorized)
			return
		}

		tokenString := splitAuthHeader[1]

		sec := os.Getenv("JWT_SECRET")
		if sec == "" {
			sec = "test-secret"
		}

		token, err := jwt.ParseWithClaims(tokenString, &userJwtClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(sec), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			log.Println("Error in check auth: ", err)
			http.Error(writer, errors.New("Token invalid").Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*userJwtClaims)

		if !ok {
			http.Error(writer, errors.New("Token invalid").Error(), http.StatusUnauthorized)
			return
		}

		log.Println("User: ", claims.UserDetails)

		ctx := context.WithValue(request.Context(), AuthUserKey, claims.UserDetails)

		reqWithUser := request.WithContext(ctx)

		handler(writer, reqWithUser)
	}
}

