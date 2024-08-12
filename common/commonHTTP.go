package common

import (
	"encoding/json"
	"net/http"
)

// Returns "unavailable" as a JSON response
func HttpResponseOfUnavailable(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("unavailable")
	return w
}
