package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"testing-server/dbInteractions"
	"testing-server/utils"
)

func signupHandler (writer http.ResponseWriter, request *http.Request) {
	var preAuthUser dbInteractions.UserPreAuth

	err := json.NewDecoder(request.Body).Decode(&preAuthUser)

	log.Println("Pre Auth User: ", preAuthUser)

	if err != nil {
		http.Error(writer, errors.New("Failed to decode JSON").Error(), http.StatusBadRequest)
		return
	}

	err = utils.ValidateEmailPwd(preAuthUser.Email, preAuthUser.UnhashedPwd)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = preAuthUser.CheckUsernameAndEmailAvailability()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}


	sessionID, err := preAuthUser.SendOTP()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "application/json")
	writer.Write([]byte(`{"session_id": "` + sessionID + `"}`))
}

func loginHandler (writer http.ResponseWriter, request *http.Request) {
	// TODO: One day, add 2 steps to login.
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
	
	writer.Header().Set("content-type", "application/json")
	writer.Write([]byte(`{"token": "` + token + `"}`))
}

func verifyOTPSignupHandler (writer http.ResponseWriter, request *http.Request) {
	var verifyOTPObj dbInteractions.OTPVerifyObj

	err := json.NewDecoder(request.Body).Decode(&verifyOTPObj)

	log.Println("Verify OTP Obj:", verifyOTPObj)

	if err != nil {
		http.Error(writer, errors.New("Failed to decode JSON").Error(), http.StatusBadRequest)
		return
	}

	user, err := verifyOTPObj.GetUser()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	err = user.SaveToDb(0)
	
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	authorisedUser, err := dbInteractions.GetUserByUsername(user.Username)

	token, err := authorisedUser.GenerateToken()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
	
	writer.Header().Set("content-type", "application/json")
	writer.Write([]byte(`{"token": "` + token + `"}`))
}
