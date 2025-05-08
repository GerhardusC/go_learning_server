package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"testing-server/cliArgs"
	"testing-server/dbInteractions"
	"testing-server/types"
)


func InitHandlers () {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /person", peopleHandler)
	mux.HandleFunc("GET /measurements", allMeasurementsHandler)
	mux.HandleFunc("GET /measurements/since", sinceMeasurementsHandler)
	mux.HandleFunc("GET /measurements/between", betweenMeasurementsHandler)

	fs := http.FileServer(http.Dir(cliargs.ServeDir))
	mux.Handle("GET /", fs)

	server := http.Server{
		Handler: mux,
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

func allMeasurementsHandler (writer http.ResponseWriter, _ *http.Request) {
	measurements, err := dbInteractions.ReadAllMeasurementsFromDB();

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func betweenMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	startTimestampStr, stopTimestampStr :=
		request.URL.Query().Get("start") , request.URL.Query().Get("stop")

	startTimestamp, err := strconv.ParseInt(startTimestampStr, 0, 64)

	if err != nil {
		http.Error(writer, errors.New("Invalid start timestamp").Error() , http.StatusBadRequest)
		return
	}

	stopTimestamp, err := strconv.ParseInt(stopTimestampStr, 0, 64)

	if err != nil {
		http.Error(writer, errors.New("Invalid stop timestamp").Error() , http.StatusBadRequest)
		return
	}

	measurements, err := dbInteractions.ReadBetweenMeasurementsFromDB(int(startTimestamp), int(stopTimestamp))

	if err != nil {
		http.Error(
			writer,
			errors.New("Failed to read data from database").Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func sinceMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	sinceTimestampStr := request.URL.Query().Get("timestamp")
	sinceTimestamp, err := strconv.ParseInt(sinceTimestampStr, 0, 64)

	if err != nil {
		http.Error(writer, errors.New("Invalid query parameters. Ensure you have a timestamp as a query parameter.").Error(), http.StatusBadRequest)
		return;
	}

	measurements, err := dbInteractions.ReadSinceMeasurementsFromDB(int(sinceTimestamp))

	if err != nil {
		http.Error(writer, errors.New("Failed to find database.").Error() , http.StatusInternalServerError)
		return;
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func peopleHandler (writer http.ResponseWriter, _ *http.Request) {
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

	writer.Header().Set("Content-Type", "application/json")
	// Write output.
	json.NewEncoder(writer).Encode(person)
}
