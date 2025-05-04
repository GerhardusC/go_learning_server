package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func allMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	measurements, err := readAllMeasurementsFromDB();

	if err != nil {
		http.Error(writer, err.Error(), 404)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func sinceMeasurementsHandler (writer http.ResponseWriter, request *http.Request) {
	sinceTimestampStr := request.URL.Query().Get("timestamp")
	sinceTimestamp, queryParamPareseErr := strconv.ParseInt(sinceTimestampStr, 0, 64)

	if queryParamPareseErr != nil {
		http.Error(writer, errors.New("Invalid query parameters. Ensure you have a timestamp as a query parameter.").Error(), http.StatusBadRequest)
		return;
	}

	measurements, dbReadErr := readSinceMeasurementsFromDB(int(sinceTimestamp))

	if dbReadErr != nil {
		http.Error(writer, errors.New("Failed to find database.").Error() , http.StatusInternalServerError)
		return;
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(measurements)
}

func peopleHandler (writer http.ResponseWriter, request *http.Request) {
	// Set headers.
	writer.Header().Set("Content-Type", "application/json")

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

	person := Person {
		Name: "John",
		Surname: "Smith",
		Siblings: names,
	}
	// Write output.
	json.NewEncoder(writer).Encode(person)
}
