package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON serializes the provided data payload as JSON and writes it to the HTTP response with the specified status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError returns an API error formatted as JSON containing an error message string.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
