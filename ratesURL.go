package main

import (
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
	ratesFile, err := os.ReadFile("testing.json")

	if err != nil {
		fmt.Println("Unable to load JSON file!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ratesFile)
}

// test w/: curl -X PUT -H "Content-Type: application/json" -d '{"key1":"value"}' http://localhost:5000/rates
// PUT Method:
func updateRate(w http.ResponseWriter, r *http.Request) {
	panic("ignore")
	//ratesFile, err := os.WriteFile("testing.json")
}
