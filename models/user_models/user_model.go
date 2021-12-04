package user_models

import (
	"database/sql"
	"time"
)

type UserModel struct{
	ID int64
	DateCreated time.Time
}

func (m *UserModel) create(db *sql.DB) error {
	dbSession, err := db.Begin()
	if err != nil{
		return err
	}
	defer dbSession.Rollback()

	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into user_model (date_created) values (?)", dt)
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
	m.DateCreated = dt
	return nil
}

func (m *UserModel) readById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from user_model where id=?", id).Scan(&m.ID, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}