package main

import (
	"os"
	"testing"
)

// Loads the initial seeded data into the priceDB
func TestLoadInitialSeededDataToPriceDB(t *testing.T) {
	// load first, then test if loaded DB is same as initial
	loadInitialSeededDataToPriceDB()

	testResult, err := os.ReadFile("priceDB.json")
	if err != nil {
		t.Errorf("Unable to load priceDB.json")
	}

	initialDB, err := os.ReadFile("initialSeededData.json")
	if err != nil {
		t.Errorf("Unable to load initialSeededData.json")
	}

	for index, resultBytes := range testResult {
		if resultBytes != initialDB[index] {
			t.Errorf("loadInitialSeededDataToPriceDB()= %v; want %v", testResult, initialDB)
		}
	}
}
