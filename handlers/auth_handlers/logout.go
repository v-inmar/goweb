package handler_auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/models/user_models"
)

/*
Provides logout function for the logout route.
NOTE: The json token will be blacklisted. For now, the
json token will be saved into mysql. But it should be in an
in-memory storage i.e. redis, memcache, etc for faster access,
auto-deletion, backups, etc
Both Access Token and Refresh Token will be blacklisted
*/
func LogoutAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		resBodyErr := models.ResponseBodyErrorMessageModel{}

		accessToken := r.Header.Get("X-Access-Token")
		refreshToken := r.Header.Get("X-Refresh-Token")
		if len(accessToken) < 1{
			rw.WriteHeader(http.StatusBadRequest)
			resBodyErr.BadRequest("No access token")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}
		if len(refreshToken) < 1{
			rw.WriteHeader(http.StatusBadRequest)
			resBodyErr.BadRequest("No refresh token")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}
	

		// create db transaction session
		dbSession, err := db.Begin()
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			resBodyErr.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}

		// initialize blacklist struct model
		blacklistModel := user_models.BlacklistModel{}

		// insert accessToken into database
		if err := blacklistModel.Create(dbSession, accessToken); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			resBodyErr.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}

		// insert refreshToken into database
		if err := blacklistModel.Create(dbSession, refreshToken); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			resBodyErr.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}

		// Commit the transaction
		if err := dbSession.Commit(); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			resBodyErr.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(resBodyErr)
			return
		}
		

		rw.WriteHeader(http.StatusOK) // Success 200 return
		okStatus := models.ResponseBodyStatusModel{}
		okStatus.OK("Logout successful")
		json.NewEncoder(rw).Encode(okStatus)
	}
}