package handler_todo

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/v-inmar/goweb/models"
)

func GetTodo(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r) // get url variables

		pid, containsKey := vars["id"] // check for key id
		if !containsKey {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(pid) < 1 { // check for valid length of id
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// grab pid model item
		pid_row := db.QueryRow("select id from pid_model where value = ?", pid)
		var pid_id int
		if err := pid_row.Scan(&pid_id); err != nil {
			if err == sql.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// grab todo_id from linker pid and todo
		linker_row := db.QueryRow("select todo_id from todo_pid_linker_model where pid_id=?", pid_id)
		var linker_todo_id int
		if err := linker_row.Scan(&linker_todo_id); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// grab todo_model id
		todo_row := db.QueryRow("select id from todo_model where id=? and date_deleted is null", linker_todo_id)
		var todo_id int
		if err := todo_row.Scan(&todo_id); err != nil {
			if err == sql.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		
		// make db call for the todo and cast it to the response body model
		todo, err := getTodoFromDB(db, &todo_id, &pid)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(todo)

	}
}

func getTodoFromDB(db *sql.DB, todo_id *int, pid *string)(models.ResponseBodyModel, error){

	var date_created time.Time
	var date_updated sql.NullTime // Initially all date_updated columns are sql null
	todo_row := db.QueryRow("select date_created, date_updated from todo_model where id=?", todo_id)
	if err := todo_row.Scan(&date_created, &date_updated); err != nil {
		return models.ResponseBodyModel{}, err
	}

	// -- deals with title
	linker_title_row := db.QueryRow("select title_id from todo_title_linker_model where todo_id=?", todo_id)
	var linker_title_title_id int
	if err := linker_title_row.Scan(&linker_title_title_id); err != nil {
		return models.ResponseBodyModel{}, err
	}

	title_row := db.QueryRow("select value from title_model where id=?", linker_title_title_id)
	var title string
	if err := title_row.Scan(&title); err != nil {
		return models.ResponseBodyModel{}, err 
	}


	// -- deals with body
	linker_body_row := db.QueryRow("select body_id from todo_body_linker_model where todo_id=?", todo_id)
	var linker_body_body_id string
	if err := linker_body_row.Scan(&linker_body_body_id); err != nil {
		return models.ResponseBodyModel{}, err
	}

	body_row := db.QueryRow("select value from body_model where id=?", linker_body_body_id)
	var body string
	if err := body_row.Scan(&body); err != nil {
		return models.ResponseBodyModel{}, err
	}

	// Check the date_updated for mysql null value
	var updated_date_value string
	if date_updated.Valid {
		updated_date_value = date_updated.Time.Format(time.RFC822)
	}else{
		updated_date_value = ""
	}

	return models.ResponseBodyModel{
		PID: *pid,
		Body: body,
		Title: title,
		Created: date_created.Format(time.RFC822),
		Updated: updated_date_value,
	}, nil


}

