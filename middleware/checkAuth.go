package middleware

import (
	"context"
	"net/http"
	"strings"
)

type userContextKey string

type User struct {
	Username string
}

const AuthUserKey userContextKey = "user context key"

func CheckAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func (writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("User-Agent"), "Mozilla") {
			writer.Write([]byte("Unauthorized"))
			return
		}

		newUser := User{Username: "John Doe"}

		ctx := context.WithValue(request.Context(), AuthUserKey, newUser)

		reqWithUser := request.WithContext(ctx)

		handler(writer, reqWithUser)
	}
}

