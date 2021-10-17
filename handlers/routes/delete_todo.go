package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/repo"
)

func DeleteTodo(rw http.ResponseWriter, r *http.Request) {

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

	// Loop through the slice
	for i, values := range repo.Todos {

		// Remove from the slice
		if values.ID == id {
			repo.Todos[i] = repo.Todos[len(repo.Todos)-1]
			repo.Todos[len(repo.Todos)-1] = models.TodoModel{}
			repo.Todos = repo.Todos[:len(repo.Todos)-1]

			// Set to no content
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// No match
	rw.WriteHeader(http.StatusNotFound)
	return
}
