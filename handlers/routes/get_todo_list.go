package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/v-inmar/goweb/repo"
)

func GetAllTodos(rw http.ResponseWriter, r *http.Request) {
	// Set the Headers
	rw.Header().Set("Content-Type", "application/json")

	// Write the status code
	rw.WriteHeader(http.StatusOK)

	json.NewEncoder(rw).Encode(repo.Todos)
}
