package common

import (
	"encoding/json"
	"io"
)

type JSONRate struct {
	Days  string `json:"days"`
	Times string `json:"times"`
	Tz    string `json:"tz"`
	Price int    `json:"price"`
}

type AllJSONRates struct {
	Rates []JSONRate `json:"rates"`
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
