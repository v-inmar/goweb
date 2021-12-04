package models

import (
	"database/sql"
	"time"
)

type EmailLinkerModel struct {
	ID int64
	UserID int64
	EmailID int64
	DateCreated time.Time
	DateUpdated time.Time
}

func (m *EmailLinkerModel) create(db *sql.DB, userID, emailID int64) error {
	dbSession, err := db.Begin()
	if err != nil{
		return err
	}
	defer dbSession.Rollback()

	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into user_email_linker_model (user_id, email_id, date_created) values (?,?,?)", userID, emailID, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = insertedID
	m.UserID = userID
	m.EmailID = emailID
	m.DateCreated = dt
	return nil
}

func (m *EmailLinkerModel) readById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from user_email_linker_model where id=?", id).Scan(&m.ID, &m.UserID, &m.EmailID, &m.DateCreated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *EmailLinkerModel) readByUserId(db *sql.DB, userID int64) error {
	err := db.QueryRow("select * from user_email_linker_model where user_id=?", userID).Scan(&m.ID, &m.UserID, &m.EmailID, &m.DateCreated, &m.DateUpdated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *EmailLinkerModel) readByEmailId(db *sql.DB, emailID int64) error {
	err := db.QueryRow("select * from user_email_linker_model where email_id=?", emailID).Scan(&m.ID, &m.UserID, &m.EmailID, &m.DateCreated, &m.DateUpdated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}