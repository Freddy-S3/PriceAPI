package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const URLDATETIMELAYOUT string = "2006-01-02T15:04:05-07:00"

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
	Weekdays  []time.Weekday `json:"-"` // #potential refactor to be a single day
	StartTime time.Time      `json:"-"`
	EndTime   time.Time      `json:"-"`
	Price     int            `json:"price"`
}

// A Map to translate the JSON "days" string into time.Weekday
var JSONWEEKDAY = map[string]time.Weekday{
	"sun":   time.Sunday,
	"mon":   time.Monday,
	"tues":  time.Tuesday,
	"wed":   time.Wednesday,
	"thurs": time.Thursday,
	"fri":   time.Friday,
	"sat":   time.Saturday,
}

// GET /price: return JSON for price given time query
// ex query: http://localhost:5000/price?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00 //"1750"
// second ex: http://localhost:5000/price?start=2015-07-01T01:00:00-05:00&end=2015-07-01T05:00:00-05:00 //"1000"
// third ex: http://localhost:5000/price?start=2015-07-01T19:00:00-05:00&end=2015-07-01T20:00:00-05:00 //"unavailable"
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

	// check if input wrong
	if endTime.After(startTime) { //check test case if equal or earlier
		//if input spans more than a day //assuming non-converted time
		if startTime.Day() != endTime.Day() {
			httpResponseOfUnavailable(w) //return as JSON?
			return
		}
	} else {
		http.Error(w, "'end' before 'start'", http.StatusBadRequest)
		return
	}

	// read the JSON file
	ratesFile, err := os.ReadFile("testing.json")
	if err != nil {
		httpResponseOfUnavailable(w)
		return
	}

	//look through the time bands //input can't span more 2 time bands //assuming DB time zone?
	AllDailyRates := AllJSONRatesToDailyRates(JSONFileToJSONRates(ratesFile))

	for _, dailyRate := range AllDailyRates {
		convertedStartTime := startTime.In(dailyRate.StartTime.Location())
		convertedEndTime := endTime.In(dailyRate.EndTime.Location())

		// "Rates will not span multiple days" so immediatly skip
		if convertedStartTime.Day() != convertedEndTime.Day() {
			continue
		}

		for _, weekDay := range dailyRate.Weekdays {
			if convertedStartTime.Weekday() == weekDay {
				//extract hour and minute for comparison //#refactor to account for seconds
				if convertedStartTime.Hour() >= dailyRate.StartTime.Hour() && convertedEndTime.Hour() <= dailyRate.EndTime.Hour() && convertedStartTime.Minute() >= dailyRate.StartTime.Minute() && convertedEndTime.Minute() <= dailyRate.EndTime.Minute() {
					fmt.Println(convertedStartTime.Weekday(), convertedStartTime)
					fmt.Println(convertedEndTime.Weekday(), convertedEndTime)
					w.Header().Set("Content-Type", "application/json")
					tester, _ := json.Marshal(dailyRate)
					w.Write(tester)
					//json.NewEncoder(w).Encode()
					return
				} else {
					continue
				}
			}
		}
	}
	httpResponseOfUnavailable(w)
}

// returns "unavailable" as a JSON response
func httpResponseOfUnavailable(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("unavailable")
	return w
}

// Decodes JSON file into JSONRates struct
func JSONFileToJSONRates(ratesFile []byte) []JSONRate {
	var decodedJSONRates AllJSONRates
	err := json.Unmarshal(ratesFile, &decodedJSONRates)

	if err != nil {
		panic("JSON decode error!")
	}

	return decodedJSONRates.Rates
}

// Converts decoded JSONRate slice into DailyRate slice for easier date and time lookup
func AllJSONRatesToDailyRates(jsonRates []JSONRate) []DailyRate {

	var allDailyRates []DailyRate
	for _, decodedRate := range jsonRates {
		dailyRate := JSONRateToDailyRate(decodedRate)
		allDailyRates = append(allDailyRates, dailyRate)
	}

	return allDailyRates
}

// Converts JSONRate to DailyRate
func JSONRateToDailyRate(jsonRate JSONRate) DailyRate {
	var dailyRate DailyRate

	dailyRate.Weekdays = JSONRateDaysToTimeWeekdaySlice(jsonRate.Days)

	startTime, endTime := JSONRateTimesToStartAndEndTime(jsonRate.Times)
	dailyRate.StartTime = JSONRateTimeToTimeTime(startTime, jsonRate.Tz)
	dailyRate.EndTime = JSONRateTimeToTimeTime(endTime, jsonRate.Tz)

	dailyRate.Price = jsonRate.Price
	return dailyRate
}

// Converts JSONRate.Days to slice of Weekdays for DailyRate.Weekdays
func JSONRateDaysToTimeWeekdaySlice(weekday string) []time.Weekday {
	timeSlice := []time.Weekday{}
	stringSlice := strings.Split(weekday, ",")
	for _, str := range stringSlice {
		timeSlice = append(timeSlice, JSONWEEKDAY[str])
	}

	return timeSlice
}

// Takes a 4 digit input for hour and minutes (from half of the JSONRate.Times)
// and a time zone from JSONRate.Tz,
// and converts the 4 to a time.Time object
// Ex: "0900" and "America/Chicago" returns time.Time object of 2000, Jan 1, 9:00:00, -6:00CST
func JSONRateTimeToTimeTime(inputTime string, timeZone string) time.Time {
	hour, err := strconv.Atoi(inputTime[:2])
	if err != nil {
		panic("hour format incorrect, please ensure the inputTime is 4 digits. ex: 1954 for 7:54PM")
	}

	minute, err := strconv.Atoi(inputTime[2:])
	if err != nil {
		panic("minute format incorrect, please ensure the inputTime is 4 digits. ex: 1954 for 7:54PM")
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		panic("timeZone input incorrect")
	}

	//broke convention for readability
	completedTime := time.Date(2000, time.January, 1, hour, minute, 0, 0, location)

	return completedTime
}

// Given "0900-2100" from JSONRate,
// returns "0900" as startTime and "2100" as endTime
func JSONRateTimesToStartAndEndTime(inputTime string) (startTime string, endTime string) {
	stringSlice := strings.Split(inputTime, "-")
	startTime, endTime = stringSlice[0], stringSlice[1]
	return startTime, endTime
}
