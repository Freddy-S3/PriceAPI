package main

import (
	"encoding/json"
	"fmt"
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

// test w/: http://localhost:5000/rates
// GET Method: gets all rates in the JSON DB
func getAllRates(w http.ResponseWriter, r *http.Request) {
	ratesFile, err := os.ReadFile("priceDB.json")
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ratesFile)
}

// test w/: curl -X PUT -H "Content-Type: application/json" -d '{"key1":"value"}' http://localhost:5000/rates
// PUT Method: Replaces JSON DB with new input JSON DB
func updateRate(w http.ResponseWriter, r *http.Request) {
	allJSONRates, err := JSONFileToJSONRates(r.Body)
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}

	fmt.Println(allJSONRates)

	//remove contents of priceDB.json
	//os.WriteFile("priceDB.json", []byte(""), os.ModePerm)

	ratesFile, err := os.OpenFile("priceDB.json", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}
	defer ratesFile.Close()
	// PUT new contents
	json.NewEncoder(ratesFile).Encode(allJSONRates)

	//w.Write([]byte("JSON data saved successfully!"))
	getAllRates(w, r)
}
