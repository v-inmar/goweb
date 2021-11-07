package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// NOTE: This does not actually delete the items from db
// it just adds the date it was deleted (basically hiding)
func DeleteTodo(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r) // Get url variables, in this case id
		pid, containsKey := vars["id"]

		// check if id is a key in the vars map
		if !containsKey{
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// check if pid is not empty
		if len(pid) < 1 {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// Grab the pid_model row
		pid_row := db.QueryRow("select id from pid_model where value=?", pid)
		var pid_id int
		if err := pid_row.Scan(&pid_id); err != nil {
			if err == sql.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Grab the linker row of todo_model and pid_model
		linker_pid_row := db.QueryRow("select todo_id from todo_pid_linker_model where pid_id=?", pid_id)
		var linker_pid_todo_id int
		if err := linker_pid_row.Scan(&linker_pid_todo_id); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// "delete" from db
		if err := deleteFromDB(db, &linker_pid_todo_id); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}


// Update the row to populate the date deleted
func deleteFromDB(db *sql.DB, todo_id *int)(error){

	// Start transaction
	dbSession, errSession := db.Begin()
	if errSession != nil{
		return errSession
	}
	defer dbSession.Rollback()

	// create datetime variable
	var dt = time.Now().UTC()

	var err error // empty error variable

	// Execute sql update query
	_, err = db.Exec("update todo_model set date_deleted = ? where id = ?", dt, todo_id)
	if err != nil {
		return err
	}

	// commit
	if err = dbSession.Commit(); err != nil {
		return err
	}
	return nil
}
