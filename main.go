package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const SERVERPORT = ":5000"

func main() {
	ratesFile, err := os.ReadFile("testing.json")

	if err != nil {
		panic("Unable to load JSON file!")
	}
	decodedRates := JSONFileToJSONRates(ratesFile)
	AllDailyRates := AllJSONRatesToDailyRates(decodedRates)

	for _, dailyRate := range AllDailyRates {
		fmt.Println(dailyRate)
	}

	//Create HTTP server to handle "price" and "rates" endpoints
	http.HandleFunc("/rates", ratesURL) //GET and PUT
	http.HandleFunc("/price", priceURL) //GET
	log.Println("Server running on port" + SERVERPORT)
	log.Fatal(http.ListenAndServe(SERVERPORT, nil))
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


Comments:
1. Potentially refactor JSON marshal to encode
*/
