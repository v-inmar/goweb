package handler_auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/models/user_linker_models"
	"github.com/v-inmar/goweb/models/user_models"
	"github.com/v-inmar/goweb/utils/hash_utils"
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

		emailModel := user_models.EmailModel{}
		err := emailModel.ReadByValue(db, strings.ToLower(reqBody.EmailAddress))
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		
		if (emailModel != user_models.EmailModel{}){
			emailLinkerModel := user_linker_models.EmailLinkerModel{}
			err := emailLinkerModel.ReadByEmailId(db, emailModel.ID)
			if err != nil{
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Compare to empty..not empty, means email not available
			if (emailLinkerModel != user_linker_models.EmailLinkerModel{}){
				rw.WriteHeader(http.StatusConflict)
				return	
			}
		}

		// Hash the incoming password
		hashed, err := hash_utils.PasswordHash(reqBody.Password)
		if err != nil{
			rw.WriteHeader(http.StatusConflict)
			return
		}

		// ----- If it gets here, the email address is available ----- //

		// Helper method for rolling back and returning response
		// with error status code
		failedAndRollback := func(tx *sql.Tx, statusCode int){
			tx.Rollback()
			rw.WriteHeader(statusCode)
			return
		}

		dbSession, err := db.Begin()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// ### User Model ### //
		userModel := user_models.UserModel{}
		if err := userModel.Create(dbSession); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		// ### Email Model ### //
		if (emailModel == user_models.EmailModel{}){
			if err := emailModel.Create(dbSession, strings.ToLower(reqBody.EmailAddress)); err != nil{
				failedAndRollback(dbSession, http.StatusInternalServerError)
			}
		}

		emailLinkerModel := user_linker_models.EmailLinkerModel{}
		if err := emailLinkerModel.Create(dbSession, userModel.ID, emailModel.ID); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}


		// ### Password Model ### //
		passwordModel := user_models.PasswordModel{}
		if err := passwordModel.Create(dbSession, hashed); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		passwordLinkerModel := user_linker_models.PasswordLinkerModel{}
		if err := passwordLinkerModel.Create(dbSession, userModel.ID, passwordModel.ID); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}


		// ### Firstname Model ### //
		firstnameModel := user_models.FirstnameModel{}
		if err := firstnameModel.ReadByValue(db, reqBody.Firstname); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		// Compare to empty
		if (firstnameModel == user_models.FirstnameModel{}) {
			// Create new
			if err := firstnameModel.Create(db, reqBody.Firstname); err != nil{
				failedAndRollback(dbSession, http.StatusInternalServerError)
			}
		}

		firstnameLinkerModel := user_linker_models.FirstnameLinkerModel{}
		if err := firstnameLinkerModel.Create(dbSession, userModel.ID, firstnameModel.ID); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		// ### Lastname Model ### //
		lastnameModel := user_models.LastnameModel{}
		if err := lastnameModel.ReadByValue(db, reqBody.Lastname); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		// Compare to empty
		if (lastnameModel == user_models.LastnameModel{}) {
			// Create new
			if err := lastnameModel.Create(dbSession, reqBody.Lastname); err != nil{
				failedAndRollback(dbSession, http.StatusInternalServerError)
			}
		}

		lastnameLinkerModel := user_linker_models.LastnameLinkerModel{}
		if err := lastnameLinkerModel.Create(dbSession, userModel.ID, firstnameModel.ID); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}


		// ### UPID Model ### //
		
		
		count := 0
		upidCreateSuccess := false
		// this will only run it 5 times
		// to make sure loop doesn't run forever
		for ok := true; ok; ok = (count != 5){
			upidModel := user_models.UPIDModel{}

			// generate 8 random characters
			upid := strings.Replace(uuid.NewString(),"-", "", -1)[0:8]
			if err := upidModel.ReadByValue(db, upid); err != nil{
				failedAndRollback(dbSession, http.StatusInternalServerError)
			}

			// Compare to empty
			if (upidModel == user_models.UPIDModel{}){
				if err := upidModel.Create(dbSession, upid); err != nil{
					failedAndRollback(dbSession, http.StatusInternalServerError)
				}

				upidLinkerModel := user_linker_models.UPIDLinkerModel{}
				if err := upidLinkerModel.Create(dbSession, userModel.ID, upidModel.ID); err != nil{
					failedAndRollback(dbSession, http.StatusInternalServerError)
				}
				upidCreateSuccess = true
				break
			}

		}

		if !upidCreateSuccess{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}



		if err := dbSession.Commit(); err != nil{
			failedAndRollback(dbSession, http.StatusInternalServerError)
		}

		// TODO: Login user here (produce jwt tokens for access and refresh)
	}
}