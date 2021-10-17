package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/repo"
)

func UpdateTodo(rw http.ResponseWriter, r *http.Request) {
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

	// Matches the request body
	// Note: similar to TodoModel struct but no ID
	// ID must not be changed
	type requestBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	// empty struct of the request body
	reqBody := requestBody{}

	// Decode the incoming body and assign it to the reqBody struct
	// Note the use of & for address
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	for index, value := range repo.Todos {
		if value.ID == id {

			// Change the repo at index
			repo.Todos[index] = models.TodoModel{
				ID:    value.ID,
				Title: reqBody.Title,
				Body:  reqBody.Body,
			}
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(repo.Todos[index])
			return
		}
	}

	// No match
	rw.WriteHeader(http.StatusNotFound)
	return
}
