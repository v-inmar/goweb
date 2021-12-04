package models

import (
	"database/sql"
	"time"
)

type PasswordLinkerModel struct {
	ID int64
	UserID int64
	PasswordID int64
	DateCreated time.Time
	DateUpdated time.Time
}

func (m *PasswordLinkerModel) create(db *sql.DB, userID, passwordID int64) error {
	dbSession, err := db.Begin()
	if err != nil{
		return err
	}
	defer dbSession.Rollback()
	dt := time.Now().UTC()


	model, err := dbSession.Exec("insert into user_password_linker_model (user_id, password_id, date_created) values (?,?,?)", userID, passwordID, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	if err := dbSession.Commit(); err != nil{
		return err
	}

	m.ID = insertedID
	m.UserID = userID
	m.PasswordID = passwordID
	m.DateCreated = dt
	return nil
}

func (m *PasswordLinkerModel) readById(db *sql.DB, id int64) error {
	model := db.QueryRow("select * from user_password_linker_model where id=?", id)
	if err := model.Scan(&m.ID, &m.UserID, &m.PasswordID, &m.DateCreated, &m.DateUpdated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil

}

func (m *PasswordLinkerModel) readByUserId(db *sql.DB, userID int64) error {
	model := db.QueryRow("select * from user_password_linker_model where user_id=?", userID)
	if err := model.Scan(&m.ID, &m.UserID, &m.PasswordID, &m.DateCreated, &m.DateUpdated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}