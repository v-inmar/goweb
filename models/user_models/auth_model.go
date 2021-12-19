package user_models

import (
	"database/sql"
	"time"
)

type AuthModel struct{
	ID int64
	Value string
	DateCreated time.Time
}

// Create new auth_model row in the database using transaction for rolling back if error occurs
func (m *AuthModel) Create(dbSession *sql.Tx, value string) error {
	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into auth_model (value, date_created) values (?, ?)", value, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = insertedID
	m.Value = value
	m.DateCreated = dt
	return nil
}

// Read row depending on the value passed in
func (m *AuthModel) ReadByValue(db *sql.DB, value string) error {
	err := db.QueryRow("select * from auth_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}