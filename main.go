package main

import (
	//	"btyes"
	"encoding/json"
	"fmt"
	//"os"
)

/* Brain Storm
use gofmt for formatting
use defer to close files after opening

*/
// module github.com/Freddy-S3/PriceAPI
//const serverPort = 5000

type Rate struct {
	Days  string `json:"days"`  //datetime
	Times string `json:"times"` //datetime
	Tz    string `json:"tz"`    //datetime?
	Price int    `json:"price"`
}

func main() {
	ratesList := []Rate{
		{
			Days:  "mon,tues,thurs",
			Times: "0900-2100",
			Tz:    "America/Chicago",
			Price: 1500,
		},
		{
			Days:  "fri,sat,sun",
			Times: "0900-2100",
			Tz:    "America/Chicago",
			Price: 2000,
		},
	}

	tester, _ := json.Marshal(ratesList)
	fmt.Printf("%s\n", tester)

	var decodedRates []Rate
	json.Unmarshal(tester, &decodedRates)
	fmt.Printf("%v", decodedRates)

}
