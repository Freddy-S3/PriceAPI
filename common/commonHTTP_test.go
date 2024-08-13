package common

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Returns "unavailable" as an http.ResponseWriter of a JSON response
// Did not finish debugging
func TestHttpResponseOfUnavailable(t *testing.T) {
	// Record result of getAllRates HttpResponseOfUnavailable
	w := httptest.NewRecorder()
	HttpResponseOfUnavailable(w)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// Check the status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the content type
	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got %s", contentType)
	}

	// Check the contents
	if string(data) != "unavailable" {
		t.Errorf("expected 'unavailable' got %v", string(data))
	}

}
