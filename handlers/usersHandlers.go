package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing-server/dbInteractions"
	"time"

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

	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_details": authorisedUser,
			"exp": now.Add(time.Duration(2*time.Hour)),
		},
	)

	tokenString, err := token.SignedString([]byte(sec))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("TokenString: ", tokenString)
	
	writer.Header().Set("Authorization", fmt.Sprint("Bearer ", tokenString))
	writer.Header().Set("content-type", "text/plain")
	writer.Write([]byte("Signup successful."))
}

func loginHandler (writer http.ResponseWriter, request *http.Request) {
	var preAuthUser dbInteractions.UserPreAuth

	err := json.NewDecoder(request.Body).Decode(&preAuthUser)

	if err != nil {
		http.Error(writer, errors.New("Failed to decode JSON").Error(), http.StatusBadRequest)
		return
	}

	authorisedUser, err := preAuthUser.GetFromDB()
	
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
		http.Error(writer, errors.New("Unauthorised").Error(), http.StatusUnauthorized)
		return
	}
	
	tokenString, err := token.SignedString([]byte(sec))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("TokenString: ", tokenString)
	
	writer.Header().Set("Authorization", fmt.Sprint("Bearer ", tokenString))
	writer.Header().Set("content-type", "text/plain")
	writer.Write([]byte("Login successful."))
}
