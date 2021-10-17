package handlers

import (
	"net/http"
)

// Handles 404 error
func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	// Set the Headers
	rw.Header().Set("Content-Type", "application/json")

	// Write the status code
	rw.WriteHeader(http.StatusNotFound)
	return
}
