package handler_auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/v-inmar/goweb/models"
)

func SignupAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		reqBody := models.RequestSignUpBodyModel{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil{
			if err.Error() == "EOF"{
				rw.WriteHeader(http.StatusBadRequest)
			}else{
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// Check for empty request body
		if (reqBody == models.RequestSignUpBodyModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		
		/*
		Makre sure all fields have values
		*/
		if len(reqBody.Firstname) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(reqBody.Lastname) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(reqBody.EmailAddress) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(reqBody.Password) < 8 { // Password must be atleast 8 characters long
			rw.WriteHeader(http.StatusBadRequest)
			return
		}


		// Check email is not in use
		// NOTE: This is probably better to do in a sql procedure
		available := false
		email_id_row := db.QueryRow("select id from email_model where value=?", strings.ToLower(reqBody.EmailAddress))
		var email_id int64
		if err := email_id_row.Scan(&email_id); err != nil {
			if err == sql.ErrNoRows {
				available = true
			}else{
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}else{
			linker_id_row := db.QueryRow("select id from user_email_linker_model where email_id=?", email_id)
			var linker_id int64
			if err := linker_id_row.Scan(&linker_id); err != nil{
				if err == sql.ErrNoRows{
					available = true
				}else{
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		// email address is already in use
		if !available{
			rw.WriteHeader(http.StatusConflict)
			return	
		}



	}
}