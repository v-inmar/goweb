package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/v-inmar/goweb/models"
)

func GetAllTodos(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		todos, err := getAllTodoFromDB(db)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(todos)

	}
}

func getAllTodoFromDB(db *sql.DB) ([]models.ResponseBodyModel, error) {
	// Initialize with empty
	// NOTE: If this was not initliazed empty
	// the return JSON from the response writer will be null if there are no data
	// this looks better
	var todos = []models.ResponseBodyModel{}

	todo_model_rows, err := db.Query("select id FROM todo_model where date_deleted is null")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer todo_model_rows.Close()

	for todo_model_rows.Next() {
		var todo_model_id int
		if err := todo_model_rows.Scan(&todo_model_id); err != nil {
			continue
		}


		linker_pid_row := db.QueryRow("select pid_id from todo_pid_linker_model where todo_id = ?", todo_model_id)
		var linker_pid_model_id int
		if err := linker_pid_row.Scan(&linker_pid_model_id); err != nil {
			continue
		}

		pid_row := db.QueryRow("select value from pid_model where id = ?", linker_pid_model_id)
		var pid string
		if err := pid_row.Scan(&pid); err != nil {
			continue
		}

		// use of the helper function from the get_todo handler
		t, err := getTodoFromDB(db, &todo_model_id, &pid)
		if err != nil {
			continue
		}

		todos = append(todos, t)

	}

	return todos, nil
}

