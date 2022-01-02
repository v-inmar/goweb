package handler_auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/v-inmar/goweb/models/request_models"
	"github.com/v-inmar/goweb/models/response_models"
	"github.com/v-inmar/goweb/models/token_models"
	"github.com/v-inmar/goweb/models/user_linker_models"
	"github.com/v-inmar/goweb/models/user_models"
	"github.com/v-inmar/goweb/utils/hash_utils"
	"github.com/v-inmar/goweb/utils/jwt_utils"
)

func LoginAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		responseModel := response_models.StatusModel{}
		loginModel := request_models.LoginModel{}
		addr := r.RemoteAddr
		// respBodErrMsgModel := models.ResponseBodyErrorMessageModel{}
		rw.Header().Set("Content-Type", "application/json")

		// reqBody := models.RequestLoginBodyModel{}
		if err := json.NewDecoder(r.Body).Decode(&loginModel); err != nil{
			if err.Error() == "EOF"{
				rw.WriteHeader(http.StatusBadRequest)
				responseModel.BadRequest("No values in the request", "")
			}else{
				rw.WriteHeader(http.StatusInternalServerError)
				responseModel.ServerError("Server error encountered", "")
			}
			// encode the error message response
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// empty request body
		if (loginModel == request_models.LoginModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("No values in request body", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// this allows a single variable only to be used for the entire handler
		loginString, err := loginModel.ToStringNoPassword()
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		if err := loginModel.HasValidValues(); err != nil{
			fmt.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("Invalid Email Address and/or Password", loginString)
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// get the email model from database
		emailModel := user_models.EmailModel{}
		if err := emailModel.ReadByValue(db, strings.ToLower(loginModel.EmailAddress)); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// did not match any email from database
		// the model is empty
		if (emailModel == user_models.EmailModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("Invalid Email Address and/or Password", loginString)
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// check if email is linked to a user
		// get the linker
		emailLinkerModel := user_linker_models.EmailLinkerModel{}
		if err := emailLinkerModel.ReadByEmailId(db, emailModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		

		// did not match any
		if (emailLinkerModel == user_linker_models.EmailLinkerModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("Invalid Email Address and/or Password", loginString)
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// create a single user id variable to be used further down the execution flow
		userID := emailLinkerModel.UserID

		// get the password linker
		passwordLinkerModel := user_linker_models.PasswordLinkerModel{}
		if err := passwordLinkerModel.ReadByUserId(db, userID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		

		// the password linker should not be empty
		if (passwordLinkerModel == user_linker_models.PasswordLinkerModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// the password model should not be empty
		passwordModel := user_models.PasswordModel{}
		if err := passwordModel.ReadById(db, passwordLinkerModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// check the password
		if !hash_utils.PasswordCheck(passwordModel.Value, loginModel.Password){
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("Invalid Email Address and/or Password", loginString)
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// get the auth value to attach to jwt claims
		authLinkerModel := user_linker_models.AuthLinkerModel{}
		if err := authLinkerModel.ReadByUserId(db, userID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// this cannot be empty at all
		// every user has should have this		
		if (authLinkerModel == user_linker_models.AuthLinkerModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		authModel := user_models.AuthModel{}
		if err := authModel.ReadById(db, authLinkerModel.AuthID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// again this cannot be empty
		// this was assigned during signup
		if (authModel == user_models.AuthModel{}){
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
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
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
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
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		jwtModel := token_models.JWTModel{
			AccessToken: accessStringToken,
			RefreshToken: refreshStringToken,
		}

		stringTokens, err := jwtModel.Stringify()
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// Set token headers
		rw.Header().Set("X-Access-Token", accessStringToken)
		rw.Header().Set("X-Refresh-Token", refreshStringToken)
		rw.WriteHeader(http.StatusOK) // Success 200 return

		// Also return tokens in the body
		responseModel.OK("User successfully logged in", stringTokens)
		json.NewEncoder(rw).Encode(responseModel)

	}
}