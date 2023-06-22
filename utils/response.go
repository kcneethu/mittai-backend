package utils

import (
	"encoding/json"
	"net/http"
)

// Function to send JSON response
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Function to send error response
func SendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	response := map[string]string{"error": errorMessage}
	SendJSONResponse(w, statusCode, response)
}
