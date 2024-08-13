package price

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/Freddy-S3/PriceAPI/common"
)

const URLDATETIMELAYOUT string = "2006-01-02T15:04:05-07:00"

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
// Ex query: http://localhost:5000/price?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00 //"1750"
// Second ex: http://localhost:5000/price?start=2015-07-01T01:00:00-05:00&end=2015-07-01T05:00:00-05:00 //"1000"
// Third ex: http://localhost:5000/price?start=2015-07-01T19:00:00-05:00&end=2015-07-01T20:00:00-05:00 //"unavailable"
func PriceURL(w http.ResponseWriter, r *http.Request) {
	// Parse start URL
	startStr := r.URL.Query().Get("start")
	startTime, err := time.Parse(URLDATETIMELAYOUT, startStr)
	if err != nil {
		http.Error(w, "Invalid 'start' time format", http.StatusBadRequest) //can replace with unavailable if required
		return
	}

	// Parse end URL #potential refactor
	endStr := r.URL.Query().Get("end")
	endTime, err := time.Parse(URLDATETIMELAYOUT, endStr)
	if err != nil {
		http.Error(w, "Invalid 'end' time format", http.StatusBadRequest)
		return
	}

	// Check if endtime is before or equal to start time
	if endTime.Before(startTime) || endTime.Equal(startTime) { //Check test case if equal or earlier
		http.Error(w, "'end' before 'start'", http.StatusBadRequest)
		return
	}

	// Check if input "spans more than a day" //assuming non-converted time
	if startTime.Day() != endTime.Day() {
		HttpResponseOfUnavailable(w) //return as JSON?
		return
	}

	// Read the JSON file
	ratesFile, err := os.Open("priceDB.json")
	if err != nil {
		HttpResponseOfUnavailable(w)
		return
	}
	defer ratesFile.Close()

	// Convert the JSON file to JSONRates struct
	allJSONRates, err := JSONFileToJSONRates(ratesFile)
	if err != nil {
		HttpResponseOfUnavailable(w)
		return
	}
	allDailyRates := AllJSONRatesToDailyRates(allJSONRates.Rates)

	// Look through each dailyRate
	for _, dailyRate := range allDailyRates {
		// Convert query to each rate's timezone
		convertedStartTime := startTime.In(dailyRate.StartTime.Location())
		convertedEndTime := endTime.In(dailyRate.EndTime.Location())

		// "Rates will not span multiple days" so immediately skip
		if convertedStartTime.Day() != convertedEndTime.Day() {
			continue
		}

		// Look through each day of each daily rate #possible refactor for faster lookup, larger DB. Possible channel use
		for _, weekDay := range dailyRate.Weekdays {
			if convertedStartTime.Weekday() == weekDay {
				// Convert Rates into the same day as the request //This assumes that a rate will not span multiple days
				convertedRateStartTime := time.Date(convertedStartTime.Year(), convertedStartTime.Month(), convertedStartTime.Day(), dailyRate.StartTime.Hour(), dailyRate.StartTime.Minute(), 0, 0, convertedStartTime.Location())
				convertedRateEndTime := time.Date(convertedEndTime.Year(), convertedEndTime.Month(), convertedEndTime.Day(), dailyRate.EndTime.Hour(), dailyRate.EndTime.Minute(), 0, 0, convertedEndTime.Location())
				if convertedStartTime.After(convertedRateStartTime) && convertedEndTime.Before(convertedRateEndTime) {
					duration := convertedEndTime.Sub(convertedStartTime).Hours()
					totalCost := float64(dailyRate.Price) * duration
					returnStruct := DailyRate{
						Weekdays:  dailyRate.Weekdays,
						StartTime: dailyRate.StartTime,
						EndTime:   dailyRate.EndTime,
						Price:     int(totalCost), //this will round down
					}
					jsonRatePrice, _ := json.Marshal(returnStruct)
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonRatePrice)
					return
				}
			}
		}
	}
	// "Input can't span more 2 rates" or rate not available for the selected time
	HttpResponseOfUnavailable(w)
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

	dailyRate.Weekdays = jsonRateDaysToTimeWeekdaySlice(jsonRate.Days)

	startTime, endTime := jsonRateTimesToStartAndEndTime(jsonRate.Times)
	dailyRate.StartTime = jsonRateTimeToTimeTime(startTime, jsonRate.Tz)
	dailyRate.EndTime = jsonRateTimeToTimeTime(endTime, jsonRate.Tz)

	dailyRate.Price = jsonRate.Price
	return dailyRate
}

// Converts JSONRate.Days (ex: "mon,tues,wed") to slice of Weekdays for DailyRate.Weekdays
// Please look at the JSONWEEKDAY const at the top of this file to see acceptable date inputs
func jsonRateDaysToTimeWeekdaySlice(weekday string) []time.Weekday {
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
func jsonRateTimeToTimeTime(inputTime string, timeZone string) time.Time {
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
func jsonRateTimesToStartAndEndTime(inputTime string) (startTime string, endTime string) {
	stringSlice := strings.Split(inputTime, "-")
	startTime, endTime = stringSlice[0], stringSlice[1]
	return startTime, endTime
}
