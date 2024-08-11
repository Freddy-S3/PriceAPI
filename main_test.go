package main

import (
	"testing"
	"time"
)

// Takes JSONRate.Times like "0900-2100" and convert to 2 strings:
// startTime and endTime
// Ex: "1900-2300" returns "1900" and "2300"
func TestJSONRateTimesToStartAndEndTime(t *testing.T) {
	tests := []struct {
		input                  string
		startOutput, endOutput string
	}{
		{"1900-2300", "1900", "2300"},
		{"2300-2000", "2300", "2000"},
	}

	for _, test := range tests {
		start, end := JSONRateTimesToStartAndEndTime(test.input)
		if start != test.startOutput || end != test.endOutput {
			t.Errorf("JSONRateTimesToStartAndEndTime(%s)= %s, %s; want %s, %s", test.input, start, end, test.startOutput, test.endOutput)
		}
	}
}

// Takes a 4 digit input for hour and minutes (from half of the JSONRate.Times)
// and a time zone from JSONRate.Tz,
// and converts the 4 to a time.Time object
// Ex: "0900" and "America/Chicago" returns
func TestJSONRateTimeToTimeTime(t *testing.T) {
	tests := []struct {
		inputTime, inputTimeZone string
		result                   time.Time
	}{
		{"1754", "America/Chicago", time.Date(2000, time.January, 1, 23, 54, 0, 0, time.UTC)},
		{"1254", "America/Chicago", time.Date(2000, time.January, 1, 18, 54, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		time := JSONRateTimeToTimeTime(test.inputTime, test.inputTimeZone)
		if !time.Equal(test.result) {
			t.Errorf("JSONRateTimeToTimeTime(%s, %s)= %v; want %v", test.inputTime, test.inputTimeZone, time, test.result)
		}
	}
}
