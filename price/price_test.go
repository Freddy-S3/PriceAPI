package price

import (
	"reflect"
	"testing"
	"time"

	. "github.com/Freddy-S3/PriceAPI/common"
)

// Converts JSONRate to DailyRate
func TestJSONRateToDailyRate(t *testing.T) {
	const testTimeZone string = "America/Chicago"
	location, _ := time.LoadLocation(testTimeZone)

	tests := []struct {
		input  JSONRate
		result DailyRate
	}{
		{
			JSONRate{ //unit test
				Days:  "mon,tues,thurs",
				Times: "0900-2100",
				Tz:    "America/Chicago",
				Price: 1500,
			}, DailyRate{
				Weekdays:  []time.Weekday{time.Monday, time.Tuesday, time.Thursday},
				StartTime: time.Date(2000, time.January, 1, 9, 0, 0, 0, location),
				EndTime:   time.Date(2000, time.January, 1, 21, 0, 0, 0, location),
				Price:     1500,
			},
		}, {
			JSONRate{ //integration test
				Days:  "wed",
				Times: "0600-1800",
				Tz:    "America/Chicago",
				Price: 1750,
			}, DailyRate{
				Weekdays:  jsonRateDaysToTimeWeekdaySlice("wed"),
				StartTime: jsonRateTimeToTimeTime("0600", "America/Chicago"),
				EndTime:   jsonRateTimeToTimeTime("1800", "America/Chicago"),
				Price:     1750,
			},
		},
	}

	for _, test := range tests {
		dailyRate := JSONRateToDailyRate(test.input)

		dailyRateValue := reflect.ValueOf(dailyRate)
		dailyRateType := dailyRateValue.Type()
		resultValue := reflect.ValueOf(test.result)

		for i := 0; i < dailyRateValue.NumField(); i++ {
			switch structFieldType := dailyRateType.Field(i).Type; structFieldType {
			case reflect.TypeOf([]time.Weekday{}):
				for index, weekday := range dailyRate.Weekdays {
					if test.result.Weekdays[index] != weekday {
						t.Errorf("WEEKDAY FAIL: JSONRateToDailyRate(%v)= %v; want %v", test.input, dailyRate, test.result)
					}
				}
			case reflect.TypeOf(time.Time{}):
				if !dailyRateValue.Field(i).Interface().(time.Time).Equal(resultValue.Field(i).Interface().(time.Time)) {
					t.Errorf("TIME FAIL: JSONRateToDailyRate(%v)= %v; want %v", test.input, dailyRate, test.result)
				}
			default:
				if dailyRateValue.Field(i).Interface() != resultValue.Field(i).Interface() {
					t.Errorf("PRICE FAIL: JSONRateToDailyRate(%v)= %v; want %v", test.input, dailyRate, test.result)
				}
			}
		}
	}
}

// Converts JSONRate.Days to slice of Weekdays for DailyRate.Weekdays
func TestJSONRateDaysToTimeWeekdaySlice(t *testing.T) {
	tests := []struct {
		inputWeekday string
		result       []time.Weekday
	}{
		{"mon,tues,thurs", []time.Weekday{time.Monday, time.Tuesday, time.Thursday}},
		{"fri,sat,sun", []time.Weekday{time.Friday, time.Saturday, time.Sunday}},
		{"wed", []time.Weekday{time.Wednesday}},
		{"sun,tues", []time.Weekday{time.Sunday, time.Tuesday}},
	}

	for _, test := range tests {
		weekdays := jsonRateDaysToTimeWeekdaySlice(test.inputWeekday)
		for index, weekday := range weekdays {
			if weekday != test.result[index] {
				t.Errorf("JSONRateDaysToTimeWeekdaySlice(%s)= %v; want %v", test.inputWeekday, weekdays, test.result)
			}
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
		{"0000", "America/Toronto", time.Date(2000, time.January, 1, 5, 0, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		time := jsonRateTimeToTimeTime(test.inputTime, test.inputTimeZone)
		if !time.Equal(test.result) {
			t.Errorf("JSONRateTimeToTimeTime(%s, %s)= %v; want %v", test.inputTime, test.inputTimeZone, time, test.result)
		}
	}
}

// Takes JSONRate.Times like "0900-2100" and convert to 2 strings
// result: startTime and endTime
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
		start, end := jsonRateTimesToStartAndEndTime(test.input)
		if start != test.startOutput || end != test.endOutput {
			t.Errorf("JSONRateTimesToStartAndEndTime(%s)= %s, %s; want %s, %s", test.input, start, end, test.startOutput, test.endOutput)
		}
	}
}
