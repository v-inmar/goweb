package models

import (
	"database/sql"
	"time"
)

type LastnameLinkerModel struct {
	ID int64
	UserID int64
	LastnameID int64
	DateCreated time.Time
	DateUpdated time.Time
}

func (m *LastnameLinkerModel) create(db *sql.DB, userID, firstnameID int64) error {
	dbSession, err := db.Begin()
	if err != nil{
		return err
	}
	defer dbSession.Rollback()
	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into user_lastname_linker_model (user_id, firstname_id, date_created) values (?,?,?)", userID, firstnameID, dt)
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
	m.LastnameID = firstnameID
	m.DateCreated = dt
	return nil
}

func (m *LastnameLinkerModel) readById(db *sql.DB, id int64) error {
	model := db.QueryRow("select * from user_lastname_linker_model where id=?", id)
	if err := model.Scan(&m.ID, &m.UserID, &m.LastnameID, &m.DateCreated, &m.DateUpdated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *LastnameLinkerModel) readByUserId(db *sql.DB, userID int64) error {
	model := db.QueryRow("select * from user_lastname_linker_model where user_id=?", userID)
	if err := model.Scan(&m.ID, &m.UserID, &m.LastnameID, &m.DateCreated, &m.DateUpdated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}