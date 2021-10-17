package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/repo"
)

func CreateTodo(rw http.ResponseWriter, r *http.Request) {
	// Set content type
	rw.Header().Set("Content-Type", "application/json")

	// Matches the request body
	// Note: similar to TodoModel struct but no ID
	type requestBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	reqBody := requestBody{}

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

	// Make sure there are values
	if len(reqBody.Title) < 1 || len(reqBody.Body) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Figure out a new id using random integer, seeded with time
	var newID int
	for {
		// setup
		isCollision := false
		rand.Seed(time.Now().UnixNano())
		newID = rand.Intn(len(repo.Todos) + 100)

		// check it doesnt exist in current repo value id
		for _, value := range repo.Todos {
			if value.ID == newID {
				// no point of continuing with loop
				// it already hit a collision
				isCollision = true
				break
			}
		}

		// check if there was a collision
		if !isCollision {
			break
		}
	}

	// Add new todo in the repo
	repo.Todos = append(repo.Todos, models.TodoModel{
		ID:    newID,
		Title: reqBody.Title,
		Body:  reqBody.Body,
	})

	// find the newly added todo and add it to the response
	for _, value := range repo.Todos {
		if value.ID == newID {
			rw.WriteHeader(http.StatusCreated)
			json.NewEncoder(rw).Encode(value)
			return
		}
	}

	// Reached here means error
	rw.WriteHeader(http.StatusInternalServerError)
	return
}
