package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing-server/cliArgs"
	"testing-server/dbInteractions"
	"testing-server/middleware"
	"testing-server/types"
)


func InitHandlers () {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /person", middleware.CheckAuth(peopleHandler))

	// Measurements
	mux.HandleFunc("GET /measurements", allMeasurementsHandler)
	mux.HandleFunc("GET /measurements/since", sinceMeasurementsHandler)
	mux.HandleFunc("GET /measurements/between", betweenMeasurementsHandler)

	// User management
	mux.HandleFunc("POST /signup", signupHandler)
	mux.HandleFunc("POST /login", loginHandler)

	fs := http.FileServer(http.Dir(cliargs.ServeDir))
	mux.Handle("GET /", fs)


	muxMiddlewareApplied := middleware.NewLogger(mux)

	server := http.Server{
		Handler: muxMiddlewareApplied,
	}
	log.Println("Serving on port 80")
	err := server.ListenAndServe()
	log.Println("Failed to serve on port 80")

	if err != nil {
		server.Addr = ":8080"
		log.Println("Falling back to 8080")
		server.ListenAndServe()
	}
}


func peopleHandler (writer http.ResponseWriter, request *http.Request) {
	names := []string{
		"Adam",
		"Sally",
		"Steven",
		"Sarah",
	}

	surnames := []string{
		"Smith",
		"Smith",
		"Phtevens",
		"Smith",
	}

	var length int

	surnameslen := len(surnames)
	nameslen := len(names)

	length = min(surnameslen, nameslen)

	for i:=range length {
		names[i] = names[i] + " " + surnames[i]
	}

	person := types.Person {
		Name: "John",
		Surname: "Smith",
		Siblings: names,
	}

	user := request.Context().Value(middleware.AuthUserKey).(dbInteractions.User)
	writer.Header().Set("Content-Type", "application/json")

	writer.Header().Add("User-Details", fmt.Sprint(user.Username))
	// Write output.
	json.NewEncoder(writer).Encode(person)
}
