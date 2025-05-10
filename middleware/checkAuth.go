package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"testing-server/dbInteractions"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userContextKey string

const AuthUserKey userContextKey = "user context key"

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

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(sec), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			http.Error(writer, errors.New("Token invalid").Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(writer, errors.New("Token invalid").Error(), http.StatusUnauthorized)
			return
		}

		decodedUser := claims["user_details"].(dbInteractions.User)
		exp := claims["exp"].(time.Time)

		log.Println("expiery: ", exp)
		log.Println("New user: ", exp)

		ctx := context.WithValue(request.Context(), AuthUserKey, decodedUser)

		reqWithUser := request.WithContext(ctx)

		handler(writer, reqWithUser)
	}
}

