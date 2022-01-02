package handler_auth

import (
	"database/sql"
	"encoding/json"
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
	"github.com/v-inmar/goweb/utils/random_utils"
)

func SignupAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		responseModel := response_models.StatusModel{}
		signupModel := request_models.SignupModel{}
		addr := r.RemoteAddr
		// respBodErrMsgModel := models.ResponseBodyErrorMessageModel{}
		rw.Header().Set("Content-Type", "application/json")

		// reqBody := models.RequestSignUpBodyModel{}
		if err := json.NewDecoder(r.Body).Decode(&signupModel); err != nil{

			// no request body
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

		// Check for empty request body
		if (signupModel == request_models.SignupModel{}){
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest("No values in request body", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// a single string to be used as a return payload
		// without exposing the password
		signupString, err := signupModel.ToStringNoPassword()
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		
		/*
		Make sure all fields have values
		*/
		if err := signupModel.HasValidValues(); err != nil{
			rw.WriteHeader(http.StatusBadRequest)
			responseModel.BadRequest(err.Error(), signupString)
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		/*
		Database checks
		*/
		emailModel := user_models.EmailModel{}
		err = emailModel.ReadByValue(db, strings.ToLower(signupModel.EmailAddress))
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		
		if (emailModel != user_models.EmailModel{}){
			emailLinkerModel := user_linker_models.EmailLinkerModel{}
			err := emailLinkerModel.ReadByEmailId(db, emailModel.ID)
			if err != nil{
				rw.WriteHeader(http.StatusInternalServerError)
				responseModel.ServerError("Server error encountered", "")
				json.NewEncoder(rw).Encode(responseModel)
				return
			}
			// Compare to empty..not empty, means email not available
			if (emailLinkerModel != user_linker_models.EmailLinkerModel{}){
				rw.WriteHeader(http.StatusConflict)
				responseModel.Conflict("Email Address is not available", signupString)
				json.NewEncoder(rw).Encode(responseModel)
				return	
			}
		}

		// Hash the incoming password
		hashed, err := hash_utils.PasswordHash(signupModel.Password)
		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// ----- If it gets here, the email address is available ----- //


		dbSession, err := db.Begin() // begin transaction session
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}
		defer dbSession.Rollback()

		// ### User Model ### //
		userModel := user_models.UserModel{}
		if err := userModel.Create(dbSession); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// ### Email Model ### //
		if (emailModel == user_models.EmailModel{}){
			if err := emailModel.Create(dbSession, strings.ToLower(signupModel.EmailAddress)); err != nil{
				rw.WriteHeader(http.StatusInternalServerError)
				responseModel.ServerError("Server error encountered", "")
				json.NewEncoder(rw).Encode(responseModel)
				return
			}
		}

		emailLinkerModel := user_linker_models.EmailLinkerModel{}
		if err := emailLinkerModel.Create(dbSession, userModel.ID, emailModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}


		// ### Password Model ### //
		passwordModel := user_models.PasswordModel{}
		if err := passwordModel.Create(dbSession, hashed); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		passwordLinkerModel := user_linker_models.PasswordLinkerModel{}
		if err := passwordLinkerModel.Create(dbSession, userModel.ID, passwordModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}


		// ### Firstname Model ### //
		firstnameModel := user_models.FirstnameModel{}
		if err := firstnameModel.ReadByValue(db, signupModel.Firstname); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// Compare to empty
		if (firstnameModel == user_models.FirstnameModel{}) {
			// Create new
			if err := firstnameModel.Create(dbSession, signupModel.Firstname); err != nil{
				rw.WriteHeader(http.StatusInternalServerError)
				responseModel.ServerError("Server error encountered", "")
				json.NewEncoder(rw).Encode(responseModel)
				return
			}
		}

		firstnameLinkerModel := user_linker_models.FirstnameLinkerModel{}
		if err := firstnameLinkerModel.Create(dbSession, userModel.ID, firstnameModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// ### Lastname Model ### //
		lastnameModel := user_models.LastnameModel{}
		if err := lastnameModel.ReadByValue(db, signupModel.Lastname); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// Compare to empty
		if (lastnameModel == user_models.LastnameModel{}) {
			// Create new
			if err := lastnameModel.Create(dbSession, signupModel.Lastname); err != nil{
				rw.WriteHeader(http.StatusInternalServerError)
				responseModel.ServerError("Server error encountered", "")
				json.NewEncoder(rw).Encode(responseModel)
				return
			}
		}

		lastnameLinkerModel := user_linker_models.LastnameLinkerModel{}
		if err := lastnameLinkerModel.Create(dbSession, userModel.ID, lastnameModel.ID); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// string generator
		randStringGen := random_utils.RandomString{}

		// ### UPID Model ### //
		success := false
		// this will only run it 5 times
		// to make sure loop doesn't run forever
		for count := 0; count < 5; count++{
			if err := randStringGen.GenerateForUPID(); err != nil{
				continue
			}

			upidModel := user_models.UPIDModel{}

			if err := upidModel.ReadByValue(db, randStringGen.Value); err != nil{
				continue
			}

			// Compare to empty
			if (upidModel == user_models.UPIDModel{}){
				if err := upidModel.Create(dbSession, randStringGen.Value); err != nil{
					continue
				}

				upidLinkerModel := user_linker_models.UPIDLinkerModel{}
				if err := upidLinkerModel.Create(dbSession, userModel.ID, upidModel.ID); err != nil{
					continue
				}
				success = true
				break
			}

		}

		// for loop finished and check for upid success
		if !success{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}


		// ### Auth Model ### //
		auth_value := ""
		success = false // reset from previous use
		// will try 5 times only
		for count := 0; count < 5; count++{
			if err := randStringGen.GenerateAuth(); err != nil{
				continue
			}

			authModel := user_models.AuthModel{}
			if err := authModel.ReadByValue(db, randStringGen.Value); err != nil{
				continue
			}

			// check if empty
			if (authModel == user_models.AuthModel{}){
				// empty means its ok to use the generated string
				if err := authModel.Create(dbSession, randStringGen.Value); err != nil{
					continue
				}

				// create the linker model
				authLinkerModel := user_linker_models.AuthLinkerModel{}
				if err := authLinkerModel.Create(dbSession, userModel.ID, authModel.ID); err != nil{
					continue
				}
				success = true
				auth_value = authModel.Value
				break
			}
		}

		// for loop finished and check for upid success
		if !success{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}



		

		// Create an access jwt
		// Claims only for non-production
		// For prod, produce better claims
		accessStringToken, err := jwt_utils.GenerateJWT(jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 60).Unix(), // 1 hour
			"addr": addr, // client's address (ip)
			"uid": auth_value, // user identification
		})

		if err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}

		// Create a refresh jwt
		// Claims only for non-production
		// For prod, produce better claims
		refreshStringToken, err := jwt_utils.GenerateJWT(jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7days
			"addr": addr, // client's address (ip)
			"uid": auth_value, // user identification
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

		// Commit the transaction
		if err := dbSession.Commit(); err != nil{
			rw.WriteHeader(http.StatusInternalServerError)
			responseModel.ServerError("Server error encountered", "")
			json.NewEncoder(rw).Encode(responseModel)
			return
		}


		// ## - Write to the Response Header the Access and Refresh Token - ## //
		rw.Header().Set("X-Access-Token", accessStringToken)
		rw.Header().Set("X-Refresh-Token", refreshStringToken)
		rw.WriteHeader(http.StatusCreated) // Success 201 return
		// Also return tokens in the body
		responseModel.Created("User successfully registered", stringTokens)
		json.NewEncoder(rw).Encode(responseModel)
	}
}
