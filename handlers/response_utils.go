package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// jsonResponse writes a JSON response
func jsonResponse(w http.ResponseWriter, response map[string]int) {
	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
