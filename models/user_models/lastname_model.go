package models

import (
	"database/sql"
	"time"
)

type LastnameModel struct {
	ID int64
	Value string
	DateCreated time.Time
}



func (m *LastnameModel) create(db *sql.DB, value string) error {
	dbSession, err := db.Begin()
	if err != nil{
		return err
	}
	defer dbSession.Rollback()

	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into lastname_model (value, date_created) values (?, ?)", value, dt)
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
	m.Value = value
	m.DateCreated = dt
	return nil
}

func (m *LastnameModel) readByValue(db *sql.DB, value string) error {
	err := db.QueryRow("select * from lastname_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *LastnameModel) readById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from lastname_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}