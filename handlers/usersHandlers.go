package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing-server/dbInteractions"

	"github.com/golang-jwt/jwt/v5"
)



func signupHandler (writer http.ResponseWriter, request *http.Request) {
	var preAuthUser dbInteractions.UserPreAuth

	err := json.NewDecoder(request.Body).Decode(&preAuthUser)

	log.Println("Pre Auth User: ", preAuthUser)

	if err != nil {
		http.Error(writer, errors.New("Failed to decode JSON").Error(), http.StatusBadRequest)
		return
	}

	authorisedUser, err := preAuthUser.SaveToDb(0)
	
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "test-secret"
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"created_at": authorisedUser.CreatedAt,
			"email": authorisedUser.Email,
			"permission_level": authorisedUser.PermissionLevel,
			"username": authorisedUser.Username,
		},
	)

	tokenString, err := token.SignedString(sec)
	
	writer.Header().Set("Authorization", fmt.Sprint("Bearer ", tokenString))
	writer.Header().Set("content-type", "text/plain")
	writer.Write([]byte(""))
}

func login (writer http.ResponseWriter, request *http.Request) {
	
}
