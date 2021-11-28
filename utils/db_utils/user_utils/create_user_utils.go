package utils

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	utils "github.com/v-inmar/goweb/utils/hash_utils"
)

/*
Creates user in the db with passed in values
as well as auto generated ones
*/
func CreateUser(db *sql.DB, firstname string, lastname string, email string, password string) (string, error){

	// helper function
	failed := func (err error) (string, error) {
		return "", err
	}

	dbSession, err := db.Begin()
	if err != nil{
		return failed(err)
	}
	defer dbSession.Rollback()


	for {
		// generate 8 random characters
		upid := strings.Replace(uuid.NewString(),"-", "", -1)[0:8]

		// check it doesnt eixst yet in db
		existing_upid_id_row := dbSession.QueryRow("select id from upid_model where value=?", upid)
		var existing_upid_id int64
		if err := existing_upid_id_row.Scan(&existing_upid_id); err != nil{
			if err == sql.ErrNoRows{

				// variable that hold the datetime value
				var dt = time.Now().UTC()

				// create user
				user_model_result, err := dbSession.Exec("insert into user_model (date_created) values (?)", dt)
				if err != nil{
					return failed(err)
				}

				user_id, err := user_model_result.LastInsertId()
				if err != nil{
					return failed(err)
				}

				// create upid
				upid_model_result, err := dbSession.Exec("insert into upid_model (value, date_created) values (?,?)", upid, dt)
				if err != nil{
					return failed(err)
				}

				upid_id, err := upid_model_result.LastInsertId()
				if err != nil{
					return failed(err)
				}

				// create user and upid linker
				_, err = dbSession.Exec("insert into user_upid_linker_model (user_id, upid_id, date_created) values (?,?,?)", user_id, upid_id, dt)
				if err != nil{
					return failed(err)
				}

				// read and/or create email
				email_id_row := dbSession.QueryRow("select id from email_model where value=?", email)
				var email_id int64
				if err := email_id_row.Scan(&email_id); err != nil{
					if err == sql.ErrNoRows{
						email_id = -1
					}else{
						return failed(err)
					}
				}

				if email_id == -1{ // create email
					email_model_result, err := dbSession.Exec("insert into email_model (value, date_created) values (?,?)", strings.ToLower(email), dt)
					if err != nil{
						return failed(err)
					}

					// not re-declaring the email_id variable
					email_id, err = email_model_result.LastInsertId()
					if err != nil {
						return failed(err)
					}

				}

				// create user and email linker
				_, err = dbSession.Exec("insert into user_email_linker_model (user_id, email_id, date_created) values (?,?,?)", user_id, email_id, dt)
				if err != nil{
					return failed(err)
				}


				// read and/or create firstname
				fname_id_row := dbSession.QueryRow("select id from firstname_model where value=?", firstname)
				var fname_id int64
				if err := fname_id_row.Scan(&fname_id); err != nil{
					if err == sql.ErrNoRows{ // does not exist yet
						
						// create firstname
						fname_model_result, err := dbSession.Exec("insert into firstname_model (value, date_created) values (?,?)", firstname, dt)
						if err != nil {
							return failed(err)
						}

						fname_id, err = fname_model_result.LastInsertId()
						if err != nil {
							return failed(err)
						}
					}else{
						return failed(err)
					}
				}

				// create user and firstname linker
				_, err = dbSession.Exec("insert into user_firstname_linker_model (user_id, firstname_id, date_created) values (?,?,?)", user_id, fname_id, dt)
				if err != nil{
					return failed(err)
				}

				// read and/or create lastname
				lname_id_row := dbSession.QueryRow("select id from lastname_model where value=?", lastname)
				var lname_id int64
				if err := lname_id_row.Scan(&lname_id); err != nil{
					if err == sql.ErrNoRows{

						// create lastname
						lname_model_result, err := dbSession.Exec("insert into lastname_model (value, date_created) values (?,?)", lastname, dt)
						if err != nil {
							return failed(err)
						}

						lname_id, err = lname_model_result.LastInsertId()
						if err != nil {
							return failed(err)
						}
					}else{
						return failed(err)
					}
				}

				// create user and lastname linker
				_, err = dbSession.Exec("insert into user_lastname_linker_model (user_id, lastname_id, date_created) values (?,?,?)", user_id, lname_id, dt)
				if err != nil{
					return failed(err)
				}

				// hashed the password
				hashed_string, err := utils.PasswordHash(password)
				if err != nil{
					return failed(err)
				}

				// create password
				pword_model_result, err := dbSession.Exec("insert into password_model (value, date_created) values (?,?)", hashed_string, dt)
				if err != nil{
					return failed(err)
				}

				pword_id, err := pword_model_result.LastInsertId()
				if err != nil {
					return failed(err)
				}

				// create user and password linker
				_, err = dbSession.Exec("insert into user_password_linker_model (user_id, password_id, date_created) values (?,?,?)", user_id, pword_id, dt)
				if err != nil {
					return failed(err)
				}

				if err := dbSession.Commit(); err != nil {
					return failed(err)
				}

				return upid, nil

			}else{
				continue
			}
		} // reaching else means it already exist
	}
	


}
