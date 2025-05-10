package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing-server/dbInteractions"
)

func signupHandler (writer http.ResponseWriter, request *http.Request) {
	var preAuthUser dbInteractions.UserPreAuth

	err := json.NewDecoder(request.Body).Decode(&preAuthUser)

	log.Println("Pre Auth User: ", preAuthUser)

	if err != nil {
		http.Error(writer, errors.New("Failed to decode JSON").Error(), http.StatusBadRequest)
		return
	}

	err = preAuthUser.SaveToDb(0)
	
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := preAuthUser.AuthenticateUsernamePwd()

	if err != nil {
		http.Error(writer, errors.New("Unauthorised").Error(), http.StatusUnauthorized)
		return
	}
	
	writer.Header().Set("Authorization", fmt.Sprint("Bearer ", token))
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

	token, err := preAuthUser.AuthenticateUsernamePwd()
	
	if err != nil {
		http.Error(writer, errors.New("Unauthorised").Error(), http.StatusUnauthorized)
		return
	}

	log.Println("TokenString: ", token)
	
	writer.Header().Set("Authorization", fmt.Sprint("Bearer ", token))
	writer.Header().Set("content-type", "text/plain")
	writer.Write([]byte("Login successful."))
}
