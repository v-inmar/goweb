package user_models

import (
	"database/sql"
	"time"
)

type EmailModel struct {
	ID int64
	Value string
	DateCreated time.Time
}


// Pass in transaction to allow single rollback when any creating user
func (m *EmailModel) Create(dbSession *sql.Tx, value string) error {
	// dbSession, err := db.Begin()
	// if err != nil{
	// 	return err
	// }
	// defer dbSession.Rollback()

	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into email_model (value, date_created) values (?, ?)", value, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	// Cannot do commit here in case it would be a series of database inserts that gets called outside of this function
	// and one failed. This needs to be rollback too
	// if err := transaction.Commit(); err != nil{
	// 	return err
	// }
	m.ID = insertedID
	m.Value = value
	m.DateCreated = dt
	return nil
}

func (m *EmailModel) ReadByValue(db *sql.DB, value string) error {
	err := db.QueryRow("select * from email_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *EmailModel) ReadById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from email_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}