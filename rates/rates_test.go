package rates

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	. "github.com/Freddy-S3/PriceAPI/common"
)

// Sends a GET request to RatesURL and tests if it returns the stored DB rates
// Did not finish debugging or adding PUT tests
func TestRatesURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	// Record result of getAllRates
	w := httptest.NewRecorder()
	RatesURL(w, req)
	res := w.Result()
	defer res.Body.Close()

	// Check the status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the content type
	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got %s", contentType)
	}

	fmt.Println(res.Body)
	// Result JSONRates Response
	resultJSONRates, err := JSONFileToJSONRates(res.Body)

	if err != nil {
		t.Errorf("Response Formatting incorrect")
	}

	// Read the expected JSON file result
	ratesFile, err := os.Open("priceDB.json")
	if err != nil {
		t.Errorf("Unable to load priceDB.json")
	}
	defer ratesFile.Close()

	// Convert the JSON file to JSONRates struct
	expectedJSONRates, err := JSONFileToJSONRates(ratesFile)
	if err != nil {
		t.Errorf("priceDB.json not formatted correctly")
	}

	//
	resultRateValue := reflect.ValueOf(resultJSONRates.Rates)
	//resultRateType := resultRateValue.Type()
	expectedRateValue := reflect.ValueOf(expectedJSONRates.Rates)

	for i := 0; i < resultRateValue.NumField(); i++ {
		if expectedRateValue.Field(i).Interface() != resultRateValue.Field(i).Interface() {
			t.Errorf("PRICE FAIL: GetAllRates(w) = %v; want %v", resultJSONRates, expectedJSONRates)
		}
	}
}

// incomplete Unit tests
func TestGetAllRates(t *testing.T) {
	//http response test
}

func TestUpdateRate(t *testing.T) {
	//http response test
}
