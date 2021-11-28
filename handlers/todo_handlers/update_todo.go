package handler_todo

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/v-inmar/goweb/models"
)



func UpdateTodo(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

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

		reqBody := models.RequestBodyModel{}
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
		if (reqBody == models.RequestBodyModel{}) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := updateDB(db, &reqBody, &todo_id); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// get the todo using the getTodoFromDB function of the get_todo handler
		todo, err := getTodoFromDB(db, &todo_id, &pid)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(todo)
	}
}


// return error when updating the database
func updateDB(db *sql.DB, reqBody *models.RequestBodyModel, todo_id *int)(error){
	// helper function for when something failed
	fail := func(err error) (error) {
		log.Printf("%v", err)
		return err
	}

	// Start transaction (context.Background by default)
	dbSession, err := db.Begin()
	if err != nil {
		return fail(err)
	}
	// Defer a rollback when something fails
	defer dbSession.Rollback()


	// get the utc datetime
	var dt = time.Now().UTC()

	

	// --- deals with the title --- //
	// Passed in Title has value
	// otherwise, just ignore the blank title value and continue using the current title
	if len(reqBody.Title) > 0 {

		// grab the current linker todo title
		linker_title_row := dbSession.QueryRow("select title_id from todo_title_linker_model where todo_id=?", todo_id)
		var current_title_id int64
		if err = linker_title_row.Scan(&current_title_id); err != nil {
			return fail(err)
		}

		// grab the current title model
		current_title_row := dbSession.QueryRow("select value from title_model where id=?", current_title_id)
		var current_title_value string
		if err = current_title_row.Scan(&current_title_value); err != nil {
			return fail(err)
		}


		// check something changed
		// only perform updates when new title does not match current title
		if current_title_value != reqBody.Title{
			// will be the golder of the title_model id
			// int64 to conform to the return of LastInsertId()
			var new_title_id int64

			title_row := dbSession.QueryRow("select id from title_model where value=?", reqBody.Title)
			if err = title_row.Scan(&new_title_id); err != nil {
				if err == sql.ErrNoRows {
					new_title_id = -1
				}else{
					// Failed to either grab the title id after checking for its existence
					return fail(err)
				}
			}

			// Title value does not exist in db and must be created
			if new_title_id == -1 {
				
				// Insert the new title into the title model
				title_result, err := dbSession.Exec("insert into title_model (value, date_created) values (?,?)", reqBody.Title, dt)
				if err != nil {
					return fail(err)
				}

				// set the title_id with the last inserted id
				// db has autoincrement on the id of the title model
				new_title_id, err = title_result.LastInsertId()
				if err != nil {
					return fail(err)
				}
			}

			// update the linker title model
			_, err = dbSession.Exec("update todo_title_linker_model set title_id=?, date_updated=? where todo_id=?", new_title_id, dt, todo_id)
			if err != nil {
				return fail(err)
			}
		}
		
	}


	// --- deals with the body --- //

	var body_id int
	linker_todo_body_row := dbSession.QueryRow("select body_id from todo_body_linker_model where todo_id=?", todo_id)
	if err := linker_todo_body_row.Scan(&body_id); err != nil {
		return fail(err)
	}

	// This is an execute statement that doesn't return any rows, as oppose to query
	// However the DB.Exec function returns values (sql.Result and error)
	// sql.Result will be ignored because it is not needed, only the error
	_, err = dbSession.Exec("update body_model set value=? where id=?", reqBody.Body, body_id)
	if err != nil {
		return fail(err)
	}

	// Update the linker model
	_, err = dbSession.Exec("update todo_body_linker_model set date_updated=? where todo_id=?", dt, todo_id)
	if err != nil {
		return fail(err)
	}


	// Update the todo model's date_updated column
	_, err = dbSession.Exec("update todo_model set date_updated=? where id=?", dt, todo_id)
	if err != nil {
		return fail(err)
	}

	if err = dbSession.Commit(); err != nil {
		return fail(err)
	}
	return nil
}