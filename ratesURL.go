package main

import (
	"encoding/json"
	"net/http"
	"os"
)

// GET and PUT '/rates'
func ratesURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getAllRates(w, r)
	} else if r.Method == http.MethodPut {
		updateRate(w, r)
	}
}

// Test with: http://localhost:5000/rates
// GET Method: gets all rates in the JSON DB priceDB.json
func getAllRates(w http.ResponseWriter, r *http.Request) {
	ratesFile, err := os.ReadFile("priceDB.json")
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ratesFile)
}

// Test with: curl -X PUT -H "Content-Type: application/json" -d '{"key1":"value"}' http://localhost:5000/rates
// Used Postman for more simplified HTTP Put testing
// PUT Method: Replaces JSON DB with new input JSON DB //Will reload original DB upon server start
func updateRate(w http.ResponseWriter, r *http.Request) {
	// Decode HTTP Request
	allJSONRates, err := JSONFileToJSONRates(r.Body)
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}

	// Open DB file in write mode and truncate old contents
	ratesFile, err := os.OpenFile("priceDB.json", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}
	defer ratesFile.Close()

	// PUT new contents into file
	encoder := json.NewEncoder(ratesFile)
	encoder.SetIndent("", "  ") //optional for readability
	encoder.Encode(allJSONRates)

	// Test if results successful
	// getAllRates(w, r)
	// w.Write([]byte("JSON data saved successfully!"))
}
