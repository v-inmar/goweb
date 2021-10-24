package handlers

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/v-inmar/goweb/models"
)

type requestBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Handles the create route
func CreateTodo(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		// Set content type
		rw.Header().Set("Content-Type", "application/json")

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

		publicModel, err := insertIntoDB(db, &reqBody)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(publicModel)
		return

	}
}

/*
Insert new todo into the database
NOTE: Glimpse of no ORM help
NOTE: Probably sql stored procedure can help here
*/
func insertIntoDB(db *sql.DB, reqBody *requestBody) (models.PublicTodoModel, error) {
	// helper function for when something failed
	fail := func(err error) (models.PublicTodoModel, error) {
		return models.PublicTodoModel{}, err
	}

	// Start transaction (context.Background by default)
	dbSession, err := db.Begin()
	if err != nil {
		return fail(err)
	}
	// Defer a rollback when something fails
	defer dbSession.Rollback()

	// generate random pid string (16 characters long value)
	// characters to use
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for { // loop forever

		// create empty rune array size 16
		pidRune := make([]rune, 16)

		// loop through the pidRune array
		// inserting random character
		for i := range pidRune {
			pidRune[i] = characters[rand.Intn(len(characters))]
		}

		//convert rune array to string
		pidString := string(pidRune)

		row := db.QueryRow("select id from pid_model where value = ?", pidString)
		var pid_row_string string
		if err := row.Scan(&pid_row_string); err != nil {
			if err == sql.ErrNoRows {

				var dt = time.Now().UTC()

				// create todo model item
				todo_result, err := dbSession.Exec("insert into todo_model (date_created) values (?)", dt)
				if err != nil {
					return fail(err)
				}

				// get the todo_model item id
				todo_id, err := todo_result.LastInsertId()
				if err != nil {
					return fail(err)
				}

				// create pid model item using the unique value string
				pid_result, err := dbSession.Exec("insert into pid_model (value, date_created) values (?,?)", pidString, dt)
				if err != nil {
					return fail(err)
				}

				// get the pid model item id
				pid_id, err := pid_result.LastInsertId()
				if err != nil {
					return fail(err)
				}

				// create todo pid linker
				// result is not needed
				_, err = dbSession.Exec("insert into todo_pid_linker_model (todo_id, pid_id, date_created, date_updated) values (?,?,?,?)", todo_id, pid_id, dt, dt)
				if err != nil {
					return fail(err)
				}

				var title_id int64
				// check if title already exist and just grab the id
				if err = dbSession.QueryRow("select id from title_model where value = ?", reqBody.Title).Scan(&title_id); err != nil {
					if err == sql.ErrNoRows {

						// create title mode item
						title_result, err := dbSession.Exec("insert into title_model (value, date_created) values (?,?)", reqBody.Title, dt)
						if err != nil {
							return fail(err)
						}

						title_id, err = title_result.LastInsertId()
						if err != nil {
							return fail(err)
						}

					} else {
						return fail(err)
					}
				}

				// create todo and title linker
				_, err = dbSession.Exec("insert into todo_title_linker_model (todo_id, title_id, date_created, date_updated) values (?,?,?,?)", todo_id, title_id, dt, dt)
				if err != nil {
					return fail(err)
				}

				// create body model item
				body_result, err := dbSession.Exec("insert into body_model (value, date_created) values (?,?)", reqBody.Body, dt)
				if err != nil {
					return fail(err)
				}

				// get the body model item id
				body_id, err := body_result.LastInsertId()
				if err != nil {
					return fail(err)
				}

				// create todo body linker item
				_, err = dbSession.Exec("insert into todo_body_linker_model (todo_id, body_id, date_created, date_updated) values (?,?,?,?)", todo_id, body_id, dt, dt)
				if err != nil {
					return fail(err)
				}

				if err = dbSession.Commit(); err != nil {
					return fail(err)
				}

				// Change this to a better approach
				var publicModel models.PublicTodoModel
				publicModel.PID = pidString
				publicModel.Title = reqBody.Title
				publicModel.Body = reqBody.Body
				return publicModel, nil

			}
		}
	}
}

// func CreateTodo(rw http.ResponseWriter, r *http.Request) {
// 	// Set content type
// 	rw.Header().Set("Content-Type", "application/json")

// 	// Matches the request body
// 	// Note: similar to TodoModel struct but no ID
// 	type requestBody struct {
// 		Title string `json:"title"`
// 		Body  string `json:"body"`
// 	}

// 	reqBody := requestBody{}

// 	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
// 		// Empty request body was sent
// 		if err.Error() == "EOF" {
// 			rw.WriteHeader(http.StatusBadRequest)
// 		} else {

// 			// Any other error
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Check for empty json request body
// 	if (reqBody == requestBody{}) {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Make sure there are values
// 	if len(reqBody.Title) < 1 || len(reqBody.Body) < 1 {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Figure out a new id using random integer, seeded with time
// 	var newID int
// 	for {
// 		// setup
// 		isCollision := false
// 		rand.Seed(time.Now().UnixNano())
// 		newID = rand.Intn(len(repo.Todos) + 100)

// 		// check it doesnt exist in current repo value id
// 		for _, value := range repo.Todos {
// 			if value.ID == newID {
// 				// no point of continuing with loop
// 				// it already hit a collision
// 				isCollision = true
// 				break
// 			}
// 		}

// 		// check if there was a collision
// 		if !isCollision {
// 			break
// 		}
// 	}

// 	// Add new todo in the repo
// 	repo.Todos = append(repo.Todos, models.TodoModel{
// 		ID:    newID,
// 		Title: reqBody.Title,
// 		Body:  reqBody.Body,
// 	})

// 	// find the newly added todo and add it to the response
// 	for _, value := range repo.Todos {
// 		if value.ID == newID {
// 			rw.WriteHeader(http.StatusCreated)
// 			json.NewEncoder(rw).Encode(value)
// 			return
// 		}
// 	}

// 	// Reached here means error
// 	rw.WriteHeader(http.StatusInternalServerError)
// 	return
// }
