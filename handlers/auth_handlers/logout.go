package handler_auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/v-inmar/goweb/models/response_models"
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
		responseModel := response_models.StatusModel{}

		accessToken := r.Header.Get("X-Access-Token")
		refreshToken := r.Header.Get("X-Refresh-Token")
		if len(accessToken) < 1{
			rw.WriteHeader(http.StatusUnauthorized)
			responseModel.Unauthorized("Missing access token", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}
		if len(refreshToken) < 1{
			rw.WriteHeader(http.StatusUnauthorized)
			responseModel.Unauthorized("Missing refresh token", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}
	

		// create db transaction session
		dbSession, err := db.Begin()
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// initialize blacklist struct model
		blacklistModel := user_models.BlacklistModel{}

		// insert accessToken into database
		if err := blacklistModel.Create(dbSession, accessToken); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// insert refreshToken into database
		if err := blacklistModel.Create(dbSession, refreshToken); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// Commit the transaction
		if err := dbSession.Commit(); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}
		

		rw.WriteHeader(http.StatusOK) // Success 200 return
		responseModel.OK("Logout successful", "")
		json.NewEncoder(rw).Encode(responseModel)
	}
}