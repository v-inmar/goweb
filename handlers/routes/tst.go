package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type TodoTST struct {
	PID   string `json:"pid"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func GetAllTodosTst(db *sql.DB) http.HandlerFunc {
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

func getAllTodoFromDB(db *sql.DB) ([]TodoTST, error) {
	// Initialize with empty
	// NOTE: If this was not initliazed empty
	// the return JSON from the response writer will be null if there are no data
	// this looks better
	var todos = []TodoTST{}

	todo_model_rows, err := db.Query("Select id FROM todo_model")
	if err != nil {
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

		linker_title_row := db.QueryRow("select title_id from todo_title_linker_model where todo_id = ?", todo_model_id)
		var linker_title_model_id int
		if err := linker_title_row.Scan(&linker_title_model_id); err != nil {
			continue
		}

		title_row := db.QueryRow("select value from title_model where id = ?", linker_title_model_id)
		var title string
		if err := title_row.Scan(&title); err != nil {
			continue
		}

		linker_body_row := db.QueryRow("select body_id from todo_body_linker_model where todo_id = ?", todo_model_id)
		var linker_body_model_id int
		if err := linker_body_row.Scan(&linker_body_model_id); err != nil {
			continue
		}

		body_row := db.QueryRow("select value from body_model where id = ?", linker_body_model_id)
		var body string
		if err := body_row.Scan(&body); err != nil {
			continue
		}
		var t TodoTST
		t.PID = pid
		t.Title = title
		t.Body = body

		todos = append(todos, t)

	}

	return todos, nil

}

// func GetAllTodos(rw http.ResponseWriter, r *http.Request) {
// 	// Set the Headers
// 	rw.Header().Set("Content-Type", "application/json")

// 	// Write the status code
// 	rw.WriteHeader(http.StatusOK)

// 	json.NewEncoder(rw).Encode(repo.Todos)
// }
