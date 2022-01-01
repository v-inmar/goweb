package handler_auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/v-inmar/goweb/models"
	"github.com/v-inmar/goweb/models/user_linker_models"
	"github.com/v-inmar/goweb/models/user_models"
	"github.com/v-inmar/goweb/utils/hash_utils"
	"github.com/v-inmar/goweb/utils/jwt_utils"
)

func LoginAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr
		respBodErrMsgModel := models.ResponseBodyErrorMessageModel{}
		rw.Header().Set("Content-Type", "application/json")

		reqBody := models.RequestLoginBodyModel{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil{
			if err.Error() == "EOF"{
				rw.WriteHeader(http.StatusBadRequest)
				respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			}else{
				rw.WriteHeader(http.StatusInternalServerError)
				respBodErrMsgModel.InternalServerError("Encountered Server Error")
			}
			// encode the error message response
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// empty request body
		if (reqBody == models.RequestLoginBodyModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// empty email field
		if len(reqBody.EmailAddress) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// empty password field
		if len(reqBody.Password) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// get the email model from database
		emailModel := user_models.EmailModel{}
		if err := emailModel.ReadByValue(db, strings.ToLower(reqBody.EmailAddress)); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// did not match any email from database
		// the model is empty
		if (emailModel == user_models.EmailModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// check if email is linked to a user
		// get the linker
		emailLinkerModel := user_linker_models.EmailLinkerModel{}
		if err := emailLinkerModel.ReadByEmailId(db, emailModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		

		// did not match any
		if (emailLinkerModel == user_linker_models.EmailLinkerModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// create a single user id variable to be used further down the execution flow
		userID := emailLinkerModel.UserID

		// get the password linker
		passwordLinkerModel := user_linker_models.PasswordLinkerModel{}
		if err := passwordLinkerModel.ReadByUserId(db, userID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		

		// the password linker should not be empty
		if (passwordLinkerModel == user_linker_models.PasswordLinkerModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// the password model should not be empty
		passwordModel := user_models.PasswordModel{}
		if err := passwordModel.ReadById(db, passwordLinkerModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// check the password
		if !hash_utils.PasswordCheck(passwordModel.Value, reqBody.Password){
			rw.WriteHeader(http.StatusBadRequest)
			respBodErrMsgModel.BadRequest("Invalid Email Address and/or Password")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// get the auth value to attach to jwt claims
		authLinkerModel := user_linker_models.AuthLinkerModel{}
		if err := authLinkerModel.ReadByUserId(db, userID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// this cannot be empty at all
		// every user has should have this		
		if (authLinkerModel == user_linker_models.AuthLinkerModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		authModel := user_models.AuthModel{}
		if err := authModel.ReadById(db, authLinkerModel.AuthID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// again this cannot be empty
		// this was assigned during signup
		if (authModel == user_models.AuthModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// generate access token
		accessStringToken, err := jwt_utils.GenerateJWT(jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 60).Unix(), // 1 hour
			"addr": addr, // client's address (ip)
			"uid": authModel.Value, // user identification
		})

		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// generate refresh token
		refreshStringToken, err := jwt_utils.GenerateJWT(jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
			"addr": addr, // client's address (ip)
			"uid": authModel.Value, // user identification
		})

		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			respBodErrMsgModel.InternalServerError("Encountered Server Error")
			json.NewEncoder(rw).Encode(respBodErrMsgModel)
			return
		}

		// Set token headers
		rw.Header().Set("X-Access-Token", accessStringToken)
		rw.Header().Set("X-Refresh-Token", refreshStringToken)
		rw.WriteHeader(http.StatusOK) // Success 200 return

		// Also return tokens in the body
		json.NewEncoder(rw).Encode(models.ResponseBodyJWTModel{
			Code: http.StatusOK,
			Status: "OK",
			AccessToken: accessStringToken,
			RefreshToken: refreshStringToken,
		})

	}
}