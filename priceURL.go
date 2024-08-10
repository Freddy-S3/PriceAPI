package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// GET /price: return JSON for price given time query
// ex query: http://localhost:5000/price?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00
func priceURL(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	// parse start URL
	startTime, err := time.Parse(URLDATETIMELAYOUT, startStr)
	if err != nil {
		http.Error(w, "Invalid 'start' time format", http.StatusBadRequest)
		return
	}

	// parse end URL #potential refactor
	endTime, err := time.Parse(URLDATETIMELAYOUT, endStr)
	if err != nil {
		http.Error(w, "Invalid 'end' time format", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(os.Stdout, []any{startTime, endTime}...)

	// check if input wrong
	if endTime.After(startTime) { //check test case if equal or earlier
		//if input spans more than a day
		if startTime.Day() != endTime.Day() {
			//w.Header().Set("Content-Type", "application/json") //check what this is
			json.NewEncoder(w).Encode("unavailable") //return as JSON?
			return
		}
	} else {
		http.Error(w, "'end' before 'start'", http.StatusBadRequest)
		return
	}

	//input can't span more 2 time bands //assuming DB time zone?
	//look through the time bands
}
