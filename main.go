package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const SERVERPORT = ":5000"

func main() {
	// Load the Database from the initial seeded values
	loadInitialSeededDataToPriceDB()

	// Create HTTP server to handle "price" and "rates" endpoints
	http.HandleFunc("/rates", ratesURL) //GET and PUT
	http.HandleFunc("/price", priceURL) //GET
	log.Println("Server running on port" + SERVERPORT)
	log.Fatal(http.ListenAndServe(SERVERPORT, nil))
}

// Loads original seeded data to the priceDB.json
func loadInitialSeededDataToPriceDB() {
	ratesFile, err := os.ReadFile("initialSeededData.json")
	if err != nil {
		panic("Unable to load initialSeededJSON file!")
	}
	os.WriteFile("priceDB.json", ratesFile, os.ModePerm)
}

// Returns "unavailable" as a JSON response
func HttpResponseOfUnavailable(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("unavailable")
	return w
}

// Decodes JSON file into JSONRates struct
func JSONFileToJSONRates(ratesFile io.Reader) (AllJSONRates, error) {
	var decodedJSONRates AllJSONRates
	err := json.NewDecoder(ratesFile).Decode(&decodedJSONRates)
	if err != nil {
		return AllJSONRates{}, err
	}

	return decodedJSONRates, nil
}

/*
Questions:
1. Input spanning more than 1 day in convertedTime or uncoverted time?
	a. Unconverted time spanning 24 hours or calendar day?
	b. "Can span multiple rates" meaning rate overlap? or time band wide enough to be in 2 bands
2. Can ignore seconds? or inclusive up till 59 seconds?

TODO:
1. refactor rate into interface?
	a. https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response
	b. can return price field more easily with map[string]interface{}
2. account for seconds in query
3. HTTP Tests
4. Packaging (refactor into folders? gomod?)
5. README - Include any instructions on how to build, run, and test your application
6. Test input for overlap?
7. Potentially add channels?
*/
