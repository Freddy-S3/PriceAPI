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
	//load the Database from the initial seeded values
	LoadInitialSeededDataToPriceDB()

	//Create HTTP server to handle "price" and "rates" endpoints
	http.HandleFunc("/rates", ratesURL) //GET and PUT
	http.HandleFunc("/price", priceURL) //GET
	log.Println("Server running on port" + SERVERPORT)
	log.Fatal(http.ListenAndServe(SERVERPORT, nil))
}

// loads original seeded data to the priceDB.json
func LoadInitialSeededDataToPriceDB() {
	ratesFile, err := os.ReadFile("initialSeededData.json")
	if err != nil {
		panic("Unable to load initialSeededJSON file!")
	}
	os.WriteFile("priceDB.json", ratesFile, os.ModePerm)
}

// returns "unavailable" as a JSON response
func httpResponseOfUnavailable(w http.ResponseWriter) http.ResponseWriter {
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
1. input spanning more than 1 day in convertedTime or uncoverted time?
	a. unconverted time spanning 24 hours or calendar day?
2. can ignore seconds? or inclusive up till 59 seconds?

TODO:
1. refactor rate into interface?
	a. https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response
	b. can return price field more easily with map[string]interface{}
2. account for seconds in query
3. HTTP Tests
4. Packaging (refactor into folders? gomod?)
5. README - Include any instructions on how to build, run, and test your application
*/
