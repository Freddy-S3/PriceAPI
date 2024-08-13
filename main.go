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
