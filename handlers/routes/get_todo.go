package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/v-inmar/goweb/repo"
)

func GetTodo(rw http.ResponseWriter, r *http.Request) {
	// Set content type
	rw.Header().Set("Content-Type", "application/json")

	// grab path params
	vars := mux.Vars(r)

	// convert to int
	id, error := strconv.Atoi(vars["id"])
	// Check for convertable
	if error != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, value := range repo.Todos {
		if value.ID == id {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(value)
			return
		}
	}

	// No match
	rw.WriteHeader(http.StatusNotFound)
	return
}
