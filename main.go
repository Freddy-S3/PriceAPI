package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Freddy-S3/PriceAPI/price"
	"github.com/Freddy-S3/PriceAPI/rates"
)

const SERVERPORT = ":5000"

func main() {
	// Load the Database from the initial seeded values
	loadInitialSeededDataToPriceDB()

	// Create HTTP server to handle "price" and "rates" endpoints
	http.HandleFunc("/rates", rates.RatesURL) //GET and PUT
	http.HandleFunc("/price", price.PriceURL) //GET
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
