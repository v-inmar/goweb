package user_models

import (
	"database/sql"
	"time"
)

type BlacklistModel struct {
	ID int64
	Value string
	DateCreated time.Time
}

/*
Create/Insert the json token into the blacklist model
*/
func (m *BlacklistModel) Create(dbSession *sql.Tx, value string) error {
	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into blacklist_model (value, date_created) values (?, ?)", value, dt)
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