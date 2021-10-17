package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		// Empty request body was sent
		if err.Error() == "EOF" {
			rw.WriteHeader(http.StatusBadRequest)
		} else {

			// Any other error
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Check for empty json request body
	if (reqBody == requestBody{}) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for index := range repo.Todos {

		// Pointer to the model to make it mutable
		valuePointer := &repo.Todos[index]
		if valuePointer.ID == id {
			if len(reqBody.Title) > 0 {
				valuePointer.Title = reqBody.Title
			}

			if len(reqBody.Body) > 0 {
				valuePointer.Body = reqBody.Body
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
