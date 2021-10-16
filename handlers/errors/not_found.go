package handlers

import (
	"encoding/json"
	"net/http"

	m "github.com/v-inmar/goweb/models"
)

// Handles 404 error
func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	// Set the Headers
	rw.Header().Set("Content-Type", "application/json")

	// Write the status code
	rw.WriteHeader(http.StatusNotFound)

	// Create the payload and response data
	pm := m.PayloadModel{
		Payload: "",
	}

	rm := m.ResponseModel{
		Status:  "404 Not Found",
		Payload: pm,
	}

	// encode to json
	json.NewEncoder(rw).Encode(rm)
}
