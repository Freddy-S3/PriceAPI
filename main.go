package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const SERVERPORT = ":5000"
const URLDATETIMELAYOUT string = "2006-01-02T15:04:05-07:00"

var JSONWEEKDAY = map[string]time.Weekday{
	"sun":   time.Sunday,
	"mon":   time.Monday,
	"tues":  time.Tuesday,
	"wed":   time.Wednesday,
	"thurs": time.Thursday,
	"fri":   time.Friday,
	"sat":   time.Saturday,
}

type JSONRate struct {
	Days  string `json:"days"`  //datetime
	Times string `json:"times"` //datetime
	Tz    string `json:"tz"`    //datetime?
	Price int    `json:"price"`
}

type AllJSONRates struct {
	Rates []JSONRate `json:"rates"`
}

type DailyRate struct {
	Weekdays  []time.Weekday // #potential refactor to be a single day
	StartTime time.Time
	EndTime   time.Time
	Price     int
}

type AllDailyRates struct {
	DailyRates []DailyRate
}

func main() {
	// read the JSON file
	ratesFile, err := os.ReadFile("testing.json")

	if err != nil {
		fmt.Println("Unable to load JSON file!")
		return
	}

	// Decode the JSON file
	var DecodedJSONRates AllJSONRates
	err = json.Unmarshal(ratesFile, &DecodedJSONRates)

	if err != nil {
		fmt.Println("JSON decode error!")
		return
	}

	// Convert to DailyRate
	//var allDailyRates AllDailyRates

	//Create HTTP server to handle "price" and "rates" endpoints
	http.HandleFunc("/rates", ratesURL) //GET and PUT
	http.HandleFunc("/price", priceURL) //GET
	log.Println("Server running on port" + SERVERPORT)
	log.Fatal(http.ListenAndServe(SERVERPORT, nil))

}

// Converts JSONRate to DailyRate
func JSONRateToDailyRate(jsonRate JSONRate) DailyRate {
	//panic("TODO")
	var d DailyRate

	d.Weekdays = JSONRateDayConversion(jsonRate.Days)

	start, end := JSONRateTimesToStartAndEndTime(jsonRate.Times)
	d.StartTime = JSONRateTimeConversion(start, jsonRate.Tz)
	d.EndTime = JSONRateTimeConversion(end, jsonRate.Tz)

	d.Price = jsonRate.Price
	return d
}

// Converts JSONRate.Days to slice of Weekdays for DailyRate.Weekdays
func JSONRateDayConversion(weekday string) []time.Weekday {
	panic("TODO")
}

// inputTime should be in the format: "2359" for 11:59pm
func JSONRateTimeConversion(inputTime string, timeZone string) time.Time {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		// ... handle error
		panic(err)
	}

	hour, err := strconv.Atoi(inputTime[:2])
	if err != nil {
		// ... handle error
		panic(err)
	}

	minute, err := strconv.Atoi(inputTime[2:])
	if err != nil {
		// ... handle error
		panic(err)
	}

	completedTime := time.Date(2000, time.January, 1, hour, minute, 0, 0, location)

	return completedTime

}

// Given "0900-2100" from JSONRate, should return "0900" as startTime and "2100" as endTime
func JSONRateTimesToStartAndEndTime(inputTime string) (startTime string, endTime string) {
	stringSlice := strings.Split(inputTime, "-")
	startTime, endTime = stringSlice[0], stringSlice[1]
	return startTime, endTime
}

// TODO: use to load time
func ParseJSONDateTime(inputTime string) time.Time {
	loc, _ := time.LoadLocation("Europe/Paris")
	parisTime := time.Date(2022, time.July, 14, 9, 30, 0, 0, loc)
	fmt.Println(parisTime)
	return parisTime
}
